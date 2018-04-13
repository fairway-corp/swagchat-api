package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/swagchat/chat-api/datastore"
	"github.com/swagchat/chat-api/logging"
	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/notification"
	"github.com/swagchat/chat-api/utils"
)

// GetDevices is get devices
func GetDevices(ctx context.Context, userID string) (*models.Devices, *models.ProblemDetail) {
	devices, err := datastore.Provider(ctx).SelectDevices(userID)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get device failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}

	return &models.Devices{
		Devices: devices,
	}, nil
}

// GetDevice is get device
func GetDevice(ctx context.Context, userID string, platform int) (*models.Device, *models.ProblemDetail) {
	device, pd := selectDevice(ctx, userID, platform)
	if pd != nil {
		return nil, pd
	}

	return device, nil
}

// PutDevice is put device
func PutDevice(ctx context.Context, put *models.Device) (*models.Device, *models.ProblemDetail) {
	if pd := put.IsValid(); pd != nil {
		return nil, pd
	}

	// User existence check
	_, pd := selectUser(ctx, put.UserId)
	if pd != nil {
		return nil, pd
	}

	isExist := true
	device, pd := selectDevice(ctx, put.UserId, put.Platform)
	if device == nil {
		isExist = false
	}

	if !isExist || (device.Token != put.Token) {
		// When using another user on the same device, delete the notification information
		// of the olderuser in order to avoid duplication of the device token
		deleteDevices, err := datastore.Provider(ctx).SelectDevicesByToken(put.Token)
		if err != nil {
			pd := &models.ProblemDetail{
				Title:  "Update device failed",
				Status: http.StatusInternalServerError,
				Error:  err,
			}
			return nil, pd
		}
		if deleteDevices != nil {
			wg := &sync.WaitGroup{}
			for _, deleteDevice := range deleteDevices {
				nRes := <-notification.Provider().DeleteEndpoint(deleteDevice.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					return nil, nRes.ProblemDetail
				}
				err := datastore.Provider(ctx).DeleteDevice(deleteDevice.UserId, deleteDevice.Platform)
				if err != nil {
					pd := &models.ProblemDetail{
						Title:  "Update device failed",
						Status: http.StatusInternalServerError,
						Error:  err,
					}
					return nil, pd
				}
				wg.Add(1)
				go unsubscribeByDevice(ctx, deleteDevice, wg)
			}
			wg.Wait()
		}

		nRes := <-notification.Provider().CreateEndpoint(put.UserId, put.Platform, put.Token)
		if nRes.ProblemDetail != nil {
			return nil, nRes.ProblemDetail
		}
		put.NotificationDeviceId = put.Token
		if nRes.Data != nil {
			put.NotificationDeviceId = *nRes.Data.(*string)
		}

		if isExist {
			err := datastore.Provider(ctx).UpdateDevice(put)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Update device failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				return nil, pd
			}
			nRes = <-notification.Provider().DeleteEndpoint(device.NotificationDeviceId)
			if nRes.ProblemDetail != nil {
				return nil, nRes.ProblemDetail
			}
			go func() {
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go unsubscribeByDevice(ctx, device, wg)
				wg.Wait()
				go subscribeByDevice(ctx, put, nil)
			}()
		} else {
			device, err = datastore.Provider(ctx).InsertDevice(put)
			if err != nil {
				pd := &models.ProblemDetail{
					Title:  "Update device failed",
					Status: http.StatusInternalServerError,
					Error:  err,
				}
				return nil, pd
			}
			go subscribeByDevice(ctx, device, nil)
		}
		return device, nil
	}

	return nil, nil
}

// DeleteDevice is delete device
func DeleteDevice(ctx context.Context, userID string, platform int) *models.ProblemDetail {
	// User existence check
	_, pd := selectUser(ctx, userID)
	if pd != nil {
		return pd
	}

	device, pd := selectDevice(ctx, userID, platform)
	if pd != nil {
		return pd
	}

	np := notification.Provider()
	nRes := <-np.DeleteEndpoint(device.NotificationDeviceId)
	if nRes.ProblemDetail != nil {
		return nRes.ProblemDetail
	}

	err := datastore.Provider(ctx).DeleteDevice(userID, platform)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Delete device failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return pd
	}

	go unsubscribeByDevice(ctx, device, nil)

	return nil
}

func selectDevice(ctx context.Context, userID string, platform int) (*models.Device, *models.ProblemDetail) {
	device, err := datastore.Provider(ctx).SelectDevice(userID, platform)
	if err != nil {
		pd := &models.ProblemDetail{
			Title:  "Get device failed",
			Status: http.StatusInternalServerError,
			Error:  err,
		}
		return nil, pd
	}
	if device == nil {
		return nil, &models.ProblemDetail{
			Title:  "Resource not found",
			Status: http.StatusNotFound,
		}
	}
	return device, nil
}

func subscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	roomUser, err := datastore.Provider(ctx).SelectRoomUsersByUserID(device.UserId)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	if roomUser != nil {
		<-subscribe(ctx, roomUser, device)
	}
	if wg != nil {
		wg.Done()
	}
}

func unsubscribeByDevice(ctx context.Context, device *models.Device, wg *sync.WaitGroup) {
	subscriptions, err := datastore.Provider(ctx).SelectDeletedSubscriptionsByUserIDAndPlatform(device.UserId, device.Platform)
	if err != nil {
		logging.Log(zapcore.ErrorLevel, &logging.AppLog{
			Error: err,
		})
	}
	<-unsubscribe(ctx, subscriptions)
	if wg != nil {
		wg.Done()
	}
}

func subscribe(ctx context.Context, roomUsers []*models.RoomUser, device *models.Device) chan bool {
	np := notification.Provider()
	dp := datastore.Provider(ctx)
	doneCh := make(chan bool, 1)
	pdCh := make(chan *models.ProblemDetail, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, roomUser := range roomUsers {
		ctx = context.WithValue(ctx, utils.CtxRoomUser, roomUser)
		d.Work(ctx, func(ctx context.Context) {
			ru := ctx.Value(utils.CtxRoomUser).(*models.RoomUser)
			room, pd := selectRoom(ctx, ru.RoomId)
			if pd != nil {
				pdCh <- pd
			} else {
				if room.NotificationTopicId == "" {
					notificationTopicID, pd := createTopic(room.RoomId)
					if pd != nil {
						pdCh <- pd
					}

					room.NotificationTopicId = notificationTopicID
					room.Modified = time.Now().Unix()
					_, err := datastore.Provider(ctx).UpdateRoom(room)
					if err != nil {
						pd := &models.ProblemDetail{
							Status: http.StatusInternalServerError,
							Title:  "Update room failed",
						}
						pdCh <- pd
					}
				}
				nRes := <-np.Subscribe(room.NotificationTopicId, device.NotificationDeviceId)
				if nRes.ProblemDetail != nil {
					pdCh <- nRes.ProblemDetail
				} else {
					if nRes.Data != nil {
						notificationSubscriptionID := nRes.Data.(*string)
						subscription := &models.Subscription{
							RoomId:                     ru.RoomId,
							UserId:                     ru.UserId,
							Platform:                   device.Platform,
							NotificationSubscriptionId: *notificationSubscriptionID,
						}
						subscription, err := dp.InsertSubscription(subscription)
						if err != nil {
							pd := &models.ProblemDetail{
								Title:  "User registration failed",
								Status: http.StatusInternalServerError,
								Error:  err,
							}
							pdCh <- pd
						} else {
							doneCh <- true
						}
					}
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case pd := <-pdCh:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         pd.Error,
				})
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}

func unsubscribe(ctx context.Context, subscriptions []*models.Subscription) chan bool {
	np := notification.Provider()
	dp := datastore.Provider(ctx)
	doneCh := make(chan bool, 1)
	pdCh := make(chan *models.ProblemDetail, 1)
	finishCh := make(chan bool, 1)

	d := utils.NewDispatcher(10)
	for _, subscription := range subscriptions {
		ctx = context.WithValue(ctx, utils.CtxSubscription, subscription)
		d.Work(ctx, func(ctx context.Context) {
			targetSubscription := ctx.Value(utils.CtxSubscription).(*models.Subscription)
			nRes := <-np.Unsubscribe(targetSubscription.NotificationSubscriptionId)
			if nRes.ProblemDetail != nil {
				pdCh <- nRes.ProblemDetail
			}
			err := dp.DeleteSubscription(targetSubscription)
			if err != nil {
				pd := &models.ProblemDetail{
					Error: err,
				}
				pdCh <- pd
			} else {
				doneCh <- true
			}

			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case pd := <-pdCh:
				logging.Log(zapcore.ErrorLevel, &logging.AppLog{
					ProblemDetail: pd,
					Error:         pd.Error,
				})
				return
			}
		})
	}
	d.Wait()
	finishCh <- true
	return finishCh
}

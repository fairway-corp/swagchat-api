package rest

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/betchi/tracer"
	"github.com/go-zoo/bone"
	"github.com/swagchat/chat-api/model"
	"github.com/swagchat/chat-api/service"
	scpb "github.com/swagchat/protobuf/protoc-gen-go"
)

func setUserMux() {
	mux.PostFunc("/users", commonHandler(adminAuthzHandler(postUser)))
	mux.GetFunc("/users", commonHandler(adminAuthzHandler(getUsers)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$", commonHandler(adminAuthzHandler(selfResourceAuthzHandler(getUser))))
	mux.PutFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(putUser)))
	mux.DeleteFunc("/users/#userId^[a-z0-9-]$", commonHandler(selfResourceAuthzHandler(deleteUser)))

	// mux.GetFunc("/users/#userId^[a-z0-9-]$/unreadCount", commonHandler(selfResourceAuthzHandler(getUserUnreadCount)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/rooms", commonHandler(selfResourceAuthzHandler(getUserRooms)))
	mux.GetFunc("/users/#userId^[a-z0-9-]$/contacts", commonHandler(selfResourceAuthzHandler(getContacts)))
	mux.GetFunc("/profiles/#userId^[a-z0-9-]$", commonHandler(contactsAuthzHandler(getProfile)))
	mux.GetFunc("/roles/#roleId^[0-9]$/users", commonHandler(adminAuthzHandler(getRoleUsers)))
}

func postUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "postUser", "rest")
	defer tracer.Finish(span)

	var req model.CreateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	user, errRes := service.CreateUser(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusCreated, "application/json", user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getUsers", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveUsersRequest{}
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	limit, offset, _, _, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	users, errRes := service.RetrieveUsers(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getUser", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, errRes := service.RetrieveUser(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "putUser", "rest")
	defer tracer.Finish(span)

	var req model.UpdateUserRequest
	if err := decodeBody(r, &req); err != nil {
		respondJSONDecodeError(w, r, "")
		return
	}

	req.UserID = bone.GetValue(r, "userId")

	user, errRes := service.UpdateUser(ctx, &req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "deleteUser", "rest")
	defer tracer.Finish(span)

	req := &model.DeleteUserRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	errRes := service.DeleteUser(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusNoContent, "", nil)
}

func getUserRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getUserRooms", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveUserRoomsRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	limit, offset, _, _, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	if filterArray, ok := params["filter"]; ok {
		filter, err := strconv.Atoi(filterArray[0])
		if err != nil {
			invalidParams := []*scpb.InvalidParam{
				&scpb.InvalidParam{
					Name:   "filter",
					Reason: "filter is incorrect.",
				},
			}
			errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
			respondError(w, r, errRes)
			return
		}
		switch int32(filter) {
		case int32(scpb.UserRoomsFilter_Online):
			req.Filter = scpb.UserRoomsFilter_Online
		case int32(scpb.UserRoomsFilter_Unread):
			req.Filter = scpb.UserRoomsFilter_Unread
		}
	}

	roomUsers, errRes := service.RetrieveUserRooms(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", roomUsers)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getContacts", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveContactsRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		errRes := model.NewErrorResponse("", http.StatusBadRequest, model.WithError(err))
		respondError(w, r, errRes)
		return
	}

	limit, offset, _, _, orders, errRes := setPagingParams(params)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	req.Limit = limit
	req.Offset = offset
	req.Orders = orders

	contacts, errRes := service.RetrieveContacts(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", contacts)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getProfile", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveProfileRequest{}

	userID := bone.GetValue(r, "userId")
	req.UserID = userID

	user, errRes := service.RetrieveProfile(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}

	respond(w, r, http.StatusOK, "application/json", user)
}

func getRoleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := tracer.StartSpan(ctx, "getRoleUsers", "rest")
	defer tracer.Finish(span)

	req := &model.RetrieveRoleUsersRequest{}

	roleIDString := bone.GetValue(r, "roleId")
	roleIDInt, err := strconv.ParseInt(roleIDString, 10, 32)
	if err != nil {
		invalidParams := []*scpb.InvalidParam{
			&scpb.InvalidParam{
				Name:   "roleId",
				Reason: "roleId must be numeric.",
			},
		}
		errRes := model.NewErrorResponse("Failed to get userIds of user role.", http.StatusBadRequest, model.WithInvalidParams(invalidParams))
		respondError(w, r, errRes)
		return
	}

	roleIDInt32 := int32(roleIDInt)
	req.RoleID = &roleIDInt32

	roleUsers, errRes := service.RetrieveRoleUsers(ctx, req)
	if errRes != nil {
		respondError(w, r, errRes)
		return
	}
	respond(w, r, http.StatusOK, "application/json", roleUsers)
}

// func getUserUnreadCount(w http.ResponseWriter, r *http.Request) {
// 	userID := bone.GetValue(r, "userId")

// 	userUnreadCount, pd := service.RetrieveUserUnreadCount(ctx, userID)
// 	if pd != nil {
// 		respondErr(w, r, pd.Status, pd)
// 		return
// 	}

// 	respond(w, r, http.StatusOK, "application/json", userUnreadCount)
// }

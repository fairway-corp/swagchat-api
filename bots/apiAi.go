package bots

import (
	"net/http"
	"net/url"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/swagchat/chat-api/models"
	"github.com/swagchat/chat-api/utils"
)

type ApiAiProvider struct {
}

func (p *ApiAiProvider) Post(m *models.Message, b *models.Bot, c utils.JSONText) BotResult {
	r := BotResult{}

	var message string
	switch m.Type {
	case "text":
		var payloadText models.PayloadText
		json.Unmarshal(m.Payload, &payloadText)
		message = payloadText.Text
	case "image":
		message = "画像を受信しました"
	default:
		message = "メッセージを受信しました"
	}

	var cred models.ApiAiCredencial
	json.Unmarshal(c, &cred)

	values := url.Values{}
	values.Set("v", "20150910")
	values.Add("timezone", "Asia/Tokyo")
	values.Add("lang", "ja")
	values.Add("sessionId", b.UserId)
	values.Add("query", message)
	log.Println(values.Encode())
	req, err := http.NewRequest(
		"GET",
		"https://api.api.ai/v1/query?"+values.Encode(),
		nil,
	)
	if err != nil {
	}
	req.Header.Set("Authorization", utils.AppendStrings("Bearer ", cred.ClientAccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
	}
	var res models.ApiAiResponse

	log.Printf("%#v\n", resp)
	//json.NewDecoder(resp.Body).Decode(res)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	log.Printf("%#v\n", string(body))
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("=================================")
	log.Printf("%#v\n", res)

	//if len(res.Results) > 0 {
	var textPayload utils.JSONText
	err = json.Unmarshal([]byte("{\"text\": \""+res.Result.Fulfillment.Speech+"\"}"), &textPayload)
	post := &models.Message{
		RoomId:  m.RoomId,
		UserId:  b.UserId,
		Type:    "text",
		Payload: textPayload,
	}
	posts := make([]*models.Message, 0)
	posts = append(posts, post)
	messages := &models.Messages{
		Messages: posts,
	}
	r.Messages = messages

	return r
}
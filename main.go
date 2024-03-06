package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"log"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	// load .env file

	// Line
	channelSecret := os.Getenv("ad77a9c6fe42d22cf61bedf3598bad8b")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("6UBugZb++eWul5dBiRjvyPVUWpLfv8AjDtMPT1ItbucizPSQiwQTt6vPPiSKBiRyTXhi+z60uK0IPAveE7nPJ+xLYicZOPP/xGzte0n4HWkBi/RnFlmzCQzN7w5j8XGKZVn44fNKZ3WSRjYEmoN4TwdB04t89/1O/w1cDnyilFU="),
	)
	if err != nil {
		log.Fatal(err)
	}
	// routess
	router.GET("/ping", ping)
	router.POST("/callback", lineCallback(bot, channelSecret))
	router.Run(":2500")
}

func lineCallback(bot *messaging_api.MessagingApiAPI, channelSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cb, err := webhook.ParseRequest(channelSecret, c.Request)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if err == linebot.ErrInvalidSignature {
				c.Status(400)
			} else {
				c.Status(500)
			}
			return
		}

		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:

					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: "testt",
								},
							},
						},
					); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent text reply.")
					}
				case webhook.StickerMessageContent:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
					if _, err = bot.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: replyMessage,
								},
							},
						}); err != nil {
						log.Print(err)
					} else {
						log.Println("Sent sticker reply.")
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			}
		}
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

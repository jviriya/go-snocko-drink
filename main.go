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
	"strconv"
	"strings"
)

var (
	orderList = map[string]int{}
	orderNo   []string
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
					if strings.HasPrefix(message.Text, "#") {
						resp := drinkCommand(message.Text)
						if resp != "" {
							if _, err = bot.ReplyMessage(
								&messaging_api.ReplyMessageRequest{
									ReplyToken: e.ReplyToken,
									Messages: []messaging_api.MessageInterface{
										messaging_api.TextMessage{
											Text: resp,
										},
									},
								},
							); err != nil {
								log.Print(err)
							} else {
								log.Println("Sent text reply.")
							}
						} else {
							log.Printf("Unsupported message content: %T\n", message.Text)
						}
					} else {
						log.Printf("Unsupported message content: %T\n", message.Text)
					}
				case webhook.StickerMessageContent:
					//replyMessage := fmt.Sprintf(
					//	"sticker id is %s, stickerResourceType is %s", message.StickerId, message.StickerResourceType)
					//if _, err = bot.ReplyMessage(
					//	&messaging_api.ReplyMessageRequest{
					//		ReplyToken: e.ReplyToken,
					//		Messages: []messaging_api.MessageInterface{
					//			messaging_api.TextMessage{
					//				Text: replyMessage,
					//			},
					//		},
					//	}); err != nil {
					//	log.Print(err)
					//} else {
					//	log.Println("Sent sticker reply.")
					//}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			}
		}
	}
}

func drinkCommand(command string) string {
	mainCommand := first2Char(command)
	switch mainCommand {
	case "#m", "#menu":
		return "MENU to be shown."
	case "#a", "#add":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l == 1 {
			return ""
		} else if l == 2 {
			if _, ok := orderList[splitCommands[1]]; !ok {
				orderNo = append(orderNo, splitCommands[1])
			}
			orderList[splitCommands[1]] += 1

		} else { //len == 3
			if _, ok := orderList[splitCommands[1]]; !ok {
				orderNo = append(orderNo, splitCommands[1])
			}
			quantity, err := strconv.Atoi(splitCommands[2])
			if err != nil {
				return ""
			}
			orderList[splitCommands[1]] += quantity
		}

	case "#r", "#rm", "#remove":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l == 1 {
			return ""
		} else if l == 2 {
			no, err := strconv.Atoi(splitCommands[1])
			if err != nil {
				return ""
			}
			delete(orderList, orderNo[no])
			orderNo = removeIndex(orderNo, no)
		} else { //len == 3
			no, err := strconv.Atoi(splitCommands[1])
			if err != nil {
				return ""
			}
			quantity, err := strconv.Atoi(splitCommands[2])
			if err != nil {
				return ""
			}
			if orderList[orderNo[no]] > quantity {
				orderList[orderNo[no]] -= quantity
			} else {
				delete(orderList, orderNo[no])
				orderNo = removeIndex(orderNo, no)
			}
		}
		return makeResponse()

	default:
		log.Printf("Unsupported message content: %T\n", command)
	}

	return ""
}

func makeResponse() string {
	resp := ""
	for i, v := range orderNo {
		resp += fmt.Sprintf("\n%d. %s %d", i, v, orderList[v])
	}
	return ""
}

func first2Char(s string) string {
	i := 0
	for j := range s {
		if i == 2 {
			return s[:j]
		}
		i++
	}
	return s
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// load .env file
	err := godotenv.Load("snockodrink.env")

	// Line

	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	//
	//fmt.Println("test: ")
	//
	//txt := "พ asdf"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "พ 1"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "พ asdfd"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "พ asdf 2"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "ล 1 2"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "ล 2 1"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "ล asdf 1"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "ล 1"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()
	//
	//txt = "พ asdf"
	//fmt.Printf("%s", drinkCommand(txt))
	//fmt.Println()

	// routess
	router.GET("/ping", ping)
	router.GET("/", handler)
	router.POST("/callback", lineCallback(bot, channelSecret))
	router.Run(":5000")
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
					//if strings.HasPrefix(message.Text, "#") {
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
					//} else {
					//	log.Printf("Unsupported message content: %T\n", message.Text)
					//}
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
	switch {
	case command == "เมนู", command == "menu":
		return "MENU to be shown."
	case command == "รายการ", command == "order":
		return makeResponse()
	case command == "เคลียร์", command == "clear":
		orderList = map[string]int{}
		orderNo = []string{}
		return makeResponse()
	case firstNChar(command, 2) == "พ ", firstNChar(command, 6) == "เพิ่ม ":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l == 1 {
			return ""
		} else if l == 2 {
			no, err := strconv.Atoi(splitCommands[1])
			if err == nil {
				no--
				if len(orderNo) > no {
					orderList[orderNo[no]] += 1
				}
			} else {
				if _, ok := orderList[splitCommands[1]]; !ok {
					orderNo = append(orderNo, splitCommands[1])
				}
				orderList[splitCommands[1]] += 1
			}

		} else if l == 3 {
			quantity, err := strconv.Atoi(splitCommands[2])
			if err != nil {
				return ""
			}
			no, err := strconv.Atoi(splitCommands[1])
			if err == nil {
				no--
				if len(orderNo) > no {
					orderList[orderNo[no]] += quantity
				}
			} else {
				if _, ok := orderList[splitCommands[1]]; !ok {
					orderNo = append(orderNo, splitCommands[1])
				}
				orderList[splitCommands[1]] += quantity
			}
		}

	case firstNChar(command, 2) == "ล ", firstNChar(command, 3) == "ลด ":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l == 1 {
			return ""
		} else if l == 2 {
			no, err := strconv.Atoi(splitCommands[1])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo { //find index
					if v == splitCommands[1] {
						no = i
						break
					}
				}
			}
			if _, ok := orderList[orderNo[no]]; ok {
				delete(orderList, orderNo[no])
				orderNo = removeIndex(orderNo, no)
			}
		} else if l == 3 {
			quantity, err := strconv.Atoi(splitCommands[2])
			if err != nil {
				return ""
			}
			no, err := strconv.Atoi(splitCommands[1])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo { //find index
					if v == splitCommands[1] {
						no = i
						break
					}
				}
			}
			if orderList[orderNo[no]] > quantity {
				orderList[orderNo[no]] -= quantity
			} else {
				if _, ok := orderList[orderNo[no]]; ok {
					delete(orderList, orderNo[no])
					orderNo = removeIndex(orderNo, no)
				}
			}
		}

	default:
		log.Printf("Unsupported message content: %T\n", command)
	}

	return makeResponse()
}

func makeResponse() string {
	resp := ""
	for i, v := range orderNo {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList[v])
	}
	return resp
}

func firstNChar(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
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

func handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

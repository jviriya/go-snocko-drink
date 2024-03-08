package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	orderList = map[string]map[string]int{
		"น": map[string]int{},
		"ข": map[string]int{},
		"ผ": map[string]int{},
	}
	orderNo = map[string][]string{
		"น": []string{},
		"ข": []string{},
		"ผ": []string{},
	}
	additionalMsg string
	bangkokTZ     *time.Location
)

func main() {
	gin.SetMode(gin.ReleaseMode)
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

	bangkokTZ, err = time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Print(err)
	}

	com := "พ ผ เทส 2"
	fmt.Println("TEST")
	fmt.Println(drinkCommand(com))

	com = "พ น เทส2 2"
	fmt.Println("TEST")
	fmt.Println(drinkCommand(com))

	com = "พ ข เทส3 2"
	fmt.Println("TEST")
	fmt.Println(drinkCommand(com))

	com = "ล ผ 1 1"
	fmt.Println("TEST")
	fmt.Println(drinkCommand(com))

	com = "clear"
	fmt.Println("TEST")
	fmt.Println(drinkCommand(com))

	c := cron.New(cron.WithLocation(bangkokTZ))

	c.AddFunc("30 11 * * *", func() {
		pushMessages(bot, "สั่งน้ำจ้าปิดบ่ายโมง!!!")
		orderList = map[string]map[string]int{
			"น": map[string]int{},
			"ข": map[string]int{},
			"ผ": map[string]int{},
		}
		orderNo = map[string][]string{
			"น": []string{},
			"ข": []string{},
			"ผ": []string{},
		}
		additionalMsg = ""
	})

	c.AddFunc("50 12 * * *", func() {
		pushMessages(bot, "อีก 10 นาทีปิดแล้วนาจา !!!")
	})

	c.AddFunc("55 12 * * *", func() {
		pushMessages(bot, "อีก 5 นาทีปิดแล้วนาจา !!!")
	})

	c.AddFunc("0 13 * * *", func() {
		pushMessages(bot, fmt.Sprintf("ปิดจ้าา! ขอสรุปออเดอร์\n\n%v", makeResponse()))
	})

	c.Start()

	router.GET("/ping", ping)
	router.GET("/", handler)
	router.POST("/callback", lineCallback(bot, channelSecret))
	router.Run(":5000")
	defer c.Stop()
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
						messages := []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: resp, //Modify text here
							},
							//messaging_api.TextMessage{
							//	Text: fmt.Sprintf("%v <%v>", e.Source.(webhook.GroupSource).GroupId, time.Now()),
							//},
						}

						if additionalMsg != "" {
							messages = append(messages, messaging_api.TextMessage{
								Text: additionalMsg,
							})
						}

						if resp != "" && isNotWeekend() {
							if _, err = bot.ReplyMessage(
								&messaging_api.ReplyMessageRequest{
									ReplyToken: e.ReplyToken,
									Messages:   messages,
								},
							); err != nil {
								log.Print(err)
							} else {
								log.Println("Sent text reply.")
							}
						} else {
							log.Printf("Unsupported message content: %T\n", message.Text)
						}
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
	command = strings.ToLower(command)
	additionalMsg = ""
	switch {
	case command == "เมนู", command == "menu":
		return "ดูในโน้ตเลยจ้า"
	case command == "รายการ", command == "order":
		additionalMsg = "ดูได้เลยจ้า"
		return makeResponse()
	case command == "เคลียร์", command == "clear":
		orderList = map[string]map[string]int{
			"น": map[string]int{},
			"ข": map[string]int{},
			"ผ": map[string]int{},
		}
		orderNo = map[string][]string{
			"น": []string{},
			"ข": []string{},
			"ผ": []string{},
		}
		additionalMsg = "clear แล้วจ้า"
		return makeResponse()
	case firstNChar(command, 2) == "พ ", firstNChar(command, 6) == "เพิ่ม ":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l < 3 {
			return ""
		} else if l == 3 {
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				no--
				if len(orderNo[splitCommands[1]]) > no {
					orderList[splitCommands[1]][orderNo[splitCommands[1]][no]] += 1
				}
			} else {
				if _, ok := orderList[splitCommands[1]][splitCommands[2]]; !ok {
					orderNo[splitCommands[1]] = append(orderNo[splitCommands[1]], splitCommands[2])
				}
				orderList[splitCommands[1]][splitCommands[2]] += 1
			}

		} else if l == 4 {
			quantity, err := strconv.Atoi(splitCommands[3])
			if err != nil {
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				no--
				if len(orderNo[splitCommands[1]]) > no {
					orderList[splitCommands[1]][orderNo[splitCommands[1]][no]] += quantity
				}
			} else {
				if _, ok := orderList[splitCommands[2]]; !ok {
					orderNo[splitCommands[1]] = append(orderNo[splitCommands[1]], splitCommands[2])
				}
				orderList[splitCommands[1]][splitCommands[2]] += quantity
			}
		}

	case firstNChar(command, 2) == "ล ", firstNChar(command, 3) == "ลบ ", firstNChar(command, 3) == "ลด ":
		additionalMsg = "รับทราบจ้า"
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l < 3 {
			return ""
		} else if l == 3 {
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo[splitCommands[1]] { //find index
					if v == splitCommands[2] {
						no = i
						break
					}
				}
			}
			if _, ok := orderList[splitCommands[1]][orderNo[splitCommands[1]][no]]; ok {
				delete(orderList[splitCommands[1]], orderNo[splitCommands[1]][no])
				orderNo[splitCommands[1]] = removeIndex(orderNo[splitCommands[1]], no)
			}
		} else if l == 4 {
			quantity, err := strconv.Atoi(splitCommands[3])
			if err != nil {
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo[splitCommands[1]] { //find index
					if v == splitCommands[2] {
						no = i
						break
					}
				}
			}
			if orderList[splitCommands[1]][orderNo[splitCommands[1]][no]] > quantity {
				orderList[splitCommands[1]][orderNo[splitCommands[1]][no]] -= quantity
			} else {
				if _, ok := orderList[orderNo[splitCommands[1]][no]]; ok {
					delete(orderList, orderNo[splitCommands[1]][no])
					orderNo[splitCommands[1]] = removeIndex(orderNo[splitCommands[1]], no)
				}
			}
		}

	default:
		log.Printf("Unsupported message content: %T\n", command)
		return ""
	}

	return makeResponse()
}

func makeResponse() string {
	resp := "รายการทั้งหมด\n\n"

	resp += "ไซส์ L\n---------------------"
	for i, v := range orderNo["น"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList["น"][v])
	}
	resp += "\n\n"

	resp += "ขนม\n---------------------"
	for i, v := range orderNo["ข"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList["ข"][v])
	}
	resp += "\n\n"

	resp += "น้ำผลไม้\n---------------------"
	for i, v := range orderNo["ผ"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList["ผ"][v])
	}
	resp += "\n"
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

func pushMessages(bot *messaging_api.MessagingApiAPI, message string) {
	if isNotWeekend() {
		_, err := bot.PushMessage(&messaging_api.PushMessageRequest{
			To: os.Getenv("GROUP_ID"),
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{Text: message},
			},
			NotificationDisabled:   true,
			CustomAggregationUnits: nil,
		}, "")
		if err != nil {
			log.Print(err)
		}
	}
}

func isNotWeekend() bool {
	now := time.Now().In(bangkokTZ)
	return now.Weekday() != time.Saturday && now.Weekday() != time.Sunday
}

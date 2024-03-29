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
	orderList     = map[string]map[string]map[string]int{}
	orderNo       = map[string]map[string][]string{}
	additionalMsg string
	bangkokTZ     *time.Location
	groupId       string
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

	//groupId = "test"
	//if _, ok := orderList[groupId]; !ok {
	//	orderList[groupId] = map[string]map[string]int{
	//		"น": map[string]int{},
	//		"ข": map[string]int{},
	//		"ผ": map[string]int{},
	//	}
	//}
	//if _, ok := orderNo[groupId]; !ok {
	//	orderNo[groupId] = map[string][]string{
	//		"น": []string{},
	//		"ข": []string{},
	//		"ผ": []string{},
	//	}
	//}
	//com := "พ น เทส 2\nพ z หยก 2"
	//
	//comArr := strings.Split(com, "\n")
	//
	//for _, v := range comArr {
	//	fmt.Println(drinkCommand(v))
	//}
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ น เทส2 2"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ น เทส1 2"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3 2"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3 2"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3 2"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "พ ตบขนมไทย เทส3"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//
	//com = "ล ผ 1 1"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))
	//fmt.Println("additionalMsg: " + additionalMsg)

	//com = "clear"
	//fmt.Println("TEST")
	//fmt.Println(drinkCommand(com))

	c := cron.New(cron.WithLocation(bangkokTZ))

	c.AddFunc("30 11 * * *", func() {
		pushMessages(bot, "สั่งน้ำจ้าปิดบ่ายโมง!!!")
		orderList = map[string]map[string]map[string]int{}
		orderNo = map[string]map[string][]string{}
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
					//groupId = e.Source.(webhook.GroupSource).GroupId
					//groupId = "test"
					if _, ok := orderList[groupId]; !ok {
						orderList[groupId] = map[string]map[string]int{
							"น": map[string]int{},
							"ข": map[string]int{},
							"ผ": map[string]int{},
						}
					}
					if _, ok := orderNo[groupId]; !ok {
						orderNo[groupId] = map[string][]string{
							"น": []string{},
							"ข": []string{},
							"ผ": []string{},
						}
					}
					texts := strings.Split(message.Text, "\n")

					var resp string
					for _, command := range texts {
						resp = drinkCommand(command)
					}

					messages := []messaging_api.MessageInterface{}
					if resp != "" {
						messages = []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: resp, //Modify text here
							},
						}
					}
					if additionalMsg != "" {
						messages = append(messages, messaging_api.TextMessage{
							Text: additionalMsg,
						})
					}

					if len(messages) > 0 && isNotWeekend() {
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
		return "ดูใน Albums เลยจ้า"
	case command == "รายการ", command == "order":
		additionalMsg = "ดูได้เลยจ้า"
		return makeResponse()
	case command == "เคลียร์", command == "clear":
		orderList[groupId] = map[string]map[string]int{
			"น": map[string]int{},
			"ข": map[string]int{},
			"ผ": map[string]int{},
		}
		orderNo[groupId] = map[string][]string{
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
			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				no--
				if len(orderNo[groupId][typ]) > no {
					orderList[groupId][typ][orderNo[groupId][typ][no]] += 1
				}
			} else {
				if _, ok := orderList[groupId][typ][splitCommands[2]]; !ok {
					orderNo[groupId][typ] = append(orderNo[groupId][typ], splitCommands[2])
				}
				orderList[groupId][typ][splitCommands[2]] += 1
			}
		} else if l == 4 {
			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			quantity, err := strconv.Atoi(splitCommands[3])
			if err != nil {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				no--
				if len(orderNo[groupId][typ]) > no {
					orderList[groupId][typ][orderNo[groupId][typ][no]] += quantity
				}
			} else {
				if _, ok := orderList[groupId][typ][splitCommands[2]]; !ok {
					orderNo[groupId][typ] = append(orderNo[groupId][typ], splitCommands[2])
				}
				orderList[groupId][typ][splitCommands[2]] += quantity
			}
			additionalMsg = fmt.Sprintf("รับออเดอร์จ้า %v จำนวน %v", splitCommands[2], quantity)

		} else {
			additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
			return ""
		}

	case firstNChar(command, 2) == "ล ", firstNChar(command, 3) == "ลบ ", firstNChar(command, 3) == "ลด ":
		additionalMsg = "รับทราบจ้า"
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)
		if l < 3 {
			return ""
		} else if l == 3 {
			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo[groupId][typ] { //find index
					if v == splitCommands[2] {
						no = i
						break
					}
				}
			}
			if _, ok := orderList[groupId][typ][orderNo[groupId][typ][no]]; ok {
				delete(orderList[groupId][typ], orderNo[groupId][typ][no])
				orderNo[groupId][typ] = removeIndex(orderNo[groupId][typ], no)
			}
		} else if l == 4 {
			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			quantity, err := strconv.Atoi(splitCommands[3])
			if err != nil {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return ""
			}
			no, err := strconv.Atoi(splitCommands[2])
			if err == nil { //remove by order number
				no--
			} else {
				for i, v := range orderNo[groupId][typ] { //find index
					if v == splitCommands[2] {
						no = i
						break
					}
				}
			}
			if len(orderNo[groupId][typ]) > 0 {
				if orderList[groupId][typ][orderNo[groupId][typ][no]] > quantity {
					orderList[groupId][typ][orderNo[groupId][typ][no]] -= quantity
				} else {
					if _, ok := orderList[groupId][orderNo[groupId][typ][no]]; ok {
						delete(orderList[groupId], orderNo[groupId][typ][no])
						orderNo[groupId][typ] = removeIndex(orderNo[groupId][typ], no)
					}
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

	resp += fmt.Sprintf("น้ำไซส์ L จำนวน (%d).\n---------------------", sumFn(orderList[groupId]["น"]))
	for i, v := range orderNo[groupId]["น"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList[groupId]["น"][v])
	}
	resp += "\n\n"

	resp += fmt.Sprintf("ขนม จำนวน (%d).\n---------------------", sumFn(orderList[groupId]["ข"]))
	for i, v := range orderNo[groupId]["ข"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList[groupId]["ข"][v])
	}
	resp += "\n\n"

	resp += fmt.Sprintf("น้ำผลไม้ จำนวน (%d).\n---------------------", sumFn(orderList[groupId]["ผ"]))
	for i, v := range orderNo[groupId]["ผ"] {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v, orderList[groupId]["ผ"][v])
	}
	resp += "\n"
	return resp
}

func convertType(typ string) string {
	switch {
	case typ == "ผ", strings.Contains(typ, "ผลไม้"):
		return "ผ"
	case typ == "น", strings.Contains(typ, "น้ำ"):
		return "น"
	case typ == "ข", strings.Contains(typ, "ขนม"):
		return "ข"
	}
	return ""
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

func sumFn(e map[string]int) int {
	sum := 0
	for _, m := range e {
		sum = sum + m
	}

	return sum
}

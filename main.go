package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type GroupOrderInterface interface {
	AddOrder(v []OrderDetails, t string)
	AddOrderByIndex(index, quantity int, t string)
	RemoveOrder(v OrderDetails, t string)
	RemoveOrderByIndex(index, quantity int, t string)
	ClearOrder()
	GetOrderNameByIndex(i int, t string) string
}

type Group struct {
	DrinksOrders []OrderDetails
	SnackOrders  []OrderDetails
	FruitOrders  []OrderDetails
	GroupId      string
}

type OrderDetails struct {
	Name     string
	Quantity int
}

var (
	Drink         = "D"
	Snack         = "S"
	Fruit         = "F"
	allGroup      = map[string]GroupOrderInterface{}
	bangkokTZ     *time.Location
	additionalMsg = ""
	reply         = false
)

func (g *Group) AddOrder(v []OrderDetails, t string) {
	switch t {
	case Drink:
		increaseOrder(v, &g.DrinksOrders)
	case Snack:
		increaseOrder(v, &g.SnackOrders)
	case Fruit:
		increaseOrder(v, &g.FruitOrders)
	}
}

func (g *Group) AddOrderByIndex(index, quantity int, t string) {
	switch t {
	case Drink:
		increaseOrderByIndex(index-1, quantity, &g.DrinksOrders)
	case Snack:
		increaseOrderByIndex(index-1, quantity, &g.SnackOrders)
	case Fruit:
		increaseOrderByIndex(index-1, quantity, &g.FruitOrders)
	}
}

func (g *Group) RemoveOrder(v OrderDetails, t string) {
	switch t {
	case Drink:
		decreaseOrder(v, &g.DrinksOrders)
	case Snack:
		decreaseOrder(v, &g.SnackOrders)
	case Fruit:
		decreaseOrder(v, &g.FruitOrders)
	}
}

func (g *Group) RemoveOrderByIndex(index, quantity int, t string) {
	switch t {
	case Drink:
		decreaseOrderByIndex(index-1, quantity, &g.DrinksOrders)
	case Snack:
		decreaseOrderByIndex(index-1, quantity, &g.SnackOrders)
	case Fruit:
		decreaseOrderByIndex(index-1, quantity, &g.FruitOrders)
	}
}

func (g *Group) ClearOrder() {
	g.DrinksOrders = []OrderDetails{}
	g.SnackOrders = []OrderDetails{}
	g.FruitOrders = []OrderDetails{}
}

func (g *Group) GetOrderNameByIndex(i int, t string) string {
	switch t {
	case Drink:
		return g.DrinksOrders[i-1].Name
	case Snack:
		return g.SnackOrders[i-1].Name
	case Fruit:
		return g.FruitOrders[i-1].Name
	}

	return ""
}

func (g *Group) String() string {
	resp := "รายการทั้งหมด\n\n"

	resp += fmt.Sprintf("น้ำไซส์ L จำนวน (%d).\n---------------------", sumFn(g.DrinksOrders))
	for i, v := range g.DrinksOrders {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v.Name, v.Quantity)
	}
	resp += "\n\n"

	resp += fmt.Sprintf("ขนม จำนวน (%d).\n---------------------", sumFn(g.SnackOrders))
	for i, v := range g.SnackOrders {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v.Name, v.Quantity)
	}
	resp += "\n\n"

	resp += fmt.Sprintf("น้ำผลไม้ จำนวน (%d).\n---------------------", sumFn(g.FruitOrders))
	for i, v := range g.FruitOrders {
		resp += fmt.Sprintf("\n%d. %s %d", i+1, v.Name, v.Quantity)
	}
	resp += "\n---------------------\n"
	resp += fmt.Sprintf("%v\n", g.GroupId)
	return resp
}

func newOrderByGroupId(groupId string) GroupOrderInterface {
	return &Group{GroupId: groupId}
}

func main() {
	//g1 := "group1"
	//allGroup[g1] = newOrderByGroupId(g1)
	//
	//drinkCommand("พ น นม 1", g1)
	//drinkCommand("พ น นม", g1)
	//drinkCommand("พ น 1", g1)
	//fmt.Println(additionalMsg)
	//
	//drinkCommand("ล น 1 1", g1)
	//drinkCommand("ล น นม 1", g1)

	//fmt.Println(allGroup[g1])
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

	c := cron.New(cron.WithLocation(bangkokTZ))

	c.AddFunc("59 23 * * *", func() {
		for gid, _ := range allGroup {
			allGroup[gid].ClearOrder()
		}
		additionalMsg = ""
	})

	c.AddFunc("30 10 * * *", func() {
		pushMessagesAllGroup(bot, "สั่งน้ำจ้าปิดเที่ยงครึ่งงง!!!")
	})

	//c.AddFunc("50 12 * * *", func() {
	//	pushMessagesAllGroup(bot, "อีก 10 นาทีปิดแล้วนาจา !!!")
	//})
	//
	//c.AddFunc("55 12 * * *", func() {
	//	pushMessagesAllGroup(bot, "อีก 5 นาทีปิดแล้วนาจา !!!")
	//})

	c.AddFunc("0 13 * * *", func() {
		pushMessagesOrderSummary(bot)
	})

	c.Start()

	router.GET("/ping", ping)
	router.GET("/", handler)
	router.POST("/callback", lineCallback(bot, channelSecret))
	router.Run(":5000")

	defer c.Stop()
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
					groupId := e.Source.(webhook.GroupSource).GroupId

					_, found := allGroup[groupId]
					if !found {
						allGroup[groupId] = newOrderByGroupId(groupId)
					}

					texts := strings.Split(message.Text, "\n")

					for _, command := range texts {
						drinkCommand(command, groupId)
					}

					messages := []messaging_api.MessageInterface{}
					messages = []messaging_api.MessageInterface{
						messaging_api.TextMessage{
							Text: fmt.Sprint(allGroup[groupId]), //Modify text here
						},
					}
					if additionalMsg != "" {
						messages = append(messages, messaging_api.TextMessage{
							Text: additionalMsg,
						})
					}

					if len(messages) > 0 && isNotWeekend() && reply == true {
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

				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			}
		}
	}
}

func drinkCommand(command, groupId string) {
	command = strings.ToLower(command)
	additionalMsg = ""
	reply = true
	switch {
	case command == "เมนู", command == "menu":
		additionalMsg = "ดูใน Albums เลยจ้า"
		return
	case command == "รายการ", command == "order":
		additionalMsg = "ดูได้เลยจ้า"
		return
	case command == "เคลียร์", command == "clear":
		allGroup[groupId].ClearOrder()
		additionalMsg = "clear แล้วจ้า"
		return
	case firstNChar(command, 2) == "พ ", firstNChar(command, 6) == "เพิ่ม ":
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)

		if l < 3 || l > 4 {
			additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
			return
		} else {
			// if len is 3 add quantity 1
			if l == 3 {
				splitCommands = append(splitCommands, "1")
			}
			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return
			}
			quantity, err := strconv.Atoi(splitCommands[3])
			if err != nil {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return
			}

			orderN, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				allGroup[groupId].AddOrderByIndex(orderN, quantity, typ)
				splitCommands[2] = allGroup[groupId].GetOrderNameByIndex(orderN, typ)
			} else {
				allGroup[groupId].AddOrder([]OrderDetails{
					{
						Name:     splitCommands[2],
						Quantity: quantity,
					},
				}, typ)
			}

			additionalMsg = fmt.Sprintf("รับออเดอร์จ้า %v จำนวน %v", splitCommands[2], quantity)
		}

	case firstNChar(command, 2) == "ล ", firstNChar(command, 3) == "ลบ ", firstNChar(command, 3) == "ลด ":
		additionalMsg = "รับทราบจ้า"
		splitCommands := strings.Split(command, " ")
		l := len(splitCommands)

		if l < 3 || l > 4 {
			additionalMsg = "เกิดข้อผิดพลาด กรุณาทำรายการใหม่"
			return
		} else {
			quantity := 0
			if l == 3 {
				quantity = -1
			} else {
				q, err := strconv.Atoi(splitCommands[3])
				if err != nil {
					additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
					return
				}

				quantity = q
			}

			typ := convertType(splitCommands[1])
			if typ == "" {
				additionalMsg = "สั่งผิด กรุณาสั่งใหม่จ้า"
				return
			}

			orderN, err := strconv.Atoi(splitCommands[2])
			if err == nil {
				allGroup[groupId].RemoveOrderByIndex(orderN, quantity, typ)
			} else {
				allGroup[groupId].RemoveOrder(OrderDetails{
					Name:     splitCommands[2],
					Quantity: quantity,
				}, typ)
			}
		}

	default:
		reply = false
		log.Printf("Unsupported message content: %T\n", command)
		return
	}
}

func convertType(typ string) string {
	switch {
	case typ == "ผ", strings.Contains(typ, "ผลไม้"):
		return Fruit
	case typ == "น", strings.Contains(typ, "น้ำ"):
		return Drink
	case typ == "ข", strings.Contains(typ, "ขนม"):
		return Snack
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

func findIndex(o []OrderDetails, order OrderDetails) int {
	for index, drinksOrder := range o {
		if drinksOrder.Name == order.Name {
			return index
		}
	}
	return -1
}

func increaseOrder(v []OrderDetails, o *[]OrderDetails) {
	a := *o
	for _, order := range v {
		index := findIndex(a, order)
		if len(a) == 0 || index == -1 {
			a = append(a, order)
		} else {
			a[index].Quantity = a[index].Quantity + order.Quantity
		}
	}
	*o = a
}

func increaseOrderByIndex(i, q int, o *[]OrderDetails) {
	a := *o
	if len(a) >= i {
		a[i].Quantity = a[i].Quantity + q
	}
	*o = a
}

func decreaseOrder(v OrderDetails, o *[]OrderDetails) {
	a := *o
	index := findIndex(a, v)
	if index != -1 {
		if v.Quantity == -1 {
			a = append(a[:index], a[index+1:]...)
		} else {
			a[index].Quantity = a[index].Quantity - v.Quantity

			if a[index].Quantity <= 0 {
				a = append(a[:index], a[index+1:]...)
			}
		}
	}
	*o = a
}

func decreaseOrderByIndex(i, q int, o *[]OrderDetails) {
	a := *o
	if len(a) >= i {
		if q == -1 {
			a = append(a[:i], a[i+1:]...)
		} else {
			a[i].Quantity = a[i].Quantity - q

			if a[i].Quantity <= 0 {
				a = append(a[:i], a[i+1:]...)
			}
		}
	}
	*o = a
}

func sumFn(e []OrderDetails) int {
	sum := 0
	for _, m := range e {
		sum = sum + m.Quantity
	}

	return sum
}

func isNotWeekend() bool {
	now := time.Now().In(bangkokTZ)
	return now.Weekday() != time.Saturday && now.Weekday() != time.Sunday
}

func pushMessagesAllGroup(bot *messaging_api.MessagingApiAPI, msg string) {
	if isNotWeekend() {
		for gid, _ := range allGroup {
			_, err := bot.PushMessage(&messaging_api.PushMessageRequest{
				To: gid,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{Text: msg},
				},
				NotificationDisabled:   true,
				CustomAggregationUnits: nil,
			}, "")

			if err != nil {
				log.Print(err)
			}
		}
	}
}

func pushMessagesOrderSummary(bot *messaging_api.MessagingApiAPI) {
	if isNotWeekend() {
		for gid, v := range allGroup {
			_, err := bot.PushMessage(&messaging_api.PushMessageRequest{
				To: gid,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{Text: fmt.Sprint(v)},
				},
				NotificationDisabled:   true,
				CustomAggregationUnits: nil,
			}, "")

			if err != nil {
				log.Print(err)
			}
		}
	}
}

package firebase

type Subscribe struct {
	FcmToken      []string `json:"fcmToken"`
	TopicName     string   `json:"topicName"`
	SubscribeType string   `json:"subscribeType"`
}

type SendNotification struct {
	FcmToken             []string     `json:"fcmToken"`
	TopicName            string       `json:"topicName"`
	FirebaseNotification Notification `json:"firebaseNotification"`
}

type Notification struct {
	NotificationHeader string `json:"notificationHeader"`
	NotificationBody   string `json:"notificationBody"`
}

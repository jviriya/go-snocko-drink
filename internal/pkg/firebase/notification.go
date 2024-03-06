package firebase

import (
	"context"
	"errors"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/common"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/pkg/firebaseadmin"
)

func GetTopicName(lang string) string {

	defaultTopicNameTH := fmt.Sprintf(common.DefaultTopic, config.Conf.State, common.LangThai)
	defaultTopicNameEN := fmt.Sprintf(common.DefaultTopic, config.Conf.State, common.LangEnglish)
	defaultTopicNameCN := fmt.Sprintf(common.DefaultTopic, config.Conf.State, common.LangChina)

	rtnLang := defaultTopicNameEN

	switch lang {
	case common.LangThai:
		rtnLang = defaultTopicNameTH
	case common.LangChina:
		rtnLang = defaultTopicNameCN
	}

	return rtnLang
}

func RegisOrUnRegisTokenToTopic(log *zerolog.Logger, fcmToken []string, topic string, subType common.SubscribeType) (config.ErrorCode, error) {
	switch subType {
	case common.FirebaseTypeSubscribe:
		_, err := firebaseadmin.FB.MessagingClient.SubscribeToTopic(context.Background(), fcmToken, topic)
		if err != nil {
			log.Error().Err(err).Msg("Firebase SubscribeToTopic")
			return config.EM.Validation.FirebaseCannotRegisTopic, nil
		}
	case common.FirebaseTypeUnSubscribe:
		_, err := firebaseadmin.FB.MessagingClient.UnsubscribeFromTopic(context.Background(), fcmToken, topic)
		if err != nil {
			log.Error().Err(err).Msg("Firebase UnSubscribeToTopic")
			return config.EM.Validation.FirebaseCannotRegisTopic, nil
		}
	default:
		return config.EM.Validation.FirebaseCannotRegisTopic, errors.New("subscribe Type Not Match")
	}

	return config.EM.Success, nil
}

func SendNotificationToUser(ctx context.Context, log *zerolog.Logger, fcmToken []string, topic, title, body string) error {
	if topic != "" {
		_, err := firebaseadmin.FB.MessagingClient.Send(ctx, message(topic, title, body))
		if err != nil {
			log.Error().Err(err).Msg("Firebase Send Err")
			return err
		}
	} else {
		if len(fcmToken) > 0 {
			resp, err := firebaseadmin.FB.MessagingClient.SendMulticast(ctx, multicastMessage(fcmToken, title, body))
			if err != nil {
				log.Error().Err(err).Msg("Firebase SendMulticast Err")
				return err
			}
			if resp.FailureCount > 0 {
				var failedTokens []string
				for idx, resp := range resp.Responses {
					if !resp.Success {
						failedTokens = append(failedTokens, fcmToken[idx])
					}
				}
				log.Error().Msgf("List of tokens that caused failures: %v\n", failedTokens)
			}
		}

	}
	return nil
}

func message(topic, title, body string) *messaging.Message {
	return &messaging.Message{
		Topic:   topic,
		Android: androidConfig(title, body),
		APNS:    apsnConfig(title, body),
		Webpush: webPushConfig(title, body),
	}
}

func multicastMessage(fcmToken []string, title, body string) *messaging.MulticastMessage {
	return &messaging.MulticastMessage{
		Tokens:  fcmToken,
		Android: androidConfig(title, body),
		APNS:    apsnConfig(title, body),
		Webpush: webPushConfig(title, body),
	}
}

func apsnConfig(title, body string) *messaging.APNSConfig {
	return &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: title,
					Body:  body,
				},
				MutableContent: true,
			},
		},
	}
}

func androidConfig(title, body string) *messaging.AndroidConfig {
	return &messaging.AndroidConfig{
		Data: map[string]string{
			"title": title,
			"body":  body,
		},
	}
}

func webPushConfig(title, body string) *messaging.WebpushConfig {
	return &messaging.WebpushConfig{
		Notification: &messaging.WebpushNotification{
			Title: title,
			Body:  body,
		},
	}
}

func notification(title, body string) *messaging.Notification {
	return &messaging.Notification{
		Title: title,
		Body:  body,
	}
}

package firebaseadmin

import (
	"context"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"go-pentor-bank/internal/clog"
	"google.golang.org/api/option"
	"os"
)

var FB fb

type fb struct {
	App             *firebase.App
	MessagingClient *messaging.Client
}

type firebaseAccess struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func InitFirebasePkg() error {
	log := clog.GetLog()

	firebaseAccessKey := os.Getenv("FIREBASE_ACCESS_KEY")

	//var accessKey firebaseAccess
	//err := json.Unmarshal([]byte(firebaseAccessKey), &accessKey)
	//if err != nil {
	//	log.Error().Err(err).Msg("json.Unmarshal")
	//	return err
	//}
	//
	//file, err := json.MarshalIndent(accessKey, "", " ")
	//if err != nil {
	//	log.Error().Err(err).Msg("json.MarshalIndent")
	//	return err
	//}
	//
	//err = os.WriteFile("accessKey.json", file, 0644)
	//if err != nil {
	//	log.Error().Err(err).Msg("ioutil.WriteFile")
	//	return err
	//}

	//opt := option.WithCredentialsFile("./accessKey.json")
	opt := option.WithCredentialsJSON([]byte(firebaseAccessKey))

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Error().Err(err).Msg("Firebase Admin Initial Error")
		return err
	}

	msgClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Firebase Admin Initial Error")
		return err
	}

	FB.App = app
	FB.MessagingClient = msgClient

	//e := os.Remove("./accessKey.json")
	//if e != nil {
	//	log.Error().Err(err).Msg("Remove File")
	//	return err
	//}

	return nil
}

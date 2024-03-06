package email

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"go-pentor-bank/internal/config"
	"html/template"
)

func SendHTMLEmail(htmlText string, subject string, toEmails, ccEmails []string) error {
	if config.Conf.State == config.StateLocal {
		return nil
	}
	sesSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(config.Conf.AwsConfig.AccessKey, config.Conf.AwsConfig.SecretKey, ""),
			Region:      aws.String(config.Conf.AwsConfig.DefaultRegion),
			MaxRetries:  aws.Int(5),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	sesClient := ses.New(sesSession)

	_, err := sesClient.SendEmail(&ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: aws.StringSlice(ccEmails),
			ToAddresses: aws.StringSlice(toEmails),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CHARSET),
					Data:    aws.String(htmlText),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CHARSET),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(Sender),
	})

	return err
}

func GenerateVerifyEmailTemplate(otp, ip, expMin, refCode string) (string, error) {
	t, err := template.ParseFiles("template/email/email_template.html")
	if err != nil {
		return "", err
	}

	var body bytes.Buffer

	err = t.Execute(&body, struct {
		Otp     string
		Ip      string
		ExpMin  string
		RefCode string
	}{
		Otp:     otp,
		Ip:      ip,
		ExpMin:  expMin,
		RefCode: refCode,
	})

	if err != nil {
		return "", err
	}
	return body.String(), err
}

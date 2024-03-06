package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

var Conf Configs

type Configs struct {
	vn            *viper.Viper
	ConfigPath    string
	State         State
	ServerSetting struct {
		RunMode              string        `mapstructure:"run_mode"`
		AllowOrigins         []string      `mapstructure:"allow_origins"`
		AccessSecret         string        `mapstructure:"access_secret"`
		RefreshSecret        string        `mapstructure:"refresh_secret"`
		TokenLifetimeMinutes time.Duration `mapstructure:"token_lifetime_minutes"`
		DiffTimestamp        time.Duration `mapstructure:"diff_timestamp"`
		PingMessage          string        `mapstructure:"ping_message"`
		Domain               string        `mapstructure:"domain"`
		MaxFileSize          int           `mapstructure:"max_file_size"`
		AllowedMimeTypes     []string      `mapstructure:"allowed_mime_types"`
		AllowedExtensions    []string      `mapstructure:"allowed_extensions"`
		Path                 struct {
			QrCode string `mapstructure:"qr_code"`
		} `mapstructure:"path"`
		AllowedMimeTypesMap  map[string]interface{}
		AllowedExtensionsMap map[string]interface{}
		RateLimit            struct {
			CommonMinute time.Duration `mapstructure:"common_minute"`
			CommonCount  int64         `mapstructure:"common_count"`
		} `mapstructure:"rate_limit"`
	} `mapstructure:"server_setting"`
	AppVersions map[string]AppVersionMap
	MongoDriver struct {
		DB struct {
			Schema             string `mapstructure:"schema" overrideEnv:"DB_SCHEMA"`
			URL                string `mapstructure:"url" overrideEnv:"DB_URL"`
			WriteMajorityCount int    `mapstructure:"write_majority_count"`
			MasterKey          string `overrideEnv:"MONGO_MASTER_KEY"`
		} `mapstructure:"db"`
		Shard struct {
			Schema             string `mapstructure:"schema" overrideEnv:"DB_SCHEMA_SHARD"`
			URL                string `mapstructure:"url" overrideEnv:"DB_URL_SHARD"`
			WriteMajorityCount int    `mapstructure:"write_majority_count"`
		} `mapstructure:"db"`
	} `mapstructure:"mongo_driver"`
	Redis struct {
		RedisPrefix     string `mapstructure:"redis_prefix"`
		PoolSize        int    `mapstructure:"pool_size"`
		MinIdleConns    int    `mapstructure:"min_idle_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		PoolTimeout     int    `mapstructure:"pool_timeout"`
		WriteTimeout    int    `mapstructure:"write_timeout"`
		ReadTimeout     int    `mapstructure:"read_timeout"`
		DialTimeout     int    `mapstructure:"dial_timeout"`
		MaxConnLifetime int    `mapstructure:"max_conn_lifetime"`
		Common          struct {
			Addr     string `mapstructure:"addr" overrideEnv:"REDIS_ADDR"`
			Port     int    `mapstructure:"port" overrideEnv:"REDIS_PRT"`
			Password string `mapstructure:"password" overrideEnv:"REDIS_PWD"`
		} `mapstructure:"common"`
	} `mapstructure:"redis"`
	SMS struct {
		Movider struct {
			Domain    string `mapstructure:"domain"`
			ApiKey    string `mapstructure:"api_key" overrideEnv:"MOVIDER_API_KEY"`
			ApiSecret string `mapstructure:"api_secret" overrideEnv:"MOVIDER_API_SECRET"`
		} `mapstructure:"movider"`
		Deecommerce struct {
			Domain       string `mapstructure:"domain"`
			SmsAccountID string `mapstructure:"sms_account_id" overrideEnv:"SMS_ACCOUNT_ID"`
			SmsSecretKey string `mapstructure:"sms_secret_key" overrideEnv:"SMS_SECRET_KEY"`
			SmsType      string `mapstructure:"sms_type" overrideEnv:"SMS_TYPE"`
			SmsSender    string `mapstructure:"sms_sender" overrideEnv:"SMS_SENDER"`
		} `mapstructure:"deecommerce"`
	} `mapstructure:"sms"`
	AwsConfig struct {
		DefaultRegion   string `mapstructure:"default_region"`
		AccessKey       string `overrideEnv:"AWS_GO_ACCESS_KEY_ID"`
		SecretKey       string `overrideEnv:"AWS_GO_SECRET_ACCESS_KEY"`
		MongoEncryption struct {
			Region string `overrideEnv:"MONGO_AWS_REGION"`
			ARN    string `overrideEnv:"MONGO_AWS_ARN"`
		}
		DynamoDB struct {
			Region string `mapstructure:"region" overrideEnv:"AWS_DYNAMODB_REGION"`
		} `mapstructure:"dynamo_db"`
		SES struct {
			Order               string `mapstructure:"order"`
			WaitingPaymentAlert string `mapstructure:"waiting_alert"`
		} `mapstructure:"ses"`
		S3 struct {
			CDN string `mapstructure:"cdn"`
		} `mapstructure:"s3"`
	} `mapstructure:"aws_config"`
	CurrencyApi struct {
		CurrencyAccessKey string `mapstructure:"currency_access_key" overrideEnv:"CURRENCY_ACCESS_KEY"`
		CurrencyApiUrl    string `mapstructure:"currency_api_url"`
	} `mapstructure:"currency_api"`
	CloudFlare struct {
		Domain string `mapstructure:"domain"`
		Images struct {
			Path struct {
				Upload    string `mapstructure:"upload"`
				ListImage string `mapstructure:"listImage"`
			} `mapstructure:"path"`
			AccountID string `overrideEnv:"CF_IMG_ACCOUNT_ID"`
			ApiToken  string `overrideEnv:"CF_IMG_API_TOKEN"`
		} `mapstructure:"images"`
		Turnstile struct {
			URL       string `mapstructure:"url"`
			SecretKey string `mapstructure:"secret_key" overrideEnv:"TURNSTILE_SECRET_KEY"`
		} `mapstructure:"turnstile"`
	} `mapstructure:"cloud_flare"`
	Image struct {
		MaxFileSize          int64    `mapstructure:"max_file_size"`
		FileBufferSize       int64    `mapstructure:"file_buffer_size"`
		AllowedMimeTypes     []string `mapstructure:"allowed_mime_types"`
		AllowedExtensions    []string `mapstructure:"allowed_extensions"`
		AllowedMimeTypesMap  map[string]interface{}
		AllowedExtensionsMap map[string]interface{}
		Path                 struct {
			Token                 string `mapstructure:"token"`
			Country               string `mapstructure:"country"`
			Currency              string `mapstructure:"currency"`
			ProfileDefault        string `mapstructure:"profile_default"`
			ProfileSupportDefault string `mapstructure:"profile_support_default"`
		} `mapstructure:"path"`
	} `mapstructure:"image"`
	ElasticSearch struct {
		CloudID      string `overrideEnv:"ELK_CLOUD_ID"`
		Username     string `overrideEnv:"ELK_USERNAME"`
		Password     string `overrideEnv:"ELK_PASSWORD"`
		TemplatesDir string `mapstructure:"templates_dir" validate:"required"`
		Index        struct {
			Order IndexTemplateES `mapstructure:"order" validate:"required"`
		} `mapstructure:"index" validate:"required"`
	} `mapstructure:"elastic_search" validate:"required"`
	GeeTestCaptchaV4 struct {
		Api               string `mapstructure:"api"`
		WebCaptchaKey     string `overrideEnv:"WEB_CAPTCHA_KEY"`
		AndroidCaptchaKey string `overrideEnv:"ANDROID_CAPTCHA_KEY"`
		IOSCaptchaKey     string `overrideEnv:"IOS_CAPTCHA_KEY"`
	} `mapstructure:"gee_test_captcha_v4"`
	Google struct {
		ReCaptchaV2 struct {
			URL       string `mapstructure:"url"`
			SecretKey string `overrideEnv:"RECAPTCHA_V2_SECRET_KEY"`
		} `mapstructure:"re_captcha_v2"`
	} `mapstructure:"google"`
	OTP struct {
		LimitationPeriod time.Duration `mapstructure:"limitation_period"`
	} `mapstructure:"otp"`
}

type AppVersionMap struct {
	Status        string
	AppVersion    string
	OSVersion     string
	AppForceUnder string
	ForceUpdate   bool
	StoreLink     string
}

type (
	IndexTemplateES struct {
		Name     string `mapstructure:"name" validate:"required"`
		Template string `mapstructure:"template" validate:"required"`
	}
	BasisDataTemplate struct {
		Value  string `mapstructure:"value"`
		Field1 Lang   `mapstructure:"field1"`
		Field2 Lang   `mapstructure:"field2"`
	}
)

func (c *Configs) InitViperWithStage(s, cfPath string) error {
	c.ConfigPath = cfPath
	c.State = parseState(s)

	name := fmt.Sprintf("config.%s", c.State)

	vn := viper.New()
	vn.SetDefault("level", "info")
	vn.SetDefault("format", "text")
	vn.SetDefault("output", "stderr")
	vn.AddConfigPath(c.ConfigPath)
	vn.SetConfigName(name)

	if err := vn.ReadInConfig(); err != nil {
		return err
	}
	c.vn = vn

	if err := c.binding(); err != nil {
		return err
	}

	readEnvironmentConfig(c)

	//vn.WatchConfig()
	//vn.OnConfigChange(func(e fsnotify.Event) {
	//	log.Println("config file changed:", e.Name)
	//	if err := c.binding(); err != nil {
	//		log.Println("binding error:", err)
	//	}
	//	//readEnvironmentConfig(c)
	//	log.Printf("config: %+v", c)
	//})

	return nil
}

func parseState(s string) State {
	switch strings.ToLower(s) {
	case "local", "l":
		return StateLocal
	case "loadtest":
		return StateLoadTest
	case "dev", "develop", "development", "d":
		return StateDEV
	case "sit", "staging", "s":
		return StateSIT
	case "uat":
		return StateUAT
	case "prod", "production", "p":
		return StateProd
	}
	return StateLocal
}

func (c *Configs) binding() error {
	if err := c.vn.Unmarshal(&c); err != nil {
		log.Println("unmarshal config error:", err)
		return err
	}

	c.Image.AllowedMimeTypesMap = make(map[string]interface{})
	for _, v := range c.Image.AllowedMimeTypes {
		c.Image.AllowedMimeTypesMap[v] = nil
	}
	c.Image.AllowedExtensionsMap = make(map[string]interface{})
	for _, v := range c.Image.AllowedExtensions {
		c.Image.AllowedExtensionsMap[v] = nil
	}
	c.ServerSetting.AllowedMimeTypesMap = make(map[string]interface{})
	for _, v := range c.ServerSetting.AllowedMimeTypes {
		c.ServerSetting.AllowedMimeTypesMap[v] = nil
	}
	c.ServerSetting.AllowedExtensionsMap = make(map[string]interface{})
	for _, v := range c.ServerSetting.AllowedExtensions {
		c.ServerSetting.AllowedExtensionsMap[v] = nil
	}

	return nil
}

func replaceEnvInConfig(body []byte) []byte {
	search := regexp.MustCompile(`\$\{([^{}]+)\}`)
	replacedBody := search.ReplaceAllFunc(body, func(b []byte) []byte {
		group1 := search.ReplaceAllString(string(b), `$1`)

		envValue := os.Getenv(group1)
		if len(envValue) > 0 {
			return []byte(envValue)
		}
		return []byte("")
	})

	return replacedBody
}

// Structure
func customFunc(fl validator.FieldLevel) bool {

	if fl.Field().String() == "invalid" {
		return false
	}

	return true
}

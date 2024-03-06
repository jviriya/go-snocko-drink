package elasticsearch

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
)

var EsClient *elasticsearch.Client

func NewElasticSearchCon(log *zerolog.Logger) error {
	log.Info().Msg("Connecting to Elasticsearh..")

	var err error
	EsClient, err = elasticsearch.NewClient(elasticsearch.Config{
		CloudID:  config.Conf.ElasticSearch.CloudID,
		Username: config.Conf.ElasticSearch.Username,
		Password: config.Conf.ElasticSearch.Password,

		// Retry on 429 TooManyRequests statuses
		RetryOnStatus: []int{502, 503, 504, 429},

		// Retry up to 5 attempts
		//
		MaxRetries: 5,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error creating the client")
		return err
	}

	// Get cluster info
	res, err := EsClient.Info()
	if err != nil {
		log.Error().Err(err).Msg("Error getting response")
		return err
	}
	defer res.Body.Close()

	// Check response status
	if res.IsError() {
		err := errors.New(res.String())
		log.Error().Err(err).Msg(res.String())
		return err
	}
	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error().Err(err).Msg("Error parsing the response body")
		return err
	}

	log.Info().Msgf("client: %s", elasticsearch.Version)
	log.Info().Msgf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Info().Msg("Connecting to Elasticsearh success!!")

	return nil
}

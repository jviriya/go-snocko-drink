package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/infra/elasticsearch"
	"strings"
	"text/template"
	"time"
)

func GetElkIndex(indexName string, smt, emt time.Time, log *zerolog.Logger) []string {
	var indexes []string
	for smt.Before(emt) || smt.Equal(emt) {
		keyDate := smt.In(config.TimeZone.Bangkok).Format("200601")
		if IsIndexExist(log, indexName, keyDate) {
			index := fmt.Sprintf("%v_%v", indexName, keyDate)
			indexes = append(indexes, index)
		}
		smt = smt.AddDate(0, 1, 0)
	}
	return indexes
}

func IsIndexExist(log *zerolog.Logger, indexName, keyDate string) bool {
	// check index exist
	indexKey := indexName + "_" + keyDate

	res, err := elasticsearch.EsClient.Indices.Exists([]string{indexKey})
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		log.Error().Err(err).Msg("Indices.Exists got error")
		return false
	}

	if res.IsError() {
		if res.StatusCode == 404 {
			err := CreateIndex(indexName, keyDate, log)
			if err != nil {
				log.Error().Err(err).Msgf("Cannot create index: %s")
				return false
			}
		} else {
			log.Error().Err(err).Msg("index : " + indexName + " + error : " + res.String())
			return false
		}
	}
	return true
}

func CreateIndex(indexName, keyDate string, log *zerolog.Logger) error {
	ms, err := GetEsMappingString(indexName, log)
	if err != nil {
		return err
	}
	indexKey := indexName + "_" + keyDate
	res, err := elasticsearch.EsClient.Indices.Create(
		indexKey,
		elasticsearch.EsClient.Indices.Create.WithBody(strings.NewReader(ms)))
	if err != nil {
		log.Error().Err(err).Msgf("Cannot create index: %s", err)
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Error().Err(err).Msgf("Cannot create index: %s", res)
		return errors.New(res.String())
	}

	return nil
}

func GetEsMappingString(templateName string, log *zerolog.Logger) (string, error) {
	template, err := template.ParseFiles(fmt.Sprintf("%s/%s.json", config.Conf.ElasticSearch.TemplatesDir, templateName))
	if err != nil {
		log.Error().Err(err).Msg("[GetEsMappingString] parse template error")
		return "", err
	}

	buff := bytes.Buffer{}
	if err := template.Execute(&buff, ""); err != nil {
		log.Error().Err(err).Msg("[GetEsMappingString] binding template error")
		return "", err
	}

	return buff.String(), nil
}

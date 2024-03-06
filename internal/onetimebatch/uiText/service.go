package uiText

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"go-pentor-bank/internal/clog"
	"go-pentor-bank/internal/config"
	"go-pentor-bank/internal/infra/mongodb"
	"go-pentor-bank/internal/repository/commondb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
)

func Run(platform commondb.UiTextPlatform, langs []string) {
	log := clog.GetLog()
	repo := commondb.NewRepositoryV2(mongodb.MongoDBCon.SecureClient, mongodb.MongoDBCon.SecureClientShard)
	toInsertMap := make(map[string]config.Lang) // [lang][code]

	_, err := repo.DeleteManyUiText(context.Background(), platform)
	if err != nil {
		log.Error().Err(err).Msg("DeleteManyUiText got err")
		return
	}

	for _, lang := range langs {
		jsonFile, err := os.Open(fmt.Sprintf("assets/uiText/%v/%v.json", platform, lang))
		if err != nil {
			log.Error().Err(err).Msg("os.Open json file got err")
			return
		}

		byteJson, err := io.ReadAll(jsonFile)
		if err != nil {
			log.Error().Err(err).Msg(" io.ReadAll got err")
			return
		}

		var data []JsonDataFile
		err = sonic.Unmarshal(byteJson, &data)
		if err != nil {
			log.Error().Err(err).Msg("sonic.Unmarshal got err")
			return
		}

		for _, row := range data {
			oldLang := toInsertMap[row.Code]
			oldLang.SetByLang(lang, row.Content)
			toInsertMap[row.Code] = oldLang
		}

		jsonFile.Close()
	}

	var writes []mongo.WriteModel
	for code, content := range toInsertMap {
		filter := bson.M{
			"code":     code,
			"platform": platform,
		}
		updater := bson.M{
			"$set": bson.M{
				"lang": content,
			},
		}
		write := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updater).SetUpsert(true)
		writes = append(writes, write)
	}

	_, err = repo.BulkUiText(context.Background(), writes)
	if err != nil {
		log.Error().Err(err).Msg("BulkUiText got err")
		return
	}
}

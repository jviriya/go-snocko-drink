package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

var encryptedMapField map[string]map[string]Field

type Field struct {
	Path     string    `json:"path"`
	BsonType string    `json:"bsonType"`
	KeyId    string    `json:"keyId"`
	Queries  []Queries `json:"queries"`
}

type Queries struct {
	QueryType string `json:"queryType"`
}

func readCollectionMapField() error {
	encryptedMapField = make(map[string]map[string]Field)
	mapDataKey = make(map[string]primitive.Binary)

	dirPath := "configs/mongodb/encryptedMapField"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(fmt.Sprintf("%v/%v", dirPath, file.Name()))
		if err != nil {
			return err
		}

		var jsonVal []Field
		err = sonic.Unmarshal(data, &jsonVal)
		if err != nil {
			return err
		}

		mapVal := make(map[string]Field)
		for i, field := range jsonVal {
			mapVal[field.Path] = jsonVal[i]
		}

		collectionName := strings.Split(file.Name(), ".")
		encryptedMapField[collectionName[0]] = mapVal
	}

	return nil
}

type EncryptionField struct {
	clientEnc      *mongo.ClientEncryption
	collectionName string
}

//func NewEncryptionField(ctx context.Context, secureClient *mongo.client, colName string) (*EncryptionField, func(), error) {
//	provider := common.AWS
//	kmsProviders := map[string]map[string]interface{}{
//		provider: {
//			"accessKeyId":     config.Conf.AwsConfig.KMS.AccessKeyID,
//			"secretAccessKey": config.Conf.AwsConfig.KMS.SecretAccessKey,
//		},
//	}
//
//	keyVaultDb := config.Conf.MongoDriver.DB.KeySpace
//	keyVaultColl := config.Conf.MongoDriver.DB.KeyVault
//
//	keyVaultNamespace := keyVaultDb + "." + keyVaultColl
//	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace(keyVaultNamespace).
//		SetKmsProviders(kmsProviders)
//	clientEnc, err := mongo.NewClientEncryption(secureClient, clientEncryptionOpts)
//	if err != nil {
//		return &EncryptionField{}, func() {}, err
//	}
//
//	closeFunc := func() {
//		_ = clientEnc.Close(ctx)
//	}
//
//	return &EncryptionField{
//		clientEnc:      clientEnc,
//		collectionName: colName,
//	}, closeFunc, nil
//}

func (ef *EncryptionField) GetInsertEncryptedField(fieldName string, value interface{}) (primitive.Binary, error) {
	rawValueType, rawValueData, err := bson.MarshalValue(value)
	if err != nil {
		return primitive.Binary{}, err
	}

	var algorithm string
	var encryptionOpts *options.EncryptOptions
	if encryptedMapField[ef.collectionName][fieldName].Queries != nil {
		algorithm = "Indexed"
		encryptionOpts = options.Encrypt().
			SetAlgorithm(algorithm).
			SetKeyID(mapDataKey[encryptedMapField[ef.collectionName][fieldName].KeyId]).
			SetContentionFactor(1)
	} else {
		algorithm = "Unindexed"
		encryptionOpts = options.Encrypt().
			SetAlgorithm(algorithm).
			SetKeyID(mapDataKey[encryptedMapField[ef.collectionName][fieldName].KeyId])
	}

	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}
	encryptedField, err := ef.clientEnc.Encrypt(
		context.TODO(),
		rawValue,
		encryptionOpts)
	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, err
}

func (ef *EncryptionField) GetFindEncryptedField(fieldName string, value interface{}) (primitive.Binary, error) {
	if encryptedMapField[ef.collectionName][fieldName].Queries == nil {
		return primitive.Binary{}, errors.New("this field is not indexed")
	}

	rawValueType, rawValueData, err := bson.MarshalValue(value)
	if err != nil {
		return primitive.Binary{}, err
	}

	encryptionOpts := options.Encrypt().
		SetAlgorithm("Indexed").
		SetKeyID(mapDataKey[encryptedMapField[ef.collectionName][fieldName].KeyId]).
		SetQueryType(encryptedMapField[ef.collectionName][fieldName].Queries[0].QueryType).
		SetContentionFactor(1)

	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}
	encryptedField, err := ef.clientEnc.Encrypt(
		context.TODO(),
		rawValue,
		encryptionOpts)
	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, err
}

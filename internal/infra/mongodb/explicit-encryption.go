package mongodb

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-pentor-bank/internal/config"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type dataKey struct {
	regularClient *mongo.Client
	clientEnc     *mongo.ClientEncryption
	keyVaultDb    string
	keyVaultColl  string
	masterKey     interface{}
	provider      string
}

func (dk dataKey) retrieve(keyAltName string) error {
	var foundDoc1 bson.M
	err := dk.regularClient.Database(dk.keyVaultDb).Collection(dk.keyVaultColl).FindOne(context.Background(), bson.D{{"keyAltNames", keyAltName}}).Decode(&foundDoc1)
	if err == mongo.ErrNoDocuments {
		switch dk.provider {
		case "local":
			dataKeyOpts1 := options.DataKey().
				//SetMasterKey(dk.masterKey).
				SetKeyAltNames([]string{keyAltName})
			dataKeyID1, err := dk.clientEnc.CreateDataKey(context.Background(), dk.provider, dataKeyOpts1)
			if err != nil {
				return fmt.Errorf("create data key error %v", err)
			}
			mapDataKey[keyAltName] = dataKeyID1
		case "aws":
			dataKeyOpts := options.DataKey().
				SetMasterKey(dk.masterKey).
				SetKeyAltNames([]string{keyAltName})

			dataKeyID, err := dk.clientEnc.CreateDataKey(context.Background(), dk.provider, dataKeyOpts)
			if err != nil {
				return fmt.Errorf("create data key error %v", err)
			}
			mapDataKey[keyAltName] = dataKeyID
		}
		return nil
	} else if err != nil {
		return err
	}
	mapDataKey[keyAltName] = foundDoc1["_id"].(primitive.Binary)
	return nil
}

const (
	keyVaultColl = "__keyVault"
)

var (
	mapDataKey map[string]primitive.Binary
)

type newMongoExpliEncConn struct {
	DbUrl     string
	Schema    string
	MasterKey string
}

func newMongoExplicitEncryptionConn(log *zerolog.Logger, connReq newMongoExpliEncConn) (SecureClient, error) {
	log.Info().Msg("Connecting to Secure MongoDB..")

	regularClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connReq.DbUrl))
	if err != nil {
		return SecureClient{}, fmt.Errorf("Connect error for regular client: %v", err)
	}

	keyVaultDb := connReq.Schema
	keyVaultNamespace := keyVaultDb + "." + keyVaultColl
	keyVaultIndex := mongo.IndexModel{
		Keys: bson.D{{"keyAltNames", 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{
				{"keyAltNames", bson.D{
					{"$exists", true},
				}},
			}),
	}

	// Drop the Key Vault Collection in case you created this collection
	// in a previous run of this application.
	//if err = regularClient.Database(keyVaultDb).Collection(keyVaultColl).Drop(context.Background()); err != nil {
	//	log.Fatal().Msgf("Collection.Drop error: %v", err)
	//}

	// create indexes
	_, err = regularClient.Database(keyVaultDb).Collection(keyVaultColl).Indexes().CreateOne(context.Background(), keyVaultIndex)
	if err != nil {
		log.Error().Err(err).Msg("Indexing keyVaultColl got err")
		return SecureClient{}, err
	}

	provider := "aws"
	kmsProviders := make(map[string]map[string]interface{})
	var masterKey interface{}
	switch provider {
	case "local":
		key, err := base64.StdEncoding.DecodeString(connReq.MasterKey)
		if err != nil {
			log.Fatal().Err(err).Msg("base64.StdEncoding.DecodeString")
		}

		kmsProviders = map[string]map[string]interface{}{"local": {"key": key}}
		masterKey = kmsProviders
	case "aws":
		kmsProviders = map[string]map[string]interface{}{
			provider: {
				"accessKeyId":     config.Conf.AwsConfig.AccessKey,
				"secretAccessKey": config.Conf.AwsConfig.SecretKey,
			},
		}
		masterKey = map[string]interface{}{
			"key":    config.Conf.AwsConfig.MongoEncryption.ARN,
			"region": config.Conf.AwsConfig.MongoEncryption.Region,
		}
	}

	clientEncryptionOpts := options.ClientEncryption().SetKeyVaultNamespace(keyVaultNamespace).
		SetKmsProviders(kmsProviders)
	clientEnc, err := mongo.NewClientEncryption(regularClient, clientEncryptionOpts)
	if err != nil {
		return SecureClient{}, fmt.Errorf("NewClientEncryption error %v", err)
	}

	datakey := dataKey{
		regularClient: regularClient,
		clientEnc:     clientEnc,
		keyVaultDb:    keyVaultDb,
		keyVaultColl:  keyVaultColl,
		masterKey:     masterKey,
		provider:      provider,
	}

	err = readCollectionMapField()
	if err != nil {
		return SecureClient{}, err
	}

	encryptedFieldsMap := bson.M{}
	for collectionName, field := range encryptedMapField {
		keyChain := []primitive.Binary{}
		for _, v := range field {
			keyID := v.KeyId
			if _, ok := mapDataKey[keyID]; !ok {
				err = datakey.retrieve(keyID)
				if err != nil {
					log.Error().Err(err).Msg("retrieve got err")
					return SecureClient{}, err
				}
			}

			keyChain = append(keyChain, mapDataKey[keyID])
		}

		var bdoc []bson.M
		count := 0
		for _, f := range field {
			querries := []bson.M{}
			for _, q := range f.Queries {
				querries = append(querries, bson.M{"queryType": q.QueryType})
			}
			bdoc = append(bdoc, bson.M{"path": f.Path, "bsonType": f.BsonType, "keyId": keyChain[count], "queries": querries})
			count++
		}

		coll := keyVaultDb + "." + collectionName
		encryptedFieldsMap[coll] = bson.M{"fields": bdoc}
	}

	//extraOptions := map[string]interface{}{
	//	"cryptSharedLibPath": "",
	//}

	autoEncryptionOpts := options.AutoEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace(keyVaultNamespace).
		//SetEncryptedFieldsMap(encryptedFieldsMap). // use with queryable encryption
		//SetExtraOptions(extraOptions). // use with queryable encryption
		SetBypassQueryAnalysis(true)

	ops := options.Client().ApplyURI(connReq.DbUrl)
	ops.SetReadPreference(readpref.SecondaryPreferred())
	ops.SetMaxPoolSize(20000)
	ops.SetConnectTimeout(5 * time.Second)
	ops.SetMaxConnIdleTime(10 * time.Second)
	ops.SetAutoEncryptionOptions(autoEncryptionOpts)
	secureClient, err := mongo.Connect(context.Background(), ops)
	if err != nil {
		log.Error().Err(err).Msg("mongo.Connect got err")
		return SecureClient{}, err
	}

	log.Info().Msg("Connecting to Secured MongoDB success!!")

	return SecureClient{
		MapDataKey:   mapDataKey,
		SecureClient: secureClient,
		ClientEnc:    clientEnc,
		Schema:       connReq.Schema,
	}, nil

	//colorValueType, colorRawValueData, err := bson.MarshalValue("55")
	//if err != nil {
	//	panic(err)
	//}
	//colorIdRawValue := bson.RawValue{Label: colorValueType, ExpectedValue: colorRawValueData}
	//colorEncryptionOpts := options.Encrypt().
	//	SetAlgorithm("Indexed").
	//	SetKeyID(mapDataKey["dataKey1"]).
	//	SetContentionFactor(1)
	//colorEncryptedField, err := clientEnc.Encrypt(
	//	context.Background(),
	//	colorIdRawValue,
	//	colorEncryptionOpts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//heightValueType, heightRawValueData, err := bson.MarshalValue("55")
	//if err != nil {
	//	panic(err)
	//}
	//heightIdRawValue := bson.RawValue{Label: heightValueType, ExpectedValue: heightRawValueData}
	//heightEncryptionOpts := options.Encrypt().
	//	SetAlgorithm("Indexed").
	//	SetKeyID(mapDataKey["dataKey2"]).
	//	SetContentionFactor(1)
	//heightEncryptedField, err := clientEnc.Encrypt(
	//	context.Background(),
	//	heightIdRawValue,
	//	heightEncryptionOpts)
	//if err != nil {
	//	panic(err)
	//}

	//raceRawValueType, raceRawValueData, err := bson.MarshalValue("5")
	//if err != nil {
	//	panic(err)
	//}
	//raceRawValue := bson.RawValue{Label: raceRawValueType, ExpectedValue: raceRawValueData}
	//raceEncryptionOpts := options.Encrypt().
	//	SetAlgorithm("Unindexed").
	//	SetKeyID(mapDataKey["dataKey3"])
	//raceEncryptedField, err := clientEnc.Encrypt(
	//	context.Background(),
	//	raceRawValue,
	//	raceEncryptionOpts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = secureClient.Database(keyVaultDb).Collection("lead_test_mongoEncrypt").InsertOne(
	//	context.Background(),
	//	bson.D{{"firstname", colorEncryptedField}, {"lastname", heightEncryptedField}, {"height", "000"}})
	//if err != nil {
	//	panic(err)
	//}

	//valueType, rawValueData, err := bson.MarshalValue("11")
	//if err != nil {
	//	panic(err)
	//}
	//rawValue := bson.RawValue{Label: valueType, ExpectedValue: rawValueData}
	//encryptionOpts := options.Encrypt().
	//	SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
	//	SetKeyID(mapDataKey["dataKey1"])
	//encryptedField, err := clientEnc.Encrypt(
	//	context.Background(),
	//	rawValue,
	//	encryptionOpts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//_, err = secureClient.Database(keyVaultDb).Collection("lead_test_explicit_mongoEncrypt").InsertOne(
	//	context.Background(),
	//	bson.D{{"firstName", "Jon"}, {"patientId", encryptedField}})
	//if err != nil {
	//	panic(err)
	//}
	//var resultSecure bson.M
	//err = secureClient.Database(keyVaultDb).Collection("lead_test_explicit_mongoEncrypt").FindOne(context.Background(), bson.D{{"patientId", encryptedField}}).Decode(&resultSecure)
	//if err != nil {
	//	panic(err)
	//}
	//outputSecure, err := json.MarshalIndent(resultSecure, "", "    ")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%s\n", outputSecure)
}

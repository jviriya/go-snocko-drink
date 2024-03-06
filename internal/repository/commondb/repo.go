package commondb

import (
	"context"
	"go-pentor-bank/internal/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB Collection name
const (
	cltTest   = "test"
	cltUiText = "ui_texts"
)

type Repossitory struct {
	DB       *mongo.Database
	DBClient *mongo.Client
}

type SecureRepository struct {
	Secure      secureDB
	SecureShard secureDB
}

type secureDB struct {
	DB        *mongo.Database
	DBClient  *mongo.Client
	DataKey   map[string]primitive.Binary
	ClientEnc *mongo.ClientEncryption
}

func NewRepositoryV2(secureClient, secureCliShard mongodb.SecureClient) *SecureRepository {
	return &SecureRepository{
		Secure: secureDB{
			DB:        secureClient.SecureClient.Database(secureClient.Schema),
			DBClient:  secureClient.SecureClient,
			DataKey:   secureClient.MapDataKey,
			ClientEnc: secureClient.ClientEnc,
		},
		SecureShard: secureDB{
			DB:        secureCliShard.SecureClient.Database(secureCliShard.Schema),
			DBClient:  secureCliShard.SecureClient,
			DataKey:   secureCliShard.MapDataKey,
			ClientEnc: secureCliShard.ClientEnc,
		},
	}
}

func (rp *SecureRepository) DeterministicEncryptionReplica(data interface{}) (primitive.Binary, error) {
	valueType, rawValueData, err := bson.MarshalValue(data)
	if err != nil {
		return primitive.Binary{}, err
	}
	rawValue := bson.RawValue{Type: valueType, Value: rawValueData}
	encryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
		SetKeyID(rp.Secure.DataKey["dataKey1"])
	encryptedField, err := rp.Secure.ClientEnc.Encrypt(
		context.TODO(),
		rawValue,
		encryptionOpts)
	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, nil
}

func (rp *SecureRepository) DeterministicEncryptionShard(data interface{}) (primitive.Binary, error) {
	valueType, rawValueData, err := bson.MarshalValue(data)
	if err != nil {
		return primitive.Binary{}, err
	}
	rawValue := bson.RawValue{Type: valueType, Value: rawValueData}
	encryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
		SetKeyID(rp.SecureShard.DataKey["dataKey1"])
	encryptedField, err := rp.SecureShard.ClientEnc.Encrypt(
		context.TODO(),
		rawValue,
		encryptionOpts)
	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, nil
}

func (rp *SecureRepository) randomEncryption(data interface{}) (primitive.Binary, error) {
	valueType, rawValueData, err := bson.MarshalValue(data)
	if err != nil {
		return primitive.Binary{}, err
	}
	rawValue := bson.RawValue{Type: valueType, Value: rawValueData}
	encryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Random").
		SetKeyID(rp.Secure.DataKey["dataKey1"])
	encryptedField, err := rp.Secure.ClientEnc.Encrypt(
		context.TODO(),
		rawValue,
		encryptionOpts)
	if err != nil {
		return primitive.Binary{}, err
	}

	return encryptedField, nil
}

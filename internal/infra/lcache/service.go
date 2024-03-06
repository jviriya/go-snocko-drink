package lcache

import (
	"github.com/bytedance/sonic"
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

func Init() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

func GetDefaultTime() time.Duration {
	return cache.DefaultExpiration
}

func Set(key string, data interface{}, duration time.Duration) error {
	var dataByte []byte
	if dstr, valid := data.(string); valid {
		dataByte = []byte(dstr)
	} else {
		var err error
		dataByte, err = sonic.Marshal(data)
		if err != nil {
			return err
		}
	}

	Cache.Set(key, string(dataByte), duration)

	return nil
}

func Parse(key string, data interface{}) (bool, error) {
	dataInf, found := Cache.Get(key)
	if found {
		switch v := dataInf.(type) {
		case string:
			err := sonic.Unmarshal([]byte(v), data)
			if err != nil {
				return false, err
			}
		case []byte:
			err := sonic.Unmarshal(v, data)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}

	return false, nil
}

func Delete(key string) {
	Cache.Delete(key)
}

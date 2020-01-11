package datasources

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/codelieche/microservice/usercenter/datamodels"
)

func TestGetRedisClient(t *testing.T) {
	client := GetRedisClient()

	log.Println(client)

	key := "user_permissions_1"
	value := datamodels.Permission{
		Name:  "权限123",
		AppID: 1,
		Code:  "code_01",
	}
	data, _ := json.Marshal(value)
	if err := client.Set(key, data, time.Second*30).Err(); err != nil {
		t.Error(err)
	} else {
		log.Println(client.Get(key).Result())
	}

}

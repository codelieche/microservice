package datasources

import (
	"log"
	"testing"
)

func TestGetMySQLDB(t *testing.T) {
	db := GetMySQLDB()
	if db != nil {
		log.Println("获取DB成功")
	} else {
		t.Error("获取db出错")
	}
}

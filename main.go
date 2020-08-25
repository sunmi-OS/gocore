package main

import "C"
import (
	"bytes"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

var str = `
[server]
mode = "development"
port = 8080
debug = true
http_timeout = 20000
jwt_timeout = 8640000
log_path = "./log/tix-log.log"
maps = "./map_data.json"
order_timeout = 600

[development]
baseURL = "http://dev.xxxxx.com:2121/partner-service/"
clientID = "xxxx"
schedule_expired_before = 900
promo_activity_id = "ttttttt"

[production]
baseURL = "http://prod.xxxxx.com:2121/partner-service/"
clientID = "yyyy"
schedule_expired_before = 900
promo_activity_id = "kkkkkkkkkk"`

func main() {

	fmt.Println("gocore")

	V := viper.New()
	V.SetConfigType("toml")
	err := V.MergeConfig(bytes.NewBuffer([]byte(str)))
	if err != nil {
		print(err)
	}

	var tmp interface{}

	if _, err := toml.Decode(str, &tmp); err != nil {

		log.Fatalf("Error decoding TOML: %s", err)
	}

	s := V.GetString("server.mode")

	fmt.Println(s)

	// use example
	//hook.ShowExample()
}

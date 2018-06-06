package main

import (
	"errors"

	"go.uber.org/zap"

	"github.com/sunmi-OS/gocore/log"
	"github.com/sunmi-OS/gocore/viper"
)

func main() {

	viper.NewConfig("config", "conf")

	log.InitLogger("example-log")
	log.Sugar.Debugw("example-log:debug")
	log.Sugar.Infow("example-log:info", zap.String("type", "log"))
	log.Sugar.Errorw("example-log:err", zap.Error(errors.New("IS ERROR")))

}

package main

import (
	"github.com/sunmi-OS/gocore/log"
	"go.uber.org/zap"
	"errors"
)

func main() {

	log.InitLogger("example-log", true)
	log.Sugar.Debugw("example-log:debug")
	log.Sugar.Infow("example-log:info", zap.String("type", "log"))
	log.Sugar.Errorw("example-log:err", zap.Error(errors.New("IS ERROR")))

}

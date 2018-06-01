package main

import (
	"errors"

	"go.uber.org/zap"

	"github.com/sunmi-OS/gocore/log"
)

func main() {

	log.InitLogger("example-log", true)
	log.Sugar.Debugw("example-log:debug")
	log.Sugar.Infow("example-log:info", zap.String("type", "log"))
	log.Sugar.Errorw("example-log:err", zap.Error(errors.New("IS ERROR")))

}

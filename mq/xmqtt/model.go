package xmqtt

import "errors"

const (
	QosAtMostOne  QosType = 0
	QosAtLeastOne QosType = 1
	QosOnlyOne    QosType = 2
)

var (
	ErrNotExists   = errors.New("mqtt not exists")
	ErrLostConnect = errors.New("mqtt connection lost")
)

type QosType byte

type Config struct {
	Host           string // host地址
	Port           int    // TCP 端口
	ClientId       string
	Uname          string
	Password       string
	KeepAlive      int // 单位秒
	IsCleanSession bool
}

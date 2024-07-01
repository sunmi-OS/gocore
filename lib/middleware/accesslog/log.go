package accesslog

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Define a custom struct that implements the gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// Override the Write([]byte) (int, error) method
func (w responseWriter) Write(b []byte) (int, error) {
	// Write a copy of the data to a bytes.buffer to use for getting the body
	w.b.Write(b)
	return w.ResponseWriter.Write(b)
}

// AccessLog middleware for accesslog
func AccessLog() gin.HandlerFunc {
	// init zap
	logFileName := utils.GetAccesslogPath()
	logger := initZap(logFileName)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return func(c *gin.Context) {
		start := time.Now()
		requestDate := start.Format(time.RFC3339)
		body := ""
		// The body is not recorded in the file upload scenario
		contentType := c.GetHeader("content-type")
		if !strings.Contains(contentType, "multipart/form-data") {
			b, err := c.GetRawData()
			if err != nil {
				body = "failed to get request body"
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
				body = string(b)
			}
		}
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery
		queryPath := path
		if rawQuery != "" {
			queryPath += "?" + rawQuery
		}
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		_, clientPort, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
		if err != nil {
			clientPort = ""
		}
		defer func() {
			var responseCode interface{}
			var responseMsg interface{}
			m := make(map[string]interface{})
			err = json.Unmarshal(writer.b.Bytes(), &m)
			if err == nil {
				responseCode = m["code"]
				responseMsg = m["msg"]
			}
			fields := []zapcore.Field{
				zap.String("r_time", requestDate),
				zap.Int64("cost", time.Since(start).Milliseconds()),
				zap.String("c_ip", c.ClientIP()),
				zap.String("c_port", clientPort),
				zap.String("c_f_ip", c.GetHeader("x-forwarded-for")),
				zap.String("schema", scheme),
				zap.String("r_host", c.Request.Host),
				zap.String("r_method", c.Request.Method),
				zap.String("r_q_path", queryPath),
				zap.String("r_path", path),
				zap.Any("r_header", c.Request.Header),
				zap.String("r_body", body),
				zap.Int("s_status", c.Writer.Status()),
				zap.Any("s_header", c.Writer.Header()),
				zap.Any("s_code", responseCode),
				zap.Any("s_msg", responseMsg),
				zap.String("s_body", writer.b.String()),
			}
			logger.Info("accesslog", fields...)
		}()
		c.Next()
	}
}

// init zap
// fileName log path ./access_log
func initZap(fileName string) *zap.Logger {
	config := zapcore.EncoderConfig{
		MessageKey: "",
		LevelKey:   "",
		TimeKey:    "",
		CallerKey:  "",
	}
	// io.Writer 使用 lumberjack
	infoWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    256, // 最大体积，单位M，超过则切割
		MaxBackups: 4,   // 最大文件保留数，超过则删除最老的日志文件
		MaxAge:     30,  // 最长保存时间30天
		Compress:   true,
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), zap.InfoLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.InfoLevel))
}

package echo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sunmi-OS/gocore/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// AccessLog middleware for accesslog
func AccessLog() echo.MiddlewareFunc {
	// init zap
	logFilePath := utils.GetAccesslogPath()
	logger := initZap(logFilePath)
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			requestDate := start.Format(time.RFC3339)
			body := ""
			b, err := io.ReadAll(c.Request().Body)
			if err != nil {
				body = "failed to get request body"
			} else {
				c.Request().Body = io.NopCloser(bytes.NewBuffer(b))
				body = string(b)
			}
			path := c.Path()
			rawQuery := c.QueryString()
			queryPath := path
			if rawQuery != "" {
				queryPath += "?" + rawQuery
			}

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer
			defer func() {
				var responseCode interface{}
				var responseMsg interface{}
				m := make(map[string]interface{})
				err = json.Unmarshal(resBody.Bytes(), &m)
				if err == nil {
					responseCode = m["code"]
					responseMsg = m["msg"]
				}

				fields := []zapcore.Field{
					zap.String("r_time", requestDate),
					zap.Int64("cost", time.Since(start).Milliseconds()),
					zap.String("c_ip", c.RealIP()),
					zap.String("c_f_ip", c.Request().Header.Get("x-forwarded-for")),
					zap.String("schema", c.Scheme()),
					zap.String("r_host", c.Request().Host),
					zap.String("r_method", c.Request().Method),
					zap.String("r_q_path", queryPath),
					zap.String("r_path", path),
					zap.Any("r_header", c.Request().Header),
					zap.String("r_body", body),
					zap.Int("s_status", c.Response().Status),
					zap.Any("s_header", c.Response().Header()),
					zap.Any("s_code", responseCode),
					zap.Any("s_msg", responseMsg),
					zap.String("s_body", strings.Trim(resBody.String(), "\n")),
				}
				logger.Info("accesslog", fields...)
			}()
			return next(c)
		}
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
	// 获取io.Writer的实现
	infoWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1024, // 最大体积，单位M，超过则切割
		MaxBackups: 5,    // 最大文件保留数，超过则删除最老的日志文件
		MaxAge:     30,   // 最长保存时间30天
		Compress:   true, // 是否压缩
	}
	// 实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), zap.InfoLevel), // 将info及以下写入logPath，NewConsoleEncoder 是非结构化输出
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.InfoLevel))
}

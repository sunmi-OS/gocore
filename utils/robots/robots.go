package robots

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/sunmi-OS/gocore/v2/utils"
	http_request "github.com/sunmi-OS/gocore/v2/utils/http-request"
)

type Robot struct {
	client *http_request.HttpClient
}

// NewWithUrl 使用url初始化
func NewWithUrl(webhookUrl string) *Robot {
	c := http_request.New()
	c.Client.SetBaseURL(webhookUrl)
	c.Client.OnAfterResponse(http_request.MustCode200)
	c.SetLog(http_request.NewGocoreLog())

	return &Robot{
		client: c,
	}
}

// NewWithClient 使用已有的client初始化
func NewWithClient(client *http_request.HttpClient) *Robot {
	return &Robot{
		client: client,
	}
}

// SendMarkdownKVs
// 该方法会把MarkdownKV，展示为key: value形式
func (r *Robot) SendMarkdownKVs(ctx context.Context, title string, kvs []MarkdownKV) error {
	msg := MarkdownMsg{
		Msgtype: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: title,
			Text:  BuildContent(kvs),
		},
	}
	_, err := r.client.Client.R().SetContext(ctx).SetBody(msg).Post("")
	return err
}

func (r *Robot) SendMarkdown(ctx context.Context, title, content string) error {
	msg := MarkdownMsg{
		Msgtype: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: title,
			Text:  content,
		},
	}
	_, err := r.client.Client.R().SetContext(ctx).SetBody(msg).Post("")
	return err
}

func (r *Robot) SendTextMsg(ctx context.Context, content string) error {
	msg := TextMsg{
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}
	_, err := r.client.Client.R().SetContext(ctx).SetBody(msg).Post("")
	return err
}

// 下面是一些通用函数

// GeneralHeader
// 通用kv，用于标识发送时间和来源
func GeneralHeader() []MarkdownKV {
	return []MarkdownKV{
		{Key: "时间", Value: time.Now().Format(utils.TimeFormat)},
		{Key: "环境", Value: utils.GetRunTime()},
		{Key: "机器", Value: utils.GetHostname()},
	}
}

// FormatError
// 将error转换成文字，并且标红
func FormatError(err error) string {
	if err == nil {
		return "无错误"
	}
	return fmt.Sprintf(`<font color="#BF4747">%s</font>`, err)
}

func BuildContent(kvs []MarkdownKV) string {
	buf := &bytes.Buffer{}
	for _, kv := range kvs {
		buf.WriteString(fmt.Sprintf("**<font color=\"#267D22\">%s</font>**: %v %v", kv.Key, kv.KeySuffix, kv.Value))
		buf.WriteString("\n\n")
	}
	str := buf.String()
	if str == "" {
		str = "无内容"
	}
	return str
}

type MarkdownKV struct {
	Key       string
	KeySuffix string
	Value     interface{}
}

type TextMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type MarkdownMsg struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

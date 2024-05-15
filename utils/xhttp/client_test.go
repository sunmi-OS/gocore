package xhttp

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/go-pay/bm"
)

type HttpGet struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var ctx = context.Background()

func TestHttpGet(t *testing.T) {
	client := NewClient()
	// test
	_, bs, err := client.Req().Get("http://www.baidu.com").EndBytes(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(bs))

	//rsp := new(HttpGet)
	//_, err = client.Type(TypeJSON).Get("http://api.igoogle.ink/app/v1/ping").EndStruct(ctx, rsp)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Println(rsp)
}

func TestHttpUploadFile(t *testing.T) {
	fileContent, err := os.ReadFile("logo.png")
	if err != nil {
		log.Println(err)
		return
	}

	bmm := make(bm.BodyMap)
	bmm.SetBodyMap("meta", func(bm bm.BodyMap) {
		bm.Set("filename", "123.jpg").
			Set("sha256", "ad4465asd4fgw5q")
	}).SetFormFile("image", &bm.File{Name: "logo.png", Content: fileContent})

	client := NewClient()

	rsp := new(HttpGet)
	_, err = client.Req(TypeMultipartFormData).
		Post("http://localhost:2233/admin/v1/oss/uploadImage").
		SendMultipartBodyMap(bmm).
		EndStruct(ctx, rsp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", rsp)
}

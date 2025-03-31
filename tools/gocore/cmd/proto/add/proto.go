package add

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Proto is a proto generator.
type Proto struct {
	Name        string
	Path        string
	Package     string
	GoPackage   string
	JavaPackage string
	Service     string
}

// Generate generate the proto files.
func (p *Proto) Generate() error {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	to := filepath.Join(wd, p.Path)
	if _, err := os.Stat(to); os.IsNotExist(err) {
		if err := os.MkdirAll(to, 0o700); err != nil {
			return err
		}
	}

	name := filepath.Join(to, p.Name)
	fmt.Println("name=" + name)
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", p.Name)
	}
	// 生成 proto 文件
	if err := p.generateProto(name); err != nil {
		return err
	}
	return nil
}

func (p *Proto) generateProto(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	t := template.Must(template.New("proto").Parse(protoTemplate))
	return t.Execute(f, p)
}

const protoTemplate = `
syntax = "proto3";

option go_package = "{{.GoPackage}}";

import "google/api/annotations.proto";

service {{.Service}} {
	rpc Create{{.Service}} (Create{{.Service}}Req) returns (Create{{.Service}}Resp) {
	  option (google.api.http) = {
	    post: "/v3/create{{.Service}}"
	    body: "*"
	  };
	}

	rpc Update{{.Service}} (Update{{.Service}}Req) returns (Update{{.Service}}Resp) {
	  option (google.api.http) = {
	    post: "/v3/update{{.Service}}"
	    body: "*"
	  };
	}

	rpc Delete{{.Service}} (Delete{{.Service}}Req) returns (Delete{{.Service}}Resp) {
	  option (google.api.http) = {
	    post: "/v3/delete{{.Service}}"
	    body: "*"
	  };
	}

	rpc Get{{.Service}} (Get{{.Service}}Req) returns (Get{{.Service}}Resp) {
	  option (google.api.http) = {
	    post: "/v3/get{{.Service}}"
	    body: "*"
	  };
	}

	rpc List{{.Service}} (List{{.Service}}Req) returns (List{{.Service}}Resp) {
	  option (google.api.http) = {
	    post: "/v3/list{{.Service}}"
	    body: "*"
	  };
	}
}

message Create{{.Service}}Req {}
message Create{{.Service}}Resp {}

message Update{{.Service}}Req {}
message Update{{.Service}}Resp {}

message Delete{{.Service}}Req {}
message Delete{{.Service}}Resp {}

message Get{{.Service}}Req {}
message Get{{.Service}}Resp {}

message List{{.Service}}Req {}
message List{{.Service}}Resp {}
`

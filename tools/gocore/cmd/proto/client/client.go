package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"crypto/md5"

	"github.com/urfave/cli/v2"
)

var protoPath string

func init() {
	if protoPath = os.Getenv("GOCORE_PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
}

// Client represents the client command.
var Client = &cli.Command{
	Name:        "client",
	Usage:       "Generate the proto client code",
	Description: `Generate the proto client code. Example: gocore proto client helloworld.proto`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "proto_path",
			Aliases:     []string{"p"},
			Value:       protoPath,
			Usage:       "proto path",
			Destination: &protoPath,
		},
	},
	Action: run,
}

// ModuleVersion returns module version.
func ModuleVersion(path string) (string, error) {
	stdout := &bytes.Buffer{}
	fd := exec.Command("go", "mod", "graph")
	fd.Stdout = stdout
	fd.Stderr = stdout
	if err := fd.Run(); err != nil {
		fmt.Println("go mod graph, err = ", err.Error())
		return "", err
	}
	rd := bufio.NewReader(stdout)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			return "", err
		}
		str := string(line)
		i := strings.Index(str, "@")
		if strings.Contains(str, path+"@") && i != -1 {
			return path + str[i:], nil
		}
	}
}

// GocoreMod returns gocore mod.
func GocoreMod() string {
	// go 1.15+ read from env GOMODCACHE
	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
	cachePath := strings.Trim(string(cacheOut), "\n")
	pathOut, _ := exec.Command("go", "env", "GOPATH").Output()
	gopath := strings.Trim(string(pathOut), "\n")
	if cachePath == "" {
		cachePath = filepath.Join(gopath, "pkg", "mod")
	}
	if path, err := ModuleVersion("github.com/sunmi-OS/gocore/v2"); err == nil {
		// $GOPATH/pkg/mod/github.com/sunmi-OS/gocore@v2
		gocorePath := filepath.Join(cachePath, path)
		gocorePath = strings.Replace(gocorePath, "sunmi-OS", "sunmi-\\!o\\!s", -1)
		return gocorePath
	}
	// $GOPATH/src/github.com/sunmi-OS/gocore
	return filepath.Join(gopath, "src", "github.com", "sunmi-OS", "gocore")
}

func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkProtocGenGoMD5() error {
	expectedMD5 := "4994a5677761d18af2bcc03f51440250"

	// 获取GOPATH
	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return fmt.Errorf("failed to get GOPATH: %v", err)
	}
	gopathStr := strings.TrimSpace(string(gopath))
	protocGenGoPath := filepath.Join(gopathStr, "bin", "protoc-gen-sm-go")

	// 检查文件是否存在
	if _, err := os.Stat(protocGenGoPath); os.IsNotExist(err) {
		return fmt.Errorf("protoc-gen-sm-go not found")
	}

	// 计算MD5
	file, err := os.Open(protocGenGoPath)
	if err != nil {
		return fmt.Errorf("failed to open protoc-gen-sm-go: %v", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate MD5: %v", err)
	}
	actualMD5 := fmt.Sprintf("%x", hash.Sum(nil))

	if actualMD5 != expectedMD5 {
		return fmt.Errorf("MD5 mismatch")
	}

	return nil
}

func downloadAndInstallProtocGenGo() error {
	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return fmt.Errorf("failed to get GOPATH: %v", err)
	}
	gopathStr := strings.TrimSpace(string(gopath))
	binPath := filepath.Join(gopathStr, "bin")

	// 创建临时文件
	tmpFile := filepath.Join(os.TempDir(), "protoc-gen-sm-go")
	defer os.Remove(tmpFile)

	// 下载文件
	resp, err := http.Get("http://qiniu.brightguo.com/sunmi/protoc-gen-go_mac13.0")
	if err != nil {
		return fmt.Errorf("failed to download protoc-gen-go: %v", err)
	}
	defer resp.Body.Close()

	// 保存到临时文件
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to save downloaded file: %v", err)
	}

	// 设置执行权限
	if err := os.Chmod(tmpFile, 0755); err != nil {
		return fmt.Errorf("failed to set executable permission: %v", err)
	}

	// 移动到GOPATH/bin
	targetPath := filepath.Join(binPath, "protoc-gen-sm-go")
	if err := os.Rename(tmpFile, targetPath); err != nil {
		return fmt.Errorf("failed to move file to GOPATH/bin: %v", err)
	}

	return nil
}

func checkAndDownloadThirdParty() error {
	thirdPartyPath := "/tmp/third_party"
	if _, err := os.Stat(thirdPartyPath); os.IsNotExist(err) {
		fmt.Println("Downloading third_party directory...")

		// 创建临时目录
		tmpDir := "/tmp/gocore_temp"
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return fmt.Errorf("failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// 克隆仓库
		cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/guoming0000/protoc-gen-go-gin.git", tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %v", err)
		}

		// 复制third_party目录
		srcPath := filepath.Join(tmpDir, "third_party")
		if err := os.MkdirAll(thirdPartyPath, 0755); err != nil {
			return fmt.Errorf("failed to create third_party directory: %v", err)
		}

		cmd = exec.Command("cp", "-r", srcPath+"/.", thirdPartyPath+"/")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy third_party directory: %v", err)
		}

		fmt.Println("Successfully downloaded third_party directory")
	}
	return nil
}

func run(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Println("Please enter the proto file or directory")
		return nil
	}
	var (
		err   error
		proto = strings.TrimSpace(c.Args().Get(0))
	)
	if err = look("yq"); err != nil {
		fmt.Println("start install yq via [ brew install yq ]")
		err0 := ExecCommand("brew", "install yq")
		if err0 != nil {
			return err0
		}
	}
	if err = look("protoc-gen-sm-go-gin", "protoc-gen-sm-openapi"); err != nil {
		// update the gocore plugins
		err0 := ExecCommand("gocore", "upgrade")
		if err0 != nil {
			return err0
		}
	}

	// 检查protoc-gen-go的MD5值
	if err = checkProtocGenGoMD5(); err != nil {
		fmt.Println("protoc-gen-sm-go MD5 check failed, downloading new version...")
		if err := downloadAndInstallProtocGenGo(); err != nil {
			return fmt.Errorf("failed to update protoc-gen-sm-go: %v", err)
		}
		fmt.Println("protoc-gen-sm-go updated successfully")
	}

	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto, c.Args().Slice())
	} else {
		err = walk(proto, c.Args().Slice())
	}
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		// 忽略docs文件夹
		if strings.HasSuffix(path, "docs/") {
			return nil
		}
		// 忽略errcode.proto文件，避免重复生成
		if strings.HasSuffix(path, "errcoe.proto") {
			return nil
		}
		fmt.Println("walk path = " + path)
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

// generate is used to execute the generate command for the specified proto file
func generate(proto string, args []string) error {
	// 检查并下载third_party目录
	if err := checkAndDownloadThirdParty(); err != nil {
		return fmt.Errorf("failed to check and download third_party: %v", err)
	}

	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	protoName := strings.TrimSuffix(path.Base(proto), ".proto")
	theDir := path.Dir(proto)

	// 开始生成golang的客户端和服务端代码
	inputExt := []string{
		//"--proto_path=" + GocoreMod(),
		"--proto_path=" + "/tmp/third_party",
		"--sm-go_out=paths=source_relative:.",
		"--sm-go-gin_out=paths=source_relative:.",
	}
	input = append(input, inputExt...)
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}

	// 生成 OpenAPI 文档
	if err := generateOpenAPI(proto, protoName, theDir); err != nil {
		return err
	}

	fmt.Printf("proto: %s\n", proto)
	return nil
}

// generateOpenAPI 生成 OpenAPI 文档
func generateOpenAPI(proto, protoName, theDir string) error {
	// 开始生成openapi
	swaggerProto := theDir + "/docs/" + protoName + ".swagger.proto"
	err := RecreateProto(proto, swaggerProto)
	if err != nil {
		return err
	}

	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}
	inputExt := []string{
		"--proto_path=" + "/tmp/third_party",
		"--sm-openapi_out=fq_schema_naming=true,default_response=false,output_mode=source_relative:.",
	}
	input = append(input, inputExt...)
	input = append(input, swaggerProto)
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	err = fd.Run()
	if err != nil {
		return err
	}

	// yq -Poj api/docs/$protoName/$protoName.swagger.yamlFile > api/docs/$protoName/$protoName.swagger.json
	yamlFile := theDir + "/docs/" + protoName + ".swagger.yaml"
	jsonFile := theDir + "/docs/" + protoName + ".swagger.json"
	file, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command("yq", "-Poj", yamlFile)
	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	cmd.Dir = "."
	if err = cmd.Run(); err != nil {
		fmt.Println("yq -Poj error:", err)
		return err
	}

	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

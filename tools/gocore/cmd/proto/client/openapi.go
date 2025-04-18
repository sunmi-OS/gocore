package client

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func RecreateProto(originProto3Path, descProto3Path string) error {
	var newFilePath string
	if descProto3Path != "" {
		newFilePath = descProto3Path
	} else {
		newFilePath = strings.ReplaceAll(originProto3Path, ".proto", ".swagger.proto")
	}

	// get fullPath path dir
	_, err := os.Stat(filepath.Dir(newFilePath))
	if err != nil {
		err = os.MkdirAll(filepath.Dir(newFilePath), 0777)
		if err != nil {
			fmt.Println("mkdir error:", err)
			return err
		}
	}

	buf := &bytes.Buffer{}

	// Read the file
	// 打开文件
	file, err := os.Open(originProto3Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个Scanner对象
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	rawMethods := []string{}
	re := regexp.MustCompile(`returns\s*\((\w+)\)`)
	structMap := map[string]bool{}
	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.Trim(rawLine, " ")
		// 忽略注释行
		if strings.HasPrefix(line, "//") {
			buf.WriteString(rawLine + "\n")
			continue
		}
		// 忽略 origin=true的行
		if strings.Contains(line, "output=origin") {
			buf.WriteString(rawLine + "\n")
			continue
		}
		match := re.FindStringSubmatch(line)
		if len(match) > 1 {
			replyStr := match[1]
			rawMethods = append(rawMethods, replyStr)
			newLine := strings.ReplaceAll(rawLine, "("+replyStr+")", "(T"+replyStr+")")
			buf.WriteString(newLine + "\n")
		} else {
			buf.WriteString(rawLine + "\n")
		}
	}

	buf.WriteString("\n")
	for _, m := range rawMethods {
		if _, ok := structMap[m]; ok {
			// 跳过一样的response
			continue
		}
		structMap[m] = true

		buf.WriteString("message T" + m + " {\n")
		buf.WriteString("   int32 code = 1; // binding:\"required\"\n")
		buf.WriteString("   string msg = 2; // binding:\"required\"\n")
		buf.WriteString("   " + m + " data = 3; // binding:\"required\"\n")
		buf.WriteString("}\n")
	}
	if err = os.WriteFile(newFilePath, buf.Bytes(), 0666); err != nil {
		return err
	}
	fmt.Println(newFilePath + " 文件生成成功")
	return nil
}

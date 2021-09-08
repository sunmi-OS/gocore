package codec

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// GzipEncode 字符串转gzip
func GzipEncode(msg string) (string, error) {
	b := new(bytes.Buffer)
	w := gzip.NewWriter(b)
	defer w.Close()

	_, err := w.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	err = w.Flush()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// GzipDecode Gzip转字符串
func GzipDecode(msg string) string {
	bf := new(bytes.Buffer)
	bf.Write([]byte(msg))

	r, err := gzip.NewReader(bf)
	if err != nil {
		return ""
	}
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	return string(undatas)
}

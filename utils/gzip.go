package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

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

func GzipEncode(msg string) string {
	b := new(bytes.Buffer)
	w := gzip.NewWriter(b)
	defer w.Close()

	w.Write([]byte(msg))
	w.Flush()
	return b.String()
}

package file

import (
	"io"
	"os"
)

// interval.间隔
const Interval = "\t"

type Writer struct {
	buff []byte
}

func NewWriter() *Writer {
	return &Writer{
		buff: make([]byte, 0),
	}
}

func (w *Writer) Add(s []byte) {
	w.buff = append(w.buff, s...)
}

func (w *Writer) Bytes() []byte {
	return w.buff
}

func (w *Writer) WriteToFile(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.WriteString(f, string(w.buff))
	w.Clear()
	if err != nil {
		panic(err)
	}
}

func (w *Writer) Clear() {
	w.buff = w.buff[:0]
}

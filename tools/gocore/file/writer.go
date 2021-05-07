package file

import (
	"io"
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

func (w *Writer) WriteToFile(file io.Writer) error {
	_, err := io.WriteString(file, string(w.buff))
	return err
}

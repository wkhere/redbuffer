package redbuffer

import (
	"bytes"
	"io"
	"strconv"
)

type Writer struct {
	llw      io.Writer
	buf      *bytes.Buffer
	colorseq []byte
}

type Color int

const (
	Red  Color = 31
	Bold       = 1
)

func New(w io.Writer) *Writer {
	return NewWithColors(w, Red, Bold)
}

func NewWithColors(w io.Writer, cc ...Color) *Writer {
	return &Writer{
		llw:      w,
		buf:      new(bytes.Buffer),
		colorseq: joinColors(cc),
	}
}

func joinColors(cc []Color) (b []byte) {
	cb := bytes.NewBuffer([]byte("\033["))
	defer func() {
		cb.WriteByte('m')
		b = cb.Bytes()
	}()

	if len(cc) == 0 {
		return
	}
	cb.WriteString(strconv.Itoa(int(cc[0])))
	for _, c := range cc[1:] {
		cb.WriteByte(';')
		cb.WriteString(strconv.Itoa(int(c)))
	}
	return
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *Writer) FlushInRed(red bool) (err error) {
	if !red {
		_, err = io.Copy(w.llw, w.buf)
		return err
	}

	_, err = w.llw.Write(w.colorseq)
	if err != nil {
		return err
	}
	_, err = io.Copy(w.llw, w.buf)
	if err != nil {
		return err
	}
	_, err = w.llw.Write([]byte("\033[0m"))
	return err
}

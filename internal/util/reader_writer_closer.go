package util

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
)

var _ io.ReadWriteCloser = &readerWriterCloser{}

type readerWriterCloser struct {
	r io.ReadCloser
	w io.WriteCloser
}

func NewReaderWriterCloser(r io.ReadCloser, w io.WriteCloser) io.ReadWriteCloser {
	return &readerWriterCloser{
		r: r,
		w: w,
	}
}

func (s *readerWriterCloser) Read(p []byte) (n int, err error) {
	fmt.Printf("-> Read(%d)\n", len(p))
	n, err = s.r.Read(p)

	fmt.Printf("<- Read %d: %v\n", n, err)
	return
}

func (s *readerWriterCloser) Write(p []byte) (n int, err error) {
	fmt.Printf("-> Write(%d)\n", len(p))
	n, err = s.w.Write(p)
	fmt.Printf("<- Write %d, %v\n", n, err)
	return
}

func (s *readerWriterCloser) Close() error {
	var result error

	if err := s.r.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	if err := s.w.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

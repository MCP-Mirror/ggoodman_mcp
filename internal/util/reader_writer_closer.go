package util

import (
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

func (s *readerWriterCloser) Read(p []byte) (int, error) {
	return s.r.Read(p)
}

func (s *readerWriterCloser) Write(p []byte) (int, error) {
	return s.w.Write(p)
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

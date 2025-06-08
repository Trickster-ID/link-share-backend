package helper

import (
	"bytes"
	"compress/gzip"
	"github.com/sirupsen/logrus"
	"io"
)

func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer func(gz *gzip.Reader) {
		err := gz.Close()
		if err != nil {
			logrus.Errorf("error while closing gzip reader: %v", err)
		}
	}(gz)
	return io.ReadAll(gz)
}

package diskcache

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type compressionType byte

const (
	CompressionNone compressionType = iota
	CompressionGzip
)

func getCompression(input compressionType) compression {
	switch input {
	case CompressionNone:
		return &compressionNone{}
	case CompressionGzip:
		return &compressionGzip{}
	default:
		panic("compression type not implemented")
	}
}

type compression interface {
	compress(input []byte) ([]byte, error)
	decompress(input []byte) ([]byte, error)
}

type compressionNone struct {
}

func (c *compressionNone) compress(input []byte) ([]byte, error) {
	return input, nil
}

func (c *compressionNone) decompress(input []byte) ([]byte, error) {
	return input, nil
}

type compressionGzip struct {
}

func (c *compressionGzip) compress(input []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write(input); err != nil {
		return []byte{}, err
	}
	if err := w.Close(); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

func (c *compressionGzip) decompress(input []byte) ([]byte, error) {
	fi := bytes.NewReader(input)
	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	return ioutil.ReadAll(fz)
}

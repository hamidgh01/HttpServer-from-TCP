package server

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
)

type CompressionAlg string

const (
	Gzip    CompressionAlg = "gzip"
	Deflate CompressionAlg = "deflate"
	// Zstd    CompressionAlg = "zstd"
	// Br      CompressionAlg = "br"
)

var ErrCompressionAlgorithmNotSupported = errors.New("compression algorithm is not supported")

func Compress(bodyBytes []byte, compressionAlg CompressionAlg) ([]byte, error) {
	switch compressionAlg {
	case Gzip:
		return compressGzip(bodyBytes)
	case Deflate:
		return compressDeflate(bodyBytes)
	}
	// case Br: (maybe add later)
	// to implement `br` compression, github.com/klauspost/compress/zstd is recommended
	// case Zstd: (maybe add later)
	// to implement `zstd` compression, github.com/andybalholm/brotli and
	// github.com/google/brotli are recommended

	return nil, ErrCompressionAlgorithmNotSupported
}

func Decompress() {
	// implement later
}

func compressGzip(bodyBytes []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	gzWriter, _ := gzip.NewWriterLevel(buf, gzip.BestCompression)

	if _, err := gzWriter.Write(bodyBytes); err != nil {
		return nil, fmt.Errorf("error while gzip compression process: %w", err)
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("error while closing (flushing) gzipWriter: %w", err)
	}

	return buf.Bytes(), nil
}

func decompressGzip() {
	// implement later
}

func compressDeflate(bodyBytes []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	flateWriter, _ := flate.NewWriter(buf, flate.BestCompression)

	if _, err := flateWriter.Write(bodyBytes); err != nil {
		return nil, fmt.Errorf("error while deflate compression process: %w", err)
	}

	if err := flateWriter.Close(); err != nil {
		return nil, fmt.Errorf("error while closing (flushing) flateWriter: %w", err)
	}

	return buf.Bytes(), nil
}

func decompressDeflate() {
	// implement later
}

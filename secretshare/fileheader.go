// Package secretshare implements Shamir secret sharing for files and directories.
// Shamir's scheme allows to split a secret value into multiple parts, so called shares, by evaluating a secret sharing polynomial at certain indices.
// The shared secret can only be reconstructed (via Lagrange interpolation) if a threshold of the participants provide their shares.
// This package provides functions to easily create secret shares for a file/directory, recover and validate secret file/directory from shares.
package secretshare

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
)

const (
	// CurVersion stores in sharing files' head, used for checking format
	CurVersion       = 0
	defaultBlockSize = 16 * 1024
	versionLen       = 4
	blockSizeLen     = 4
	knumLen          = 1
	fileIDLen        = 11
	headerLen        = versionLen + blockSizeLen + fileIDLen
)

type shareFileHeader struct {
	version   uint32
	blockSize uint32
	knum      uint8
	fileID    [fileIDLen]byte
}

func newShareFileHeader() *shareFileHeader {
	var hd shareFileHeader
	hd.version = CurVersion
	hd.blockSize = defaultBlockSize
	io.ReadFull(rand.Reader, hd.fileID[:])
	return &hd
}

func (hd *shareFileHeader) searialize() []byte {
	data := make([]byte, headerLen)
	binary.BigEndian.PutUint32(data, hd.version)
	binary.BigEndian.PutUint32(data[versionLen:], hd.blockSize)
	data[versionLen+blockSizeLen] = hd.knum
	copy(data[versionLen+blockSizeLen+knumLen:], hd.fileID[:])
	return data
}

func (hd *shareFileHeader) unsearialize(data []byte) {
	hd.version = binary.BigEndian.Uint32(data[:versionLen])
	hd.blockSize = binary.BigEndian.Uint32(data[versionLen : versionLen+blockSizeLen])
	hd.knum = data[versionLen+blockSizeLen]
	copy(hd.fileID[:], data[versionLen+blockSizeLen+knumLen:])
}

func (hd *shareFileHeader) equal(hd2 *shareFileHeader) bool {
	return hd.version == hd2.version && hd.blockSize == hd2.blockSize && bytes.Equal(hd.fileID[:], hd2.fileID[:])
}

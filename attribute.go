package jvmgo

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	AttributeInfo struct {
		AttributeNameIndex uint16
		AttributeLength    uint32
		Info               []byte
	}
)

func readAttribute(r io.Reader) (*AttributeInfo, error) {
	nBuf := make([]byte, 2)
	if _, err := io.ReadFull(r, nBuf); err != nil {
		return nil, fmt.Errorf("read attribute: %w", err)
	}

	lBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lBuf); err != nil {
		return nil, fmt.Errorf("read attribute length: %w", err)
	}
	length := binary.BigEndian.Uint32(lBuf)

	info := make([]byte, length)
	buf := make([]byte, 1)
	for i := range info {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, fmt.Errorf("read attribute info idx=%d: %w", i, err)
		}
		info[i] = buf[0]
	}

	return &AttributeInfo{
		AttributeNameIndex: binary.BigEndian.Uint16(nBuf),
		AttributeLength:    length,
		Info:               info,
	}, nil
}

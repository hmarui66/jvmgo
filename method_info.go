package jvmgo

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	MethodInfo struct {
		AccessFlags     uint16
		NameIndex       uint16
		DescriptorIndex uint16
		AttributesCount uint16
		Attributes      []*AttributeInfo
	}
)

func readMethod(r io.Reader) (*MethodInfo, error) {
	ret := &MethodInfo{}
	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read method access flags: %w", err)
	}
	ret.AccessFlags = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read method name index: %w", err)
	}
	ret.NameIndex = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read method descriptor index: %w", err)
	}
	ret.DescriptorIndex = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read method attributes count: %w", err)
	}
	ret.AttributesCount = binary.BigEndian.Uint16(buf)
	if err := ret.readAttributes(r); err != nil {
		return nil, fmt.Errorf("read method attributes: %w", err)
	}

	return ret, nil
}

func (m *MethodInfo) readAttributes(r io.Reader) error {
	var err error
	m.Attributes = make([]*AttributeInfo, m.AttributesCount)
	for i := range m.Attributes {
		if m.Attributes[i], err = readAttribute(r); err != nil {
			return fmt.Errorf("read method attribute idx=%d: %w", i, err)
		}
	}

	return nil
}

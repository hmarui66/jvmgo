package jvmgo

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	FieldInfo struct {
		AccessFlags     uint16
		NameIndex       uint16
		DescriptorIndex uint16
		AttributesCount uint16
		Attributes      []*AttributeInfo
	}
)

func readField(r io.Reader) (*FieldInfo, error) {
	ret := &FieldInfo{}
	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read field access flags: %w", err)
	}
	ret.AccessFlags = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read field name index: %w", err)
	}
	ret.NameIndex = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read field descriptor index: %w", err)
	}
	ret.DescriptorIndex = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read field attributes count: %w", err)
	}
	ret.AttributesCount = binary.BigEndian.Uint16(buf)
	if err := ret.readAttributes(r); err != nil {
		return nil, fmt.Errorf("read field attributes: %w", err)
	}

	return ret, nil
}

func (f *FieldInfo) readAttributes(r io.Reader) error {
	var err error
	f.Attributes = make([]*AttributeInfo, f.AttributesCount)
	for i := range f.Attributes {
		if f.Attributes[i], err = readAttribute(r); err != nil {
			return fmt.Errorf("read field attribute idx=%d: %w", i, err)
		}
	}

	return nil
}

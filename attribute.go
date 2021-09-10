package jvmgo

import (
	"bytes"
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
	if _, err := io.ReadFull(r, info); err != nil {
		return nil, fmt.Errorf("read attribute info: %w", err)
	}

	return &AttributeInfo{
		AttributeNameIndex: binary.BigEndian.Uint16(nBuf),
		AttributeLength:    length,
		Info:               info,
	}, nil
}

func (a *AttributeInfo) toCodeAttribute() (*CodeAttribute, error) {
	r := bytes.NewBuffer(a.Info)
	ret := &CodeAttribute{
		AttributeNameIndex: a.AttributeNameIndex,
		AttributeLength:    a.AttributeLength,
	}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read max stack: %w", err)
	}
	ret.MaxStack = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read max locals: %w", err)
	}
	ret.MaxLocals = binary.BigEndian.Uint16(buf)

	lBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lBuf); err != nil {
		return nil, fmt.Errorf("read code length: %w", err)
	}
	ret.CodeLength = binary.BigEndian.Uint32(lBuf)

	code := make([]byte, ret.CodeLength)
	if _, err := io.ReadFull(r, code); err != nil {
		return nil, fmt.Errorf("read code: %w", err)
	}
	ret.Code = code

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read exception table length: %w", err)
	}
	ret.ExceptionTableLength = binary.BigEndian.Uint16(buf)
	if err := ret.readExceptions(r); err != nil {
		return nil, fmt.Errorf("read exceptions: %w", err)
	}

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read attributes count: %w", err)
	}
	ret.AttributesCount = binary.BigEndian.Uint16(buf)
	if err := ret.readAttributes(r); err != nil {
		return nil, fmt.Errorf("read attributes: %w", err)
	}

	return ret, nil
}
func (c *CodeAttribute) readExceptions(r io.Reader) error {
	var err error
	c.ExceptionTable = make([]*Exception, c.ExceptionTableLength)
	for i := range c.ExceptionTable {
		if c.ExceptionTable[i], err = readException(r); err != nil {
			return fmt.Errorf("read exception idx=%d: %w", i, err)
		}
	}

	return nil
}

func (c *CodeAttribute) readAttributes(r io.Reader) error {
	var err error
	c.Attributes = make([]*AttributeInfo, c.AttributesCount)
	for i := range c.Attributes {
		if c.Attributes[i], err = readAttribute(r); err != nil {
			return fmt.Errorf("read attribute idx=%d: %w", i, err)
		}
	}

	return nil
}

func readException(r io.Reader) (*Exception, error) {
	ret := &Exception{}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read start pc: %w", err)
	}
	ret.StartPC = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read end pc: %w", err)
	}
	ret.EndPC = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read handler pc: %w", err)
	}
	ret.HandlerPC = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read catch type: %w", err)
	}
	ret.CatchType = binary.BigEndian.Uint16(buf)

	return ret, nil
}

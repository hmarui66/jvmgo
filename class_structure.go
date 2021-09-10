package jvmgo

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	magic                 = []byte{0xCA, 0xFE, 0xBA, 0xBE}
	minorVersion          = []byte{0x00, 0x00}
	majorVersion          = []byte{0x00, 0x37}
	ErrInvalidMagicNumber = errors.New("invalid magic number")
	ErrInvalidVersion     = errors.New("invalid version")
)

type (
	ClassStructure struct {
		Magic             []byte
		MinorVersion      []byte
		MajorVersion      []byte
		ConstantPoolCount uint16
		ConstantPool      []*CpInfo
		AccessFlags       uint16
		ThisClass         uint16
		SuperClass        uint16
		InterfacesCount   uint16
		Interfaces        []uint16
		FieldsCount       uint16
		Fields            []*FieldInfo
		MethodsCount      uint16
		Methods           []*MethodInfo
		AttributesCount   uint16
		Attributes        []*AttributeInfo
	}
)

func DecodeClassStructure(r io.Reader) (*ClassStructure, error) {
	ret := &ClassStructure{}
	mbuf := make([]byte, 4)
	if _, err := io.ReadFull(r, mbuf); err != nil {
		return nil, ErrInvalidMagicNumber
	}

	for i := range mbuf {
		if mbuf[i] != magic[i] {
			return nil, ErrInvalidMagicNumber
		}
	}
	ret.Magic = mbuf

	ret.MinorVersion = make([]byte, 2)
	if _, err := io.ReadFull(r, ret.MinorVersion); err != nil {
		return nil, ErrInvalidVersion
	}
	for i := range ret.MinorVersion {
		if ret.MinorVersion[i] != minorVersion[i] {
			return nil, ErrInvalidVersion
		}
	}

	ret.MajorVersion = make([]byte, 2)
	if _, err := io.ReadFull(r, ret.MajorVersion); err != nil {
		return nil, ErrInvalidVersion
	}
	for i := range ret.MajorVersion {
		if ret.MajorVersion[i] != majorVersion[i] {
			return nil, ErrInvalidVersion
		}
	}

	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read constant pool count: %w", err)
	}
	ret.ConstantPoolCount = binary.BigEndian.Uint16(buf)
	if err := ret.readConstantPool(r); err != nil {
		return nil, fmt.Errorf("read constant pool: %w", err)
	}

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read access flags: %w", err)
	}
	ret.AccessFlags = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read this class: %w", err)
	}
	ret.ThisClass = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read super class: %w", err)
	}
	ret.SuperClass = binary.BigEndian.Uint16(buf)

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read interfaces count: %w", err)
	}
	ret.InterfacesCount = binary.BigEndian.Uint16(buf)
	if err := ret.readInterfaces(r); err != nil {
		return nil, fmt.Errorf("read interfaces: %w", err)
	}

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read field count: %w", err)
	}
	ret.FieldsCount = binary.BigEndian.Uint16(buf)
	if err := ret.readFields(r); err != nil {
		return nil, fmt.Errorf("read fields: %w", err)
	}

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("read methods count: %w", err)
	}
	ret.MethodsCount = binary.BigEndian.Uint16(buf)
	if err := ret.readMethods(r); err != nil {
		return nil, fmt.Errorf("read methods: %w", err)
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

func (c *ClassStructure) readConstantPool(r io.Reader) error {
	var err error
	c.ConstantPool = make([]*CpInfo, c.ConstantPoolCount-1)
	for i := range c.ConstantPool {
		if c.ConstantPool[i], err = readCpInfo(r); err != nil {
			return fmt.Errorf("read constant pool idx=%d: %w", i, err)
		}
	}

	return nil
}

func (c *ClassStructure) readInterfaces(r io.Reader) error {
	buf := make([]byte, 2)
	c.Interfaces = make([]uint16, c.InterfacesCount)
	for i := range c.Interfaces {
		if _, err := io.ReadFull(r, buf); err != nil {
			return fmt.Errorf("read interface idx=%d: %w", i, err)
		}
		c.Interfaces[i] = binary.BigEndian.Uint16(buf)
	}

	return nil
}

func (c *ClassStructure) readFields(r io.Reader) error {
	var err error
	c.Fields = make([]*FieldInfo, c.FieldsCount)
	for i := range c.Fields {
		if c.Fields[i], err = readField(r); err != nil {
			return fmt.Errorf("read field idx=%d: %w", i, err)
		}
	}

	return nil
}

func (c *ClassStructure) readMethods(r io.Reader) error {
	var err error
	c.Methods = make([]*MethodInfo, c.MethodsCount)
	for i := range c.Methods {
		if c.Methods[i], err = readMethod(r); err != nil {
			return fmt.Errorf("read method idx=%d: %w", i, err)
		}
	}

	return nil
}

func (c *ClassStructure) readAttributes(r io.Reader) error {
	var err error
	c.Attributes = make([]*AttributeInfo, c.AttributesCount)
	for i := range c.Attributes {
		if c.Attributes[i], err = readAttribute(r); err != nil {
			return fmt.Errorf("read attribute idx=%d: %w", i, err)
		}
	}

	return nil
}

func (c *ClassStructure) GetCpInfo(idx uint16) *CpInfo {
	return c.ConstantPool[idx-1]
}

package jvmgo

import (
	"encoding/binary"
	"fmt"
	"io"
)

type (
	ConstantKind byte
	CpInfo       struct {
		Tag  ConstantKind
		Info []byte
	}
	Fieldref struct {
		Tag              ConstantKind
		ClassIndex       uint16
		NameAndTypeIndex uint16
	}
	Methodref struct {
		Tag              ConstantKind
		ClassIndex       uint16
		NameAndTypeIndex uint16
	}
	Class struct {
		Tag       ConstantKind
		NameIndex uint16
	}
	NameAndType struct {
		Tag             ConstantKind
		NameIndex       uint16
		DescriptorIndex uint16
	}
)

const (
	ConstantKindClass              ConstantKind = 7
	ConstantKindFieldref           ConstantKind = 9
	ConstantKindMethodref          ConstantKind = 10
	ConstantKindInterfaceMethodref ConstantKind = 11
	ConstantKindString             ConstantKind = 8
	ConstantKindInteger            ConstantKind = 3
	ConstantKindFloat              ConstantKind = 4
	ConstantKindLong               ConstantKind = 5
	ConstantKindDouble             ConstantKind = 6
	ConstantKindNameAndType        ConstantKind = 12
	ConstantKindUTF8               ConstantKind = 1
	ConstantKindMethodHandle       ConstantKind = 15
	ConstantKindMethodType         ConstantKind = 16
	ConstantKindDynamic            ConstantKind = 17
	ConstantKindInvokeDynamic      ConstantKind = 18
	ConstantKindModule             ConstantKind = 19
	ConstantKindPackage            ConstantKind = 20
)

func readCpInfo(r io.Reader) (*CpInfo, error) {
	tBuf := make([]byte, 1)
	if _, err := io.ReadFull(r, tBuf); err != nil {
		return nil, fmt.Errorf("read cp info: %w", err)
	}

	tag := ConstantKind(tBuf[0])

	var info []byte
	switch tag {
	case ConstantKindClass:
		info = make([]byte, 2)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind class: %w", err)
		}
	case ConstantKindFieldref:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind field ref: %w", err)
		}
	case ConstantKindMethodref:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind method ref: %w", err)
		}
	case ConstantKindInterfaceMethodref:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind interface method ref: %w", err)
		}
	case ConstantKindString:
		info = make([]byte, 2)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind string: %w", err)
		}
	case ConstantKindInteger, ConstantKindFloat:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind integer or float: %w", err)
		}
	case ConstantKindLong, ConstantKindDouble:
		info = make([]byte, 8)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind long or double: %w", err)
		}
	case ConstantKindNameAndType:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind name and type: %w", err)
		}
	case ConstantKindUTF8:
		length := make([]byte, 2)
		if _, err := io.ReadFull(r, length); err != nil {
			return nil, fmt.Errorf("read constant kind utf8 length: %w", err)
		}
		buf := make([]byte, binary.BigEndian.Uint16(length))
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, fmt.Errorf("read constant kind utf8 bytes: %w", err)
		}
		info = append(info, length...)
		info = append(info, buf...)
	case ConstantKindMethodHandle:
		info = make([]byte, 3)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind method handle: %w", err)
		}
	case ConstantKindMethodType:
		info = make([]byte, 2)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind method type: %w", err)
		}
	case ConstantKindDynamic, ConstantKindInvokeDynamic:
		info = make([]byte, 4)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind dynamic or invoke dynamic: %w", err)
		}
	case ConstantKindModule:
		info = make([]byte, 2)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind module: %w", err)
		}
	case ConstantKindPackage:
		info = make([]byte, 2)
		if _, err := io.ReadFull(r, info); err != nil {
			return nil, fmt.Errorf("read constant kind package: %w", err)
		}
	}

	return &CpInfo{
		Tag:  tag,
		Info: info,
	}, nil
}

func (c *CpInfo) GetAsUTF8String() (string, error) {
	if c.Tag != ConstantKindUTF8 {
		return "", fmt.Errorf("constant kind mismatch. kind should be UTF8")
	}
	if len(c.Info) < 3 {
		return "", fmt.Errorf("cp info is invalid as kind UTF8")
	}
	return string(c.Info[2:]), nil
}

func (c *CpInfo) ToFieldRef() (*Fieldref, error) {
	if c.Tag != ConstantKindFieldref {
		return nil, fmt.Errorf("constant kind mismatch. kind should be field ref")
	}
	if len(c.Info) < 2 {
		return nil, fmt.Errorf("cp info is invalid as kind field ref")
	}

	return &Fieldref{
		Tag:              c.Tag,
		ClassIndex:       binary.BigEndian.Uint16(c.Info[:2]),
		NameAndTypeIndex: binary.BigEndian.Uint16(c.Info[2:]),
	}, nil
}

func (c *CpInfo) ToClass() (*Class, error) {
	if c.Tag != ConstantKindClass {
		return nil, fmt.Errorf("constant kind mismatch. kind should be class")
	}
	if len(c.Info) < 1 {
		return nil, fmt.Errorf("cp info is invalid as kind class")
	}

	return &Class{
		Tag:       c.Tag,
		NameIndex: binary.BigEndian.Uint16(c.Info),
	}, nil
}

func (c *CpInfo) ToNameAndType() (*NameAndType, error) {
	if c.Tag != ConstantKindNameAndType {
		return nil, fmt.Errorf("constant kind mismatch. kind should be name and type")
	}
	if len(c.Info) < 2 {
		return nil, fmt.Errorf("cp info is invalid as kind name and type")
	}
	return &NameAndType{
		Tag:             c.Tag,
		NameIndex:       binary.BigEndian.Uint16(c.Info[:2]),
		DescriptorIndex: binary.BigEndian.Uint16(c.Info[2:]),
	}, nil
}

func (c *CpInfo) ToMethodRef() (*Methodref, error) {
	if c.Tag != ConstantKindMethodref {
		return nil, fmt.Errorf("constant kind mismatch. kind should be method ref")
	}
	if len(c.Info) < 2 {
		return nil, fmt.Errorf("cp info is invalid as kind method ref")
	}

	return &Methodref{
		Tag:              c.Tag,
		ClassIndex:       binary.BigEndian.Uint16(c.Info[:2]),
		NameAndTypeIndex: binary.BigEndian.Uint16(c.Info[2:]),
	}, nil
}

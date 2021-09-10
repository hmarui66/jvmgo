package jvmgo

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestVirtualMachine_ExecMain(t *testing.T) {
	buf, err := ioutil.ReadFile("HelloWorld.class")
	require.NoError(t, err)

	class, err := DecodeClassStructure(bytes.NewBuffer(buf))
	require.NoError(t, err)
	vm := NewVM(class)

	err = vm.ExecMain()
	require.NoError(t, err)

	//	c := CodeAttribute{
	//		AttributeNameIndex:   9,
	//		AttributeLength:      55,
	//		MaxStack:             2,
	//		MaxLocals:            1,
	//		CodeLength:           9,
	//		Code:                 []byte{178, 0, 2, 18, 3, 182, 0, 4, 177},
	//		ExceptionTableLength: 0,
	//		ExceptionTable:       nil,
	//		AttributesCount:      2,
	//		Attributes: []*AttributeInfo{{
	//			AttributeNameIndex: 0xc0,
	//			AttributeLength:    0x00,
	//			Info:               []byte{0x15, 0x87, 0x80},
	//		}, {
	//			AttributeNameIndex: 0xc0,
	//			AttributeLength:    0x00,
	//			Info:               []byte{0x15, 0x87, 0xa0},
	//		}},
	//	}
}

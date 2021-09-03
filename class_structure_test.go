package jvmgo

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestDecodeClassStructure(t *testing.T) {
	buf, err := ioutil.ReadFile("HelloWorld.class")
	require.NoError(t, err)

	class, err := DecodeClassStructure(bytes.NewBuffer(buf))
	require.NoError(t, err)
	require.Equal(t, magic, class.Magic)
	require.Equal(t, minorVersion, class.MinorVersion)
	require.Equal(t, majorVersion, class.MajorVersion)

	require.Equal(t, uint16(29), class.ConstantPoolCount)
	require.Len(t, class.ConstantPool, 28) // 仕様によると constant_pool のエントリ数は constant_pool_count - 1 とのこと

	require.Equal(t, uint16(33), class.AccessFlags)
	require.Equal(t, uint16(5), class.ThisClass)
	require.Equal(t, uint16(6), class.SuperClass)

	require.Equal(t, uint16(0), class.FieldsCount)
	require.Len(t, class.Fields, 0)

	require.Equal(t, uint16(2), class.MethodsCount)
	require.Len(t, class.Methods, 2)

	require.Equal(t, uint16(1), class.AttributesCount)
	require.Len(t, class.Attributes, 1)
}

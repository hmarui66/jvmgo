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
}

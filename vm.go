package jvmgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	OpCodeGetStatic     OpCode = 0xb2
	OpCodeLdc           OpCode = 0x12
	OpCodeInvokeVirtual OpCode = 0xb6
	OpCodeReturn        OpCode = 0xB1
)

type (
	VirtualMachine struct {
		Class        *ClassStructure
		OperandStack []uint64
	}
	OpCode uint8
)

func NewVM(class *ClassStructure) *VirtualMachine {
	return &VirtualMachine{
		Class: class,
	}
}

func (vm *VirtualMachine) ExecMain() error {
	for _, methodInfo := range vm.Class.Methods {
		methodName, err := vm.Class.GetCpInfo(methodInfo.NameIndex).GetAsUTF8String()
		if err != nil {
			return fmt.Errorf("get method name: %w", err)
		}
		if methodName == "main" {
			if err := vm.execute(methodInfo); err != nil {
				return fmt.Errorf("execute main. %v: %w", methodInfo, err)
			}
			fmt.Printf("finished!: %v\n", methodInfo)
			return nil
		}
	}

	return fmt.Errorf("main method does not exist")
}

func (vm *VirtualMachine) execute(m *MethodInfo) error {
	for i := range m.Attributes {
		code, err := m.Attributes[i].toCodeAttribute()
		if err != nil {
			return fmt.Errorf("convert to code attribute :%w", err)
		}
		if err := vm.executeCode(code); err != nil {
			return fmt.Errorf("execute code :%w", err)
		}
		fmt.Printf("executed %v\n", code)
	}

	return nil
}

func (vm *VirtualMachine) executeCode(c *CodeAttribute) error {
	r := bytes.NewBuffer(c.Code)
	for {
		op := make([]byte, 1)
		if _, err := io.ReadFull(r, op); err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("execute remaining buf len=%d :%w", r.Len(), err)
		}

		switch OpCode(op[0]) {
		case OpCodeGetStatic:
			buf := make([]byte, 2)
			if _, err := io.ReadFull(r, buf); err != nil {
				return fmt.Errorf("execute get static :%w", err)
			}
			symbol, err := vm.Class.GetCpInfo(binary.BigEndian.Uint16(buf)).ToFieldRef()
			if err != nil {
				return fmt.Errorf("execute get static :%w", err)
			}

			fmt.Printf("executing getstatic. class: %v\n", vm.Class.GetCpInfo(symbol.ClassIndex))
			fmt.Printf("executing getstatic. name and type: %v\n", vm.Class.GetCpInfo(symbol.NameAndTypeIndex))

			// TODO: continue execution
			return nil
		case OpCodeLdc:
		case OpCodeInvokeVirtual:
		case OpCodeReturn:
		}
	}
}

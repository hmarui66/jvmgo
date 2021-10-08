package jvmgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
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
		OperandStack *OperandStack
	}
	OperandStack []*CpInfo
	OpCode       uint8
)

func NewVM(class *ClassStructure) *VirtualMachine {
	return &VirtualMachine{
		Class:        class,
		OperandStack: &OperandStack{},
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
				return fmt.Errorf("execute get static symbol: %w", err)
			}
			vm.OperandStack.push(vm.Class.GetCpInfo(binary.BigEndian.Uint16(buf)))
		case OpCodeLdc:
			buf := make([]byte, 1)
			if _, err := io.ReadFull(r, buf); err != nil {
				return fmt.Errorf("execute ldc: %w", err)
			}
			vm.OperandStack.push(vm.Class.GetCpInfo(binary.BigEndian.Uint16([]byte{0, buf[0]})))
		case OpCodeInvokeVirtual:
			buf := make([]byte, 2)
			if _, err := io.ReadFull(r, buf); err != nil {
				return fmt.Errorf("execute invoke virtual: %w", err)
			}
			symbol, err := vm.Class.GetCpInfo(binary.BigEndian.Uint16(buf)).ToMethodRef()
			if err != nil {
				return fmt.Errorf("execute invoke virtual parse symbol: %w", err)
			}
			calleeInfo, err := vm.Class.GetCpInfo(symbol.NameAndTypeIndex).ToNameAndType()
			if err != nil {
				return fmt.Errorf("execute invoke virtual name and type: %w", err)
			}
			methodName, err := vm.Class.GetCpInfo(calleeInfo.NameIndex).GetAsUTF8String()
			if err != nil {
				return fmt.Errorf("execute invoke virtual method name: %w", err)
			}
			argumentInfo, err := vm.Class.GetCpInfo(calleeInfo.DescriptorIndex).GetAsUTF8String()
			if err != nil {
				return fmt.Errorf("execute invoke virtual arguments: %w", err)
			}

			arguments := make([]*CpInfo, len(strings.Split(argumentInfo, ";"))-1)
			var exists bool
			for i := range arguments {
				arguments[i], exists = vm.OperandStack.pop()
				if !exists {
					return fmt.Errorf("execute invoke virtual pop arguments: arguments not exist. idx: %d", i)
				}
			}
			callableInfo, exists := vm.OperandStack.pop()
			if !exists {
				return fmt.Errorf("execute invoke virtual pop callable info: callable not exist")
			}
			callable, err := callableInfo.ToFieldRef()
			if err != nil {
				return fmt.Errorf("execute invoke virtual parse callable: %w", err)
			}
			callableClass, err := vm.Class.GetCpInfo(callable.ClassIndex).ToClass()
			if err != nil {
				return fmt.Errorf("execute invoke virtual callable class: %w", err)
			}
			callableClassName, err := vm.Class.GetCpInfo(callableClass.NameIndex).GetAsUTF8String()
			if err != nil {
				return fmt.Errorf("execute invoke virtual callable class name: %w", err)
			}

			nameAndType, err := vm.Class.GetCpInfo(callable.NameAndTypeIndex).ToNameAndType()
			if err != nil {
				return fmt.Errorf("execute invoke virtual name and type: %w", err)
			}
			fieldType, err := vm.Class.GetCpInfo(nameAndType.DescriptorIndex).GetAsUTF8String()
			if err != nil {
				return fmt.Errorf("execute invoke virtual field type: %w", err)
			}

			switch callableClassName {
			case "java/lang/System":
				switch fieldType {
				case "Ljava/io/PrintStream;":
					f := PrintStream{}
					switch methodName {
					case "println":
						args := make([]interface{}, len(arguments))
						for i := range arguments {
							args[i], _ = vm.Class.GetCpInfo(binary.BigEndian.Uint16(arguments[i].Info)).GetAsUTF8String()
						}
						f.println(args...)
					}
				}
			}
		case OpCodeReturn:
			break
		}
	}
}

func (s *OperandStack) push(ope *CpInfo) {
	*s = append(*s, ope)
}

func (s *OperandStack) pop() (*CpInfo, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	idx := len(*s) - 1
	res := (*s)[idx]
	*s = (*s)[:idx]
	return res, true
}

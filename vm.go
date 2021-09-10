package jvmgo

import "fmt"

type (
	VirtualMachine struct {
		Class *ClassStructure
	}
)

func NewVM(class *ClassStructure) *VirtualMachine {
	return &VirtualMachine{
		Class: class,
	}
}

func (vm *VirtualMachine) ExecMain() error {
	for _, methodInfo := range vm.Class.Methods {
		methodName, err := vm.Class.ConstantPool[methodInfo.NameIndex-1].GetAsUTF8String()
		if err != nil {
			return fmt.Errorf("get method name: %w", err)
		}
		if methodName == "main" {
			fmt.Printf("execute main!: %v\n", methodInfo)
			return nil
		}
	}

	return fmt.Errorf("main method does not exist")
}

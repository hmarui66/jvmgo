package jvmgo

type (
	CodeAttribute struct {
		AttributeNameIndex   uint16
		AttributeLength      uint32
		MaxStack             uint16
		MaxLocals            uint16
		CodeLength           uint32
		Code                 []byte
		ExceptionTableLength uint16
		ExceptionTable       []*Exception
		AttributesCount      uint16
		Attributes           []*AttributeInfo
	}

	Exception struct {
		StartPC   uint16
		EndPC     uint16
		HandlerPC uint16
		CatchType uint16
	}
)

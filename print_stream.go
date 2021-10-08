package jvmgo

import "fmt"

type PrintStream struct{}

func (PrintStream) println(args ...interface{}) {
	fmt.Println(args...)
}

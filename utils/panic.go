package utils

import (
	"fmt"
	"runtime/debug"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RecoverWithStack() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		fmt.Printf("panic: %v\n%s", r, string(stack))
	}
}

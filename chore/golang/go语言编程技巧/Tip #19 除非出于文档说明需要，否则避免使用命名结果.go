package main

import "fmt"

// !长函数中避免使用裸露(naked)返回，简短的函数中可以使用。
func calculateStatus(a, b int) (sum, product int) {
	sum = a + b
	product = a * b
	return
}

// !若需在延迟函数调用中修改返回值，为结果参数命名至关重要。
func operation() (res int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("caught panic: %v", r)
		}
	}()

	res = 1 // some operation that may panic
	return
}

// 当函数返回相同类型的对象时，特别是在某一类型的成员方法中，为每个返回的对象命名可能会造成冗余，并使我们的文档显得杂乱。
// !或者，该类型本身可能已经具有自解释性。

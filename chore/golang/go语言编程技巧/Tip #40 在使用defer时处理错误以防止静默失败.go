package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	err1 := badDoSomething()
	fmt.Println(err1)

	err2 := goodDoSomething()
	fmt.Println(err2)
}

func badDoSomething() error {
	file, err := os.Open("filename")
	if err != nil {
		return err
	}
	defer file.Close() // !忘记检查延迟调用中的错误

	// do something with file

	return nil
}

func goodDoSomething() error {
	file, err := os.Open("filename")
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close()) // !检查延迟调用中的错误
	}()

	// do something with file

	return nil
}

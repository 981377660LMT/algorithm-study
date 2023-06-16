package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	inStream := new(bytes.Buffer)
	w := bufio.NewWriter(inStream)
	_, err := w.Write([]byte("foo"))
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte("bar"))
	if err != nil {
		panic(err)
	}
	err = w.Flush()
	if err != nil {
		panic(err)
	}

	inStream.WriteTo(os.Stdout)
	fmt.Print("\n")

	/////////////////////////////////////////////////////
	/////////////////////////////////////////////////////

	outStream := new(bytes.Buffer)
	outStream.Write([]byte("abc\n"))
	outStream.Write([]byte("xyc\n"))
	pipeIn, pipeOut := io.Pipe()

	go func() {
		sc := bufio.NewScanner(pipeIn)
		for sc.Scan() { // hasNext
			fmt.Println("received: " + sc.Text()) // next
		}
		if err := sc.Err(); err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	go func() {
		defer pipeOut.Close()
		io.Copy(pipeOut, outStream)
	}()

	defer runtime.Goexit()
}

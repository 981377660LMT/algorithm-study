// Tip #64: Make main() clean and testable
// 让main()函数更清晰并且易于测试
//
// !main()函数现在只是替run函数把东西准备好，用不同的参数集来测试“run函数”部分是如何独立工作的。

package main

import (
	"log"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	// loadConfig()
	// connectDB()
	return nil
}

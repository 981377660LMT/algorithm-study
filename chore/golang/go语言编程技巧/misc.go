package main

import (
	rand2 "crypto/rand"
	"errors"
	"fmt"
	rand1 "math/rand"
	"time"
)

func main() {
	// sliceToArray()
	// multiError()
	// avoidNakedParameters()
	// numericSeparators()
	// avoidMathRand()
	// nilSliceFirst()
}

// !sliceToArray
func sliceToArray() {
	a := []int{1, 2, 3, 4, 5}
	b := [3]int(a[0:3]) // Go 1.20 后支持
	fmt.Println(b)
}

// !Underscore Import 下划线导入
// 它会在不创建那个包的引用的情况下，执行那个包里的初始化代码（init()函数）。

// !错误封装 Wrapping Errors
func multiError() {
	err := errors.Join(fmt.Errorf("err1"), fmt.Errorf("err2"))
	fmt.Println(err)
}

// !编译时接口检查 Compile-Time Interface Verification
// var _ MyImpl = (*MyInterface)(nil)

// !避免裸露参数 Avoid Naked Parameters
func avoidNakedParameters() {
	var f func(string, bool, bool)
	f("hello", true, false)                            // bad
	f("hello", true /* verbose */, false /* dryRun */) // good
}

// !数字分隔符 Numeric separators
// !快速识别数字有几位
func numericSeparators() {
	const OneMillion = 1_000_000
	const OneBillion = 1_000_000_000
	const OneTrillion = 1_000_000_000_000
}

// !使用crypto/rand生成密钥，避免使用math/rand.
// Avoid using math/rand, use crypto/rand for keys instead.
//
// math/rand这个包生成的是伪随机数.
// 这意味着如果你知道那些数字是怎么生成的（就是知道用于生成随机数序列的种子），那你就能预知到会生成哪些数字.
// crypto/rand提供了一个生成密码学安全随机数的方式.
func avoidMathRand() {
	key1 := func() string {
		r := rand1.New(rand1.NewSource(time.Now().UnixNano()))
		buf := make([]byte, 16)
		for i := range buf {
			buf[i] = byte(r.Intn(256))
		}
		return fmt.Sprintf("%x", buf)
	}

	key2 := func() string {
		buf := make([]byte, 16)
		_, err := rand2.Read(buf)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%x", buf)
	}

	fmt.Println(key1())
	fmt.Println(key2())
}

// !优先使用nil切片而不是空切片
// Empty slice or, even better, NIL SLICE.
//
// !nil切片没有分配任何的内存，而空切片（[]int{}）则分配了很小的内存去指向一个空数组.
//
// 在设计代码的时候，你应该同等对待非空切片、空切片和nil切片
// !在使用JSON的时候，nil切片和空切片的表现是不一样的。
func nilSliceFirst() {
	var s1 []int
	s2 := []int{}

	fmt.Println(s1 == nil) // true
	fmt.Println(s2 == nil) // false
}

// !错误信息不要大写或者以标点结尾
// 为什么要小写？
// !错误信息经常会被封装或者合并到其他错误信息里。
// 如果一条错误信息以大写字母开头，那么当它出现在句子中间的时候，看起来就会很怪或者显得格格不入。

// !什么时候使用空白导入和点导入？
// Blank imports are used when you want to import a package solely for its side-effects.
//  一个常见的例子是在使用 database/sql 包的程序中导入数据库驱动包。
//  数据库驱动包被导入是因为其副作用（例如，将自己注册为 database/sql 的驱动）。

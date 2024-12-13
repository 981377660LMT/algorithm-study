package main

import (
	"context"
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
// Dot imports are used when you want to use a package's exported identifiers without a qualifier.
//  这种形式在测试中特别有用。
//  !尤其是在处理难以轻易解决的循环依赖时。

// !不要通过返回 -1 或者 nil 来表示错误，这被称为“带内错误”（in-band errors）。
// 使用 in-band 错误的主要问题是，需要调用者记住每次都要检查返回的特殊值。
// 但这是...非常容易出错的。
// Go的解决方案是：多返回值
//
// 现在，您的代码便拥有了 3 个优势（您甚至不需要额外关心）：
// 明确的关注点分离
// 返回值明确的表示了的哪部分是实际结果，哪部分表示操作的成功或失败。
// 强制错误处理
// Go 编译器要求开发人员处理错误的可能性，从而降低忽略错误的风险（因此，请勿使用 “_” 来忽略错误）。
// 提高可读性和可维护性
// 代码可以明确地解释自身的行为（documents itself）。

// !“尽快返回、尽早返回”，避免代码嵌套。"Return fast, return early" to avoid nested code.
// 当你写代码的时候，你会想让它尽可能的清晰易懂。
// !要做到这点，其中一个方法就是组织你的代码，让它的“快乐路径”（预期的或者正常的执行流程）更加的突出和简单明了。
// 所以，指导原则是什么？
// !很简单：提前处理错误，别让他们碍事。

// !在使用者的包中定义接口，而不是提供者的包中定义. `Define interfaces in the consumer package`, not the producer.
// 在使用者的包中定义接口，而不是提供者的包中定义；在提供者的包中使用具体类型作为返回值；避免过早定义接口；只在有明确使用场景下定义接口，确保它们是必要的且设计得当的。
// !通过在使用接口的地方(在高级模块中)定义接口，可以确保这些模块依赖于抽象接口而不是具体的实现。
// golang 的隐式接口真方便.

// !除非出于文档说明需要，否则避免使用命名结果. Avoid named results unless necessary for documentation..
// 哪些场景是需要的呢？
// 1. 延迟闭包
// 2. 短小的函数

// !传递值，而不是指针 Pass values, not pointers.
// https://colobu.com/gotips/020.html
// 传递值的速度很快，而且很少比传递指针慢
// 这可能会因为复制而显得有悖常理，但原因如下：
// 	 - 复制少量数据非常高效，通常比使用指针时所需的间接操作更快。
// 	 - 当值直接传递时，垃圾收集器的工作量会减少，因为它需要跟踪的指针引用更少。
// 	 - 通过值传递的数据在内存中往往存储得更紧密，这使得CPU能够更快地访问数据。
// 你很少会遇到一个足够大的结构体，以至于通过指针传递对其有利。

// !定义方法时，优先使用指针作为接收器(receiver). Prefer using a pointer receiver when defining methods.
// 何时适合使用值接收器？
//
// 小型且不会被改变的类型
// 如果你的类型是 map、func、channel，或者涉及到切片，而且切片的大小和容量不会改变（尽管元素可能会改变）
// !避免在给定结构体中混合使用不同的接收器类型，以保持一致性。

// !使用Builder模式、optionFunc模式构建复杂对象

// !省略 getter 方法的'Get'前缀
// 在编写代码时，我们通常以动词开头给函数命名，比如 get、set、fetch、update、calculate 等等...
// 但是在Go语言中 getter 方法是一个例外。
// !在 Go 语言 中，封装是通过方法的可见性和命名约定来实现的，这巧妙地支持了封装，而不需要严格使用 getter/setter 方法。
// !然而，如果需要额外的逻辑，或者我们想要访问一个计算字段，手动定义 getter 和 setter 方法也是没有什么问题的。

// !在 goroutines 之间进行信号传递时，使用 'chan struct{}' 而不是 'chan bool'
// 减少了歧义.

// !避免使用 context.Background()，使你的协程具备承诺性
// !协程给了使用方一个承诺：我要么执行成功，要么因为超时等原因取消执行，但最终在有限时间内一定会有一个明确的状态。
// 一般来说，有两种方法可以使你的协程具有承诺性：取消和超时。
// 我要么完成我的任务，要么及时告诉你为什么我不能完成，并且你可以在任何时候取消我的任务。
// 如果想让某些任务不受上下文的影响，可以使用 context.WithoutCancel(ctx)。
func avoidContextBackground() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// do something with ctx
	_ = ctx
}

// !尽量不要使用panic
// 1.只有panic触发和调用recover()是在同一个goroutine中时，使用recover()才有效。
// 因此，在主函数中的defer函数无法捕获或恢复panic，尽管尝试了恢复，程序仍然会崩溃。
// 2.在生产环境中，代码必须具备极高的稳健性 程序意外崩溃是绝对要避免的，因为它会导致系统宕机
// 3.系统中某一部分的panic可能会引发连锁反应 这可能导致系统（尤其是在微服务或者分布式系统）中其他部分接连出现故障（可能是级联失败）。
// 当程序返回错误而非panic时，你的程序可以根据错误做进行相应的处理，例如：
//
// - 重试操作
// - 使用默认值
// - 记录详细的调试信息
// - 程序终止
// - 等等...
//
// 这种灵活性对于构建健壮的系统至关重要。
// 应当把panic作为最后的手段
// !仅在遇到真正无法恢复的错误时才使用panic，即如果继续运行程序可能会引发更严重的问题，比如数据损坏或未知行为。
// !在程序初始化阶段，如果一个关键组件启动失败，panic或许是“可接受的”，因为它表明程序无法按预期运行。
// !参考Java规范：只有mustXXX和initXXX的函数才允许panic.

// !Tip #34 以context开头，以options结尾，并且总是用error来关闭
// Lead with context, end with options, and always close with an error
// - context 通常与一个请求或者操作的生命周期息息相关 -> 取消、截止、上下文；这种一致性有助于提高代码的可读性，并且让代码库变得易于浏览。
//   !此外，别把context.Context放到 struct 中。
//   context 本质上意味着它注定是短暂的，旨在贯穿于一段程序，而非成为对象状态的一部分（这里有一些例外情况，比如，HTTP 的 handler，大家习惯地从请求中提取 context，这是因为 context 早已跟请求的生命周期相关联了）。
// - options 结构体置于最后
//   !参数的顺序可能表明了这个参数的重要性。
//   把这个结构体作为一个函数的最后一个参数有两个目的：
//   保持一致性（与可变参 options 模式一致）
//   表明这些是可选配置项，而非函数操作逻辑的核心部分
// - 以 error（或者bool）结尾
//   Go习惯通过最后一个返回值(通常是一个error)来表明一个操作是成功还是失败。如果兼而有之，那么优先级应当是（x, bool, error）。
//   ```go
//   var func FetchData(ctx context.Context, url string, opts FetchOptions) (data []byte, err error)
//   var func TryFetchData(ctx context.Context, url string, opts FetchOptions) (data []byte, ok bool, err error)
//   ```

// Tip #35 转换字符串时优先使用 strconv 而非 fmt

// !使用空结构体模拟Symbol(Using Unexported Empty Struct as Context Key)
// !这背后的原理归结为 Go 如何比较 interface{}，只有当两个 interface{} 的类型和值都匹配时，它们才相等。
// 一个空结构体不会分配内存，它没有字段因而不包含数据，但它的类型仍然可以唯一地标识上下文值。

// !Tip #38 使用 fmt.Errorf 使你的错误信息清晰明了，不要让它们过于赤裸

// !Tip #39 避免在循环中使用defer，否则可能会导致内存溢出(Avoid defer in loops, or your memory might blow up.)
// https://colobu.com/gotips/039.html
// 最好提取出一个函数.

// !Tip #40 在使用defer时处理错误以防止静默失败

// !Tip #41 将你结构体中的字段按从大到小的顺序排列

// !Tip #42 单点错误处理，降低噪音。Single Touch Error Handling, Less Noise.

// !Tip #43 优雅关闭你的应用程序

// !Tip #44 有意地使用Must函数来停止程序 Intentionally Stop with Must Functions
// 这类函数有一个特定的命名模式，它们以“Must”（或“must”）开头，这就是提醒你需要警惕一下，如果程序没有按照预期执行的话就会导致panic。
//
// Must函数主要用于：
// 1. 不应该失败的初始化任务
//    在应用程序开始时设置包级变量、设置正则表达式、连接数据库等。
// 2. 用于测试的辅助函数
//    使用t.Fatal/Fatalf立即失败这个测试用例

// !Tip #45 始终管理您协程的生命周期
// Golang 中的协程是有栈协程，这意味着相比较于其他语言中的类似结构，Golang 中的协程会占用更多的内存，每个 Golang 协程至少会占用 2KB 的内存。
// !不要小看这 2KB 的内存占用量，因为在 Gloang 中，协程的创建是非常便捷的，很容易就快速增长到一个庞大的数量，当协程数量达到 10K 时，其内存占用将达到 20MB。
// 因此，对于那些本质上没有明确终点的任务（例如：网络连接服务、配置文件监视等），应该使用取消信号或条件来明确定义这些任务何时应该结束。

// !Tip #46 避免在 switch 语句的 case 中使用 break，除非与标签一起使用。
// !Go 的 switch 语句中的每个 case 自带一个隐式的 break。

// !Tip #47 表驱动测试，测试集和并行运行测试.Table-driven tests, subtests, and parallel tests.
// 如果有一个测试失败，我不想运行其余的测试，因为会很慢。 => 我们可以使用 t.Fatalf 而不是 t.Errorf，它相当于 t.Logf + t.FailNow

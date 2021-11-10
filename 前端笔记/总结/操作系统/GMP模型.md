Golang 的协程调度器原理及 GMP 设计思想

1. gorutine 与 js/python 里的协程 corutine 的 区别
   js/Python 的协程是用生成器实现的，上下文切换通过操作栈帧来实现，单个生成器对象只持有栈帧而不是整个调用栈。
   如今协程已经成为大多数语言的标配，例如 Golang 里的 goroutine，JavaScript 里的 async/await。尽管名称可能不同，但它们都可以被划分为两大类，一类是有栈（stackful）协程，例如 goroutine；一类是无栈（stackless）协程，例如 async/await。
   **有栈（stackful）协程与无栈（stackless）协程**
   所谓的有栈，无栈并不是说这个协程运行的时候有没有栈，而是说协程之间是否存在调用栈（callbackStack）
   **有栈协程**
   实现一个协程的关键点在于如何保存、恢复和切换上下文。已知函数运行在调用栈上；如果将一个函数作为协程，我们很自然地联想到，保存上下文即是保存从这个函数及其嵌套函数的（连续的）栈帧存储的值，以及此时寄存器存储的值；恢复上下文即是将这些值分别重新写入对应的栈帧和寄存器；而切换上下文无非是保存当前正在运行的函数的上下文，恢复下一个将要运行的函数的上下文。有栈协程便是这种朴素思想下的产物。
   **无栈协程**
   相比于有栈协程直接切换栈帧的思路，无栈协程在**不改变函数调用栈**的情况下，采用类似生成器（generator）的思路实现了上下文切换
   尽管有栈协程和无栈协程是根据它们存储上下文的机制区分命名的，但二者的本质区别还是在于**是否可以在其任意嵌套函数中被挂起**。这也决定了有栈协程被挂起时的自由度要比无栈协程高。比如使用无栈协程的 JavaScript 就不能这么写：
   ```JS
    async function processArray(array) {
        // 显然这里 forEach 是个嵌套函数
        array.forEach(item => {
            // Uncaught SyntaxError:
            // await is only valid in async function
            const result = await processItem(item)
            ...
        })
    }
   ```
   但使用有栈协程的 Golang 就可以轻松实现类似的逻辑：
   ```Go
    func processArray(array []int) {
        for i := 0; i < len(array); i++ {
            ch := make(chan int)
            go processItem(array[i], ch)
            result := <- ch
            ...
        }
    }
   ```
2. 为什么 Java 坚持多线程不选择协程？
   JCP 根本难辞其咎

3. Fiber/coroutine
   纤程（Fiber）就是协程（coroutine)，也叫做用户级线程或者轻线程之类的。
   纤程的概念中有两个关键点：

   - 纤程拥有独立的栈空间和寄存器环境；
   - 纤程在用户态实现调调度，也就是说完全由程序员控制；

   协程的特点

   1. 协程可以自动让出 CPU 时间片。注意，不是当前线程让出 CPU 时间片，而是线程内的某个协程让出时间片供同**线程内其他协程运行**
   2. 协程可以恢复 CPU 上下文。当另一个协程继续执行时，其需要恢复 CPU 上下文环境
   3. 协程有个管理者，管理者可以选择一个协程来运行，其他协程要么阻塞，要么 ready，或者 died
   4. 运行中的协程将占有当前线程的所有计算资源
   5. 协程天生有栈属性，而且是 lock free

4. python 和 js 异步对比
   跟 js 的异步对比就很明显, js 是为网页而生, 必须异步才能不影响 UI 线程, 所以 js 天生就是异步的, js 的运行时都有一个内置的 event loop, 这是异步代码的运行基础. 正常情况下, js 的代码是不能脱离 event loop 运行的, 原生 api 甚至不存在同步的 io.
   核心的差异在于, **js 默认运行 IO 不阻塞的异步代码, python 默认运行 IO 阻塞的同步代码.**

5. goroutine
   协程跟线程是有区别的，线程由 CPU 调度是抢占式的，协程由用户态调度是协作式的，一个协程让出 CPU 后，才执行下一个协程。
   goroutine 来自协程的概念，让一组可复用的函数运行在一组线程之上，**即使有协程阻塞，该线程的其他协程也可以被 runtime 调度，转移到其他可运行的线程上。**
6. goroutine 调度器:GMP 模型
   G:goroutine(go 协程)
   M:machine(线程)
   P:Processor(处理器),它包含了运行 goroutine 的资源、G 队列
   在 Go 中，线程是运行 goroutine 的实体，调度器的功能是把可运行的 goroutine 分配到工作线程上。

   **调度器的设计策略**

   - 复用线程：避免频繁的创建、销毁线程，而是对线程的复用。

   1）work stealing 机制

   ​ 当本线程无可运行的 G 时，尝试从其他线程绑定的 P 偷取 G，而不是销毁线程。

   2）hand off 机制

   ​ 当本线程因为 G 进行系统调用阻塞时，线程释放绑定的 P，把 P 转移给其他空闲的线程执行。

   - 利用并行：GOMAXPROCS 设置 P 的数量，最多有 GOMAXPROCS 个线程分布在多个 CPU 上同时运行。GOMAXPROCS 也限制了并发的程度，比如 GOMAXPROCS = 核数/2，则最多利用了一半的 CPU 核进行并行。

   - 抢占：在 coroutine 中要等待一个协程主动让出 CPU 才执行下一个协程，在 Go 中，一个 goroutine 最多占用 CPU 10ms，防止其他 goroutine 被饿死，这就是 goroutine 不同于 coroutine 的一个地方。

   - 全局 G 队列：在新的调度器中依然有全局 G 队列，但功能已经被弱化了，当 M 执行 work stealing 从其他 P 偷不到 G 时，它可以从全局 G 队列获取 G。

**go func() 调度流程**

7. Go 调度器调度场景过程全解析
   Go 调度本质是把大量的 goroutine 分配到少量线程上去执行，并利用多核并行，实现更强大的并发。

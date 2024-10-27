React Fiber架构如何从JS引擎手中“夺回”调度权
https://segmentfault.com/a/1190000020110045
https://www.zhihu.com/question/388457689
https://www.zhihu.com/question/388457689/answer/1165782397
https://segmentfault.com/a/1190000020110166
https://github.com/reactjs/rfcs/pull/68#issuecomment-433158179
https://github.com/facebook/react/issues/7942

- React 16采用新的Fiber架构对React进行完全重写，同时保持向后兼容。

- concurrent rendering
  又叫 async rendering，包含几个feature：

  - Time Slicing
    为了保证60fps的刷新率，React将渲染工作切分，当这一帧的时间片用完，`React框架就将控制权交还给浏览器`，剩下的渲染工作等到下一帧再做。
  - Suspense
    如果一个组件的渲染要等待一个异步任务的完成（比如data-fetching、或者lazy-loading），那么`React组件可以将控制权交给React框架`，React框架会等待异步任务完成以后继续渲染这个组件。
    组件通过throw Promise的方式来将控制权交还调度器。
  - React Hooks
    允许`React组件将控制权交给React框架`，从而React能够帮你维护、注入state、context等。只不过目前React hooks的控制权交接是同步完成的（也就是hooks执行一些同步代码以后就将控制权交回到你的组件代码），因此看上去就是一个函数调用而已，但是它的心智模型实际上是Algebraic Effects （Suspense的心智模型也是）。

  实现这几个feature的关键前提是：
  **React的渲染能够被中止（interrupt）、恢复**
  为此，我们需要fiber架构。
  因为JS的协程(async/await)是无栈协程，不好(原因在下面讲)，React 实现了一个有栈协程，叫做Fiber。

- React Fiber架构：可控的“调用栈”
  `相当于手写一个有栈协程。`
  React Fiber架构的核心能力就是让程序控制权能够在React开发者、React框架、浏览器之间合理切换(调度)。

  - 之前的React版本中，调用栈是递归的，一旦开始渲染，就无法中断。
  - 一个fiber 类似一个函数栈帧frame，一个组件实例对应一个fiber。

  | 函数栈帧     | Fiber    |
  | ------------ | -------- |
  | 返回指针地址 | 父组件   |
  | 当前函数     | 当前组件 |
  | 调用参数     | props    |
  | 局部变量     | state    |

  https://github.com/facebook/react/blob/c1d3f7f1a97adad9441287a92dcd4ac5d2478c38/packages/react-reconciler/src/ReactFiber.js#L251

  fiber 有三个指针: child, sibling, return，有这三个指针的链表就能够实现深度优先遍历
  https://github.com/facebook/react/issues/7942

  fiber与调用栈的另一个区别是，栈帧在函数返回以后就销毁了，而`fiber会在渲染结束以后继续存在，保存组件实例的信息（比如state）。`

  React能够在这个执行模型之上实现更多的编程模式，`而不必受到过多来自JavaScript语言的约束：`
  例如：

  - 没有函数着色问题(async 的传染性)
  - 随时暂停、恢复渲染
  - 并发渲染：一个渲染还没有完成就开始另一个渲染，并控制各个渲染的优先级
  - Suspense
  - Algebraic Effects，以及基于它的React hooks编程模式

- 为什么不使用generator来实现协作式调度

  1. 性能问题
     您必须`将每个函数包装在生成器中`。这不仅增加了大量语法开销，而且增加了任何现有实现中的运行时开销。
     https://www.typescriptlang.org/zh/play/?target=1&module=7#code/IYZwngdgxgBAZgV2gFwJYHsI2AB1QCgEoYBvAWACgYYAnAU2QRqwEZKBfSy0SWRFDFgC2wVBCKlK1KJhDJadEDAC82AO6j5uAoSkwZEEOgA2dAHTH0Ac3z0QuipwpA
  2. generator只能从上次yield的地方恢复，不能回到更早的执行状态
     而Fiber架构能够更加灵活地让React从任意一个Fiber恢复执行（不只是从上次中断的地方恢复，而且能够从更早的Fiber恢复）

     ```js
     function* doWork(a, b, c) {
       var x = doExpensiveWorkA(a)
       yield
       var y = x + doExpensiveWorkB(b)
       yield
       var z = y + doExpensiveWorkC(c)
       return z
     }
     ```

     如果用yield，假设现在处于第二个yield，发现doExpensiveWorkB的计算结果需要更新（比如组件B内部发生了setState，需要重新渲染），`此时你只能同时丢弃A和B的计算结果，从头开始。注意，此时doExpensiveWorkA的计算被白白浪费了。`
     而如果使用纯函数+memo（即fiber选择的方案）：
     如果需要交出控制权（比如发现时间片用完），直接throw一个信号即可退出，然后下次重新执行doWork。由于React频繁使用memoization，因此计算浪费可以忽略不计。

     ```js
     function doWork(a, b, c) {
       var x = withMemo(doExpensiveWorkA, a)
       var y = x + withMemo(doExpensiveWorkB, b)
       var z = y + withMemo(doExpensiveWorkC, c)
       return z
     }
     ```

     Generator的状态是执行引擎管理的，React必须服从；而doWork的memoization是可以被React完全掌控的。

  3. generator 传染性，给用户增加心智负担
  4. typescript对于generator函数的类型推导能力很弱

- 只有Render Phase是可以打断的
  React Fiber的更新过程被分为两个阶段(Phase)：第一阶段Render Phase；第二阶段Commit Phase。
  在第一阶段Render Phase，`React Fiber会找出需要更新哪些DOM，这个阶段是可以被打断的，甚至可以中途放弃；`但是到了第二阶段Commit Phase，就必须一鼓作气把DOM更新完，绝不会被打断。

- Algebraic Effects(代数效应)，以及它在React中的应用
  最近的一些新兴研究型语言提供了一种叫Algebraic Effects的`函数远程通信机制`。React团队借鉴过这个概念。
  https://segmentfault.com/a/1190000020110166

  类似async/await ，但关键区别在于，

  1. async/await 耦合了（写死了）副作用方法，而代数效应则是动态的。
  2. 异步性会`感染`所有上层调用者，而代数效应则不会。

  `在实际实现中，通过 throw promise 来模拟 Algebraic Effect`
  React团队将hooks都看做Algebraic Effect，React调度器提供了各种基本hooks的“effect handler”，在调用render函数的时候为其提供上下文信息。

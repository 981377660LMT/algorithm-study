https://juejin.cn/post/6914936575507988494
https://juejin.cn/post/7217487697676697655
https://github.com/wzhudev/blog/issues/9

Angular 通过引入 ZoneJS 使得其变更检测机制更加简单与可靠
在 ZoneJS 中有一个核心概念：Zone（域）。一个 Zone 表示一个 JavaScript `执行过程的上下文`，其可以在异步任务之间进行持久性传递。
Zone 库附带的代码会对所有浏览器的异步 API 进行 monkey patch，并将它们重定向通过 Zone 进行拦截。
Zone 允许拦截调度和调用异步操作，并在异步任务之前和之后执行附加代码。
拦截规则使用 [ZoneConfig] 进行配置。 系统中可以有许多不同的 Zone 实例，但在任何给定时间只有一个 Zone 是活动的，可以使用 [Zone#current] 进行检索。
如果将一系列异步操作视为执行线程（有点牵强），那么 [Zone#current] 相当于 thread-local 变量.

能力：

- thread-local
- 异步任务、错误跟踪

---

[Angular 框架解读--Zone 区域之 zone.js](https://godbasin.github.io/2021/05/01/angular-design-zonejs/)

zone.js 是针对异步编程提出的作用域解决方案

- 提供异步操作之间的执行上下文
- 提供异步生命周期挂钩
- 提供统一的异步错误处理机制

zone 具有当前区域的概念：当前区域是随所有异步操作一起传播的异步上下文，它表示与当前正在执行的堆栈帧/异步任务关联的区域。

- zone.fork(zoneSpec): 创建一个新的子区域，并将其 parent 设置为用于分支的区域
- zone.run(callback, ...)：在给定区域中同步调用一个函数
- zone.runGuarded(callback, ...)：与 run 捕获运行时错误相同，并提供了一种拦截它们的机制。如果任何父区域未处理错误，则将其重新抛出。
- zone.wrap(callback)：产生一个新的函数，该函数将区域绑定在一个闭包中，并在执行 zone.runGuarded(callback)时执行，与 JavaScript 中的 Function.prototype.bind 工作原理类似。

---

zone.js 是如何`识别出异步任务`的呢？
其实 zone.js 主要是通过猴子补丁`拦截异步 API`（包括 DOM 事件、XMLHttpRequest 和 NodeJS 的 API 如 EventEmitter、fs 等）来实现这些功能

```ts
// 为指定的本地模块加载补丁
static __load_patch(name: string, fn: _PatchFn, ignoreDuplicate = false): void {
  // 检查是否已经加载补丁
  if (patches.hasOwnProperty(name)) {
    if (!ignoreDuplicate && checkDuplicate) {
      throw Error('Already loaded patch: ' + name);
    }
  // 检查是否需要加载补丁
  } else if (!global['__Zone_disable_' + name]) {
    const perfName = 'Zone:' + name;
    // 使用 performance.mark 标记时间戳
    mark(perfName);
    // 拦截指定异步 API，并进行相关处理
    patches[name] = fn(global, Zone, _api);
    // 使用 performance.measure 计算耗时
    performanceMeasure(perfName, perfName);
  }
}
```

---

任务执行的生命周期
zone.js 提供了异步操作生命周期钩子，有了这些钩子，Zone 可以监视和拦截异步操作的所有生命周期：

```ts
interface ZoneSpec {
  // 允许拦截 Zone.fork，对该区域进行 fork 时，请求将转发到此方法以进行拦截
  onFork?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    zoneSpec: ZoneSpec
  ) => Zone
  // 允许拦截回调的 wrap
  onIntercept?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    delegate: Function,
    source: string
  ) => Function
  // 允许拦截回调调用
  onInvoke?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    delegate: Function,
    applyThis: any,
    applyArgs?: any[],
    source?: string
  ) => any
  // 允许拦截错误处理
  onHandleError?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    error: any
  ) => boolean
  // 允许拦截任务计划
  onScheduleTask?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    task: Task
  ) => Task
  // 允许拦截任务回调调用
  onInvokeTask?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    task: Task,
    applyThis: any,
    applyArgs?: any[]
  ) => any
  // 允许拦截任务取消
  onCancelTask?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    task: Task
  ) => any
  // 通知对任务队列为空状态的更改
  onHasTask?: (
    parentZoneDelegate: ZoneDelegate,
    currentZone: Zone,
    targetZone: Zone,
    hasTaskState: HasTaskState
  ) => void
}
```

这些生命周期的钩子回调会在 zone.fork()时，通过 new Zone()创建子区域并创建和传入到 ZoneDelegate 中：

```ts
class Zone implements AmbientZone {
  constructor(parent: Zone | null, zoneSpec: ZoneSpec | null) {
    this._zoneDelegate = new ZoneDelegate(
      this,
      this._parent && this._parent._zoneDelegate,
      zoneSpec
    )
  }
}

class ZoneDelegate implements AmbientZoneDelegate {
  constructor(zone: Zone, parentDelegate: ZoneDelegate | null, zoneSpec: ZoneSpec | null) {
    // 管理 onFork 钩子回调
    this._forkZS = zoneSpec && (zoneSpec && zoneSpec.onFork ? zoneSpec : parentDelegate!._forkZS)
    this._forkDlgt = zoneSpec && (zoneSpec.onFork ? parentDelegate : parentDelegate!._forkDlgt)
    this._forkCurrZone = zoneSpec && (zoneSpec.onFork ? this.zone : parentDelegate!._forkCurrZone)
  }
  // fork 调用时，会检查是否有 onFork 钩子回调注册，并进行调用
  fork(targetZone: Zone, zoneSpec: ZoneSpec): AmbientZone {
    return this._forkZS
      ? this._forkZS.onFork!(this._forkDlgt!, this.zone, targetZone, zoneSpec)
      : new Zone(targetZone, zoneSpec)
  }
}
```

---

# Zone.js 源码简读

两种 Patch 方式：Wrap 和 Task
对于 Api 的 Patch 方式不同，控制的颗粒度是不同的：

- Wrap 方式：onInvoke 和 onIntercept
- Task 方式：onScheduleTask（Zone 内配置了 Task 就会触发）， onInvokeTask（Task 任务触发前）， onCancelTask（Task 任务取消前）， onHasTask(Zone 内 Task 状态变化后触发)

如果想尝试更多的新功能，需要单独引入 patch

Delegate 类定义了钩子函数的执行规则：冒泡。
fork 时传入 ZoneSpec 的话，parent 的 Delegate 就会被保存。

---

# zone.js 由入门到放弃之一——通过一场游戏认识 zone.js

A Zone is an execution context that persists across async tasks.
You can think of it as thread-local storage for JavaScript VMs.

`每个子 zone 都保存了其父 zone 的引用；每个父 zone 也能监听到子 zone 的事件。`

每个子 zone 都保存了其父 zone 的引用这个好理解，那么每个父 zone 也能监听到子 zone 的事件怎么理解？其实这个就是 zone.js 最神奇的地方，zone.js 在初始化的时候对很多 API 都做了“手脚”——Monkey Patch，将这些异步方法封装成了 zone.js 中的异步任务。同时，由于在这些任务中定义很多勾子函数，导致 zone.js 可以完全监控这些异步任务的整个生命周期。

# zone.js 由入门到放弃之二——zone.js API 大练兵

zone.js 中最重要的三个定义为：Zone，ZoneDelegate，ZoneTask。
搞清楚了这三个类的 API 及它们之间关系，基本上对 zone.js 就通了。
而 Zone，ZoneDelegate，ZoneTask 三者中，Zone，ZoneDelegate 其实半差不差的可以先当成一个东西。

```ts
interface Zone {
  // 通用API
  name: string

  get(key: string): any

  getZoneWith(key: string): Zone | null

  // 给当前Zone创建一个子Zone，函数接受一个ZoneSpec的参数，参数规定了当前Zone的一些基本信息以及需要注入的钩子
  fork(zoneSpec: ZoneSpec): Zone

  run<T>(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): T

  runGuarded<T>(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): T

  runTask(task: Task, applyThis?: any, applyArgs?: any): any

  cancelTask(task: Task): any

  // Wrap类包装API
  wrap<F extends Function>(callback: F, source: string): F

  // Task类包装API
  scheduleMicroTask(
    source: string,
    callback: Function,
    data?: TaskData,
    customSchedule?: (task: Task) => void
  ): MicroTask

  scheduleMacroTask(
    source: string,
    callback: Function,
    data?: TaskData,
    customSchedule?: (task: Task) => void,
    customCancel?: (task: Task) => void
  ): MacroTask

  scheduleEventTask(
    source: string,
    callback: Function,
    data?: TaskData,
    customSchedule?: (task: Task) => void,
    customCancel?: (task: Task) => void
  ): EventTask

  scheduleTask<T extends Task>(task: T): T
}
```

Zone 中的 API 大致分了三类：通用 API、Wrap 类和 Task 类。
Wrap 和 Task 类分别对应 zone.js 对异步方法的两种打包方式（Patch），不同的打包方式对异步回调提供了不同粒度的"监听"方式，即不同的打包方式会暴露出不同的拦截勾子。
你可以根据自身对异步的控制精度选择不同的打包方式。

- Wrap 方式：

onIntercept：当在注册回调函数时被触发，简单点理解在调用 wrap 的时候，该勾子被调用
onInvoke：当通过 wrap 包装的函数调用时被触发

- Task 方式

- onScheduleTask：在调用真正的 setTimeout 之前会触发 onScheduleTask
- onInvokeTask：Task 会将 setTimeout 回调通过 wrap 打包，当回调被执行之前，onInvokeTask 勾子会被触发
- onHasTask：它记录了任务队列的状态。当任务队列中有 MacroTask、MicroTask 或 EventTask 进队或出队时都会触发该勾子函数
- onCancelTask

runXXX 方法可以指定函数运行在特定的 zone 上，这里可以把该方法类比成 JS 中的 call 或者 apply，它可以指定函数所运行的上下文环境；而 zone 在这里可以类比成特殊的 this，只不过 zone 上下文可以跨执行栈保存，而 this 不行。与此同时，runXXX 在回调执行结束后，会自动地恢复 zone 的执行环境。

```ts
run(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): T;
runGuarded(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): T;
runTask(task: Task, applyThis?: any, applyArgs?: any；
```

`runXXX 方法类似于 call 和 apply 的作用，那么 wrap 方法类似于 JS 中的 bind 方法`
wrap 可以将执行函数绑定到当前的 zone 中，使得函数也能执行在特定的 zone 中。下面是我简化以后的 wrap 源码

```ts
public wrap<T extends Function>(callback: T, source: string): T {
  // 省略若干无关紧要的代码
  const zone: Zone = Zone.current;
  return function() {
    return zone.runGuarded(callback, (this as any), <any>arguments, source);
  } as any as T;
}
```

runGuarded 对比 run 就多了一个 catch

```ts
public runGuarded<T>(
    callback: (...args: any[]) => T, applyThis: any = null, applyArgs?: any[],
    source?: string) {
  _currentZoneFrame = {parent: _currentZoneFrame, zone: this};
  try {
    try {
      return this._zoneDelegate.invoke(this, callback, applyThis, applyArgs, source);
    } catch (error) {
      if (this._zoneDelegate.handleError(this, error)) {
        throw error;
      }
    }
  } finally {
    _currentZoneFrame = _currentZoneFrame.parent!;
  }
}
```

## ZoneTask

Task 形式比 Wrap 形式有更丰富的生命周期勾子，使得你可以更精细化地控制每个异步任务。好比 Angular，它可以通过这些勾子决定在何时进行脏值检测，何时渲染 UI 界面。
zone.js 任务分成 MacroTask、MicroTask 和 EventTask 三种：

- MicroTask：在当前 task 结束之后和下一个 task 开始之前执行的，`不可取消`，如 Promise，MutationObserver、process.nextTick
- MacroTask：一段时间后才执行的 task，`可以取消`，如 setTimeout, setInterval, setImmediate, I/O, UI rendering
- EventTask：监听未来的事件，可能执行 0 次或多次，执行时间是不确定的

# zone.js 由入门到放弃之三——zone.js 源码分析【setTimeout 篇】

Zone 其实主要只负责两件事：

- `维护 Zone 的上下文栈`：我们知道 Zone 是个具有继承关系的链式结构。zone.js 在全局会维护一个 Zone 栈帧，每当我们在某个 Zone 中执行代码时，Zone 要负责将当前的 Zone 上下文置于栈帧中；当代码执行完毕，又要负责将 Zone 栈帧恢复回去。
- `Zone还负责ZoneTask的状态切换`。上文说过，Zone 可以对宏任务、微任务、事件进行管理。那么每个任务在 Zone 中处于何种阶段、何种状态也是由 Zone 负责的。Zone 会在适当时候`调用 ZoneTask 的_transitionTo 方法切换 ZoneTask 的状态`。

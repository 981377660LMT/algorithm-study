https://juejin.cn/post/6914936575507988494
https://juejin.cn/post/7217487697676697655
https://github.com/wzhudev/blog/issues/9

Angular 通过引入 ZoneJS 使得其变更检测机制更加简单与可靠
在 ZoneJS 中有一个核心概念：Zone（域）。一个 Zone 表示一个 JavaScript `执行过程的上下文`，其可以在异步任务之间进行持久性传递。
Zone 库附带的代码会对所有浏览器的异步 API 进行 monkey patch，并将它们重定向通过 Zone 进行拦截。
Zone 允许拦截调度和调用异步操作，并在异步任务之前和之后执行附加代码。
拦截规则使用 [ZoneConfig] 进行配置。 系统中可以有许多不同的 Zone 实例，但在任何给定时间只有一个 Zone 是活动的，可以使用 [Zone#current] 进行检索。
如果将一系列异步操作视为执行线程（有点牵强），那么 [Zone#current] 相当于 thread-local 变量.

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

````ts
interface ZoneSpec {
  // 允许拦截 Zone.fork，对该区域进行 fork 时，请求将转发到此方法以进行拦截
  onFork?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, zoneSpec: ZoneSpec) => Zone;
  // 允许拦截回调的 wrap
  onIntercept?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, delegate: Function, source: string) => Function;
  // 允许拦截回调调用
  onInvoke?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, delegate: Function, applyThis: any, applyArgs?: any[], source?: string) => any;
  // 允许拦截错误处理
  onHandleError?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, error: any) => boolean;
  // 允许拦截任务计划
  onScheduleTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task) => Task;
  // 允许拦截任务回调调用
  onInvokeTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task, applyThis: any, applyArgs?: any[]) => any;
  // 允许拦截任务取消
  onCancelTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task) => any;
  // 通知对任务队列为空状态的更改
  onHasTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, hasTaskState: HasTaskState) => void;
}```
````

这些生命周期的钩子回调会在 zone.fork()时，通过 new Zone()创建子区域并创建和传入到 ZoneDelegate 中：

```ts
class Zone implements AmbientZone {
  constructor(parent: Zone | null, zoneSpec: ZoneSpec | null) {
    this._zoneDelegate = new ZoneDelegate(this, this._parent && this._parent._zoneDelegate, zoneSpec)
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
    return this._forkZS ? this._forkZS.onFork!(this._forkDlgt!, this.zone, targetZone, zoneSpec) : new Zone(targetZone, zoneSpec)
  }
}
```

---

zone.js 由入门到放弃之一——通过一场游戏认识 zone.js

A Zone is an execution context that persists across async tasks. You can think of it as thread-local storage for JavaScript VMs.

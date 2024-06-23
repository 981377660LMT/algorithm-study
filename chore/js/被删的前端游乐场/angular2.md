# angular2 和 Angular 源码

## ng2 中的依赖注入

1. 当 Angular 创建组件时，会首先为组件所需的服务找一个注入器（Injector）
2. 注入器是一个`维护服务实例的容器`，存放着以前创建的实例
3. 如果容器中还没有所请求的服务实例，注入器就会`创建`一个服务实例，并且添加到容器中，然后把这个服务返回给 Angular
4. 当所有的服务都被解析完并返回时，Angular 会`以这些服务为参数去调用组件的构造函数`

## Angular2 中的基础概念

1. Angular 应用中的构造块

   - 1.1. 模块

   - 1.2. 组件

   - 1.3. 模板

   - 1.4. 元数据

   - 1.5. 数据绑定

   - 1.6. 指令

   - 1.7. 服务

   - 1.8. 依赖注入

# 谈谈 Angular2 的依赖注入

1. 依赖注入

   - 1.1. 简要说明
   - 1.2. 理解依赖注入

     而所有的设计和架构，都是为了使得我们的工作更加高效。
     项目提供了这样一个注入机制，有人负责提供服务，有人负责消耗服务，而这样的机制提供了中间的接口，并替使用者进行了创建并初始化这样的处理。
     我们只需要知道，`拿到的是完整可用的服务就好了，至于这个服务内部的实现，甚至是它又依赖了怎样的其他服务，都不需要关注`。

2. Angular 与依赖注入
   - 2.1. Angular 依赖注入
     Angular2 官网上有句话：Angular 的依赖注入系统能够即时地创建和交付所依赖的服务。
   - 2.2. NgModule 模块类
   - 2.3. provider 服务提供
   - 2.4. 多级依赖注入
   - 2.5. 对比其他框架谈谈依赖注入

# Angular 框架解读--预热篇

1. 前端框架

   - 1.1. 三大前端“框架”
     React/Vue 专注于构建用户界面，在一定程度上来说为一个 Javascript 库；而 Angular 则提供了前端项目开发中较完整的解决方案
   - 1.2. 低热度的 Angular

   - 1.3. 我对 Angular 的理解
     `大型项目如何设计和管理`这块领域对前端来说依然比较陌生。我们可以借助常见的后台系统架构设计来进行参考和反思，比如`微服务、领域驱动设计、职责驱动设计`等。但这些终究是设计思想，如何才能很好地落地，对前端开发都是不小的考验。
     虽然新人的加入、每个人都按照自己的想法去开发，最终总会变得难以维护，历史债务十分严重。而 Angular 则是唯一一个能限制开发的自由发挥的，可以`让经验不足和经验丰富的开发都写出一样易维护的代码。`

2. Angular 框架解读
   阅读源码，并不是为了熟悉掌握源码本身，更是为了`掌握其中的一些值得借鉴的思考方式和设计`
   依赖注入整体框架的设计
   组件设计与管理
   Provider 与 Service 的设计
   NgModule 模块化组织（多级/分层）的设计
   模板引擎/模板编译过程的整体设计
   Zone 设计：提升计算速度
   JIT/AOT 设计
   元数据设计：（Reflect.metadata）的引入和使用思考
   响应式编程：Rxjs 的引入和使用思考

# Angular 框架解读--元数据和装饰器

1. 装饰器与元数据

- 1.1. 元数据（Metadata）

  在通用的概念中，`元数据是描述用户数据的数据`
  它总结了有关数据的基本信息，可以使查找和使用特定数据实例更加容易。
  例如，`作者，创建日期，修改日期和文件大小`是非常基本的`文档元数据`的示例。

  在用于类的场景下，元数据用于装饰类，来描述类的定义和行为，以便可以配置类的预期行为。

- 1.2. 装饰器（Decorator）
  装饰器可用于对值进行元编程和向其添加功能，而无需从根本上改变其外部行为

2. Angular 中的装饰器和元数据

- 2.1. 使用装饰器和元数据来改变类的行为
  可以用下列装饰器来声明 Angular 的类：
  @Component()、@Directive()、@Pipe()、@Injectable()、@NgModule()
- 2.2. 装饰器的创建过程
- 2.3. 根据装饰器元数据编译组件
- 2.4. 编译过程中的元数据

# Angular 框架解读--视图抽象定义

Angular 版本可在不同的平台上运行：在浏览器中、在移动平台上或在 Web Worker 中。
因此，需要特定级别的抽象来介于平台特定的 API 和框架接口之间。

# Angular 框架解读--Zone 区域之 zone.js

zone 具有当前区域的概念：当前区域是随所有异步操作一起传播的异步上下文，它表示与当前正在执行的堆栈帧/异步任务关联的区域。

```js
zone.fork(zoneSpec): 创建一个新的子区域，并将其parent设置为用于分支的区域
zone.run(callback, ...)：在给定区域中同步调用一个函数
zone.runGuarded(callback, ...)：与run捕获运行时错误相同，并提供了一种拦截它们的机制。如果任何父区域未处理错误，则将其重新抛出。
zone.wrap(callback)：产生一个新的函数，该函数将区域绑定在一个闭包中，并在执行zone.runGuarded(callback)时执行，与 JavaScript 中的Function.prototype.bind工作原理类似。
```

我们可以看到 Zone 的主要实现逻辑（new Zone()/fork()/run()）也相对简单：

```ts
class Zone implements AmbientZone {
  // 获取根区域
  static get root(): AmbientZone {
    let zone = Zone.current
    // 找到最外层，父区域为自己
    while (zone.parent) {
      zone = zone.parent
    }
    return zone
  }
  // 获取当前区域
  static get current(): AmbientZone {
    return _currentZoneFrame.zone
  }
  private _parent: Zone | null // 父区域
  private _name: string // 区域名字
  private _properties: { [key: string]: any }
  // 拦截区域操作时的委托，用于生命周期钩子相关处理
  private _zoneDelegate: ZoneDelegate

  constructor(parent: Zone | null, zoneSpec: ZoneSpec | null) {
    // 创建区域时，设置区域的属性
    this._parent = parent
    this._name = zoneSpec ? zoneSpec.name || 'unnamed' : '<root>'
    this._properties = (zoneSpec && zoneSpec.properties) || {}
    this._zoneDelegate = new ZoneDelegate(
      this,
      this._parent && this._parent._zoneDelegate,
      zoneSpec
    )
  }
  // fork 会产生子区域
  public fork(zoneSpec: ZoneSpec): AmbientZone {
    if (!zoneSpec) throw new Error('ZoneSpec required!')
    // 以当前区域为父区域，调用 new Zone() 产生子区域
    return this._zoneDelegate.fork(this, zoneSpec)
  }
  // 在区域中同步运行某段代码
  public run(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): any
  public run<T>(
    callback: (...args: any[]) => T,
    applyThis?: any,
    applyArgs?: any[],
    source?: string
  ): T {
    // 准备执行，入栈处理
    _currentZoneFrame = { parent: _currentZoneFrame, zone: this }
    try {
      // 使用 callback.apply(applyThis, applyArgs) 实现
      return this._zoneDelegate.invoke(this, callback, applyThis, applyArgs, source)
    } finally {
      // 执行完毕，出栈处理
      _currentZoneFrame = _currentZoneFrame.parent!
    }
  }
}
```

除了上面介绍的，Zone 还提供了许多方法来运行、计划和取消任务，包括：

```ts
interface Zone {
  ...
  // 通过在任务区域中恢复 Zone.currentTask 来执行任务
  runTask<T>(task: Task, applyThis?: any, applyArgs?: any): T;
  // 安排一个 MicroTask
  scheduleMicroTask(source: string, callback: Function, data?: TaskData, customSchedule?: (task: Task) => void): MicroTask;
  // 安排一个 MacroTask
  scheduleMacroTask(source: string, callback: Function, data?: TaskData, customSchedule?: (task: Task) => void, customCancel?: (task: Task) => void): MacroTask;
  // 安排一个 EventTask
  scheduleEventTask(source: string, callback: Function, data?: TaskData, customSchedule?: (task: Task) => void, customCancel?: (task: Task) => void): EventTask;
  // 安排现有任务（对重新安排已取消的任务很有用）
  scheduleTask<T extends Task>(task: T): T;
  // 允许区域拦截计划任务的取消，使用 ZoneSpec.onCancelTask​​ 配置拦截
  cancelTask(task: Task): any;
}
```

zone.js 是如何识别出异步任务的呢？
其实 zone.js 主要是`通过猴子补丁拦截异步 API`（包括 DOM 事件、XMLHttpRequest 和 NodeJS 的 API 如 EventEmitter、fs 等）来实现这些功能：

任务执行的生命周期:

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

zone.js 提供了丰富的生命周期钩子，可以使用 zone.js 的区域能力以及生命周期钩子解决前面我们提到的这些问题：

- 异步代码执行时，`上下文发生了变更`，导致预期不一致：使用 Zone 来执行相关代码
- throw Error 时，无法准确`定位到上下文：使用生命周期钩子 onHandleError 进行处理和跟踪`
- `测试某个函数的执行耗时`，但因为函数内有`异步`逻辑，无法得到准确的执行时间：使用生命周期钩子配合可得到具体的耗时

# Angular 框架解读--Zone 区域之 ngZone

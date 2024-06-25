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
    this._zoneDelegate = new ZoneDelegate(this, this._parent && this._parent._zoneDelegate, zoneSpec)
  }
  // fork 会产生子区域
  public fork(zoneSpec: ZoneSpec): AmbientZone {
    if (!zoneSpec) throw new Error('ZoneSpec required!')
    // 以当前区域为父区域，调用 new Zone() 产生子区域
    return this._zoneDelegate.fork(this, zoneSpec)
  }
  // 在区域中同步运行某段代码
  public run(callback: Function, applyThis?: any, applyArgs?: any[], source?: string): any
  public run<T>(callback: (...args: any[]) => T, applyThis?: any, applyArgs?: any[], source?: string): T {
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
  onFork?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, zoneSpec: ZoneSpec) => Zone
  // 允许拦截回调的 wrap
  onIntercept?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, delegate: Function, source: string) => Function
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
  onHandleError?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, error: any) => boolean
  // 允许拦截任务计划
  onScheduleTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task) => Task
  // 允许拦截任务回调调用
  onInvokeTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task, applyThis: any, applyArgs?: any[]) => any
  // 允许拦截任务取消
  onCancelTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, task: Task) => any
  // 通知对任务队列为空状态的更改
  onHasTask?: (parentZoneDelegate: ZoneDelegate, currentZone: Zone, targetZone: Zone, hasTaskState: HasTaskState) => void
}
```

zone.js 提供了丰富的生命周期钩子，可以使用 zone.js 的区域能力以及生命周期钩子解决前面我们提到的这些问题：

- 异步代码执行时，`上下文发生了变更`，导致预期不一致：使用 Zone 来执行相关代码
- throw Error 时，无法准确`定位到上下文：使用生命周期钩子 onHandleError 进行处理和跟踪`
- `测试某个函数的执行耗时`，但因为函数内有`异步`逻辑，无法得到准确的执行时间：使用生命周期钩子配合可得到具体的耗时

# Angular 框架解读--Zone 区域之 ngZone

NgZone 基于 zone.js 之上再做了一层封装，通过 fork 创建出子区域作为 Angular 区域，
使得在 Angular Zone 内函数中的`所有异步操作可以在正确的时间自动触发变更检测`。

默认情况下，所有异步操作都在 Angular Zone 内，这会自动触发变更检测(handleChange)。
另一个常见的情况是我们不想触发变更检测（`比如不希望像 scroll 等事件过于频繁地进行变更检测，从而导致性能问题`），此时可以使用 NgZone 的 runOutsideAngular()方法，来自己实现变更检测的逻辑。

# Angular 框架解读--模块化组织

说到模块化，前端开发首先会想到 ES6 的模块，这两者其实并没有什么关联：

ES6 模块以文件为单位；Angular 模块则是以 NgModule 为单位。
ES6 模块用于跨文件的功能调用；Angular 模块用于组织有特定意义的功能块。
ES6 模块在编译阶段确认各个模块的依赖关系，模块间关系扁平；Angular 模块则可以带有深度的层次结构。

NgModule 把组件、指令和管道打包成内聚的功能块，每个模块聚焦于一个特性区域、业务领域、工作流或通用工具。运行时，模块相关的信息存储在 NgModuleDef 中：

```ts
// NgModuleDef 是运行时用于组装组件、指令、管道和注入器的内部数据结构
export interface NgModuleDef<T> {
  // 表示模块的令牌，由DI使用
  type: T
  // 要引导的组件列表
  bootstrap: Type<any>[] | (() => Type<any>[])
  // 此模块声明的组件、指令和管道的列表
  declarations: Type<any>[] | (() => Type<any>[])
  // 此模块导入的模块列表或 ModuleWithProviders
  imports: Type<any>[] | (() => Type<any>[])
  // 该模块导出的模块、ModuleWithProviders、组件、指令或管道的列表
  exports: Type<any>[] | (() => Type<any>[])
  // 为该模块计算的 transitiveCompileScopes 的缓存值
  transitiveCompileScopes: NgModuleTransitiveScopes | null
  // 声明 NgModule 中允许的元素的一组模式
  schemas: SchemaMetadata[] | null
  // 应为其注册模块的唯一ID
  id: string | null
}
```

宏观来讲，NgModule 是组织 Angular 应用的一种方式，它们通过@NgModule 装饰器中的元数据来实现这一点，这些元数据可以分成三类：

- 静态的：编译器配置，通过 `declarations` 数组来配置。用于告诉编译器指令的选择器，并通过选择器匹配的方式，决定要把该指令应用到模板中的什么位置
- 运行时：通过 `providers` 数组提供给注入器的配置
- 组合/分组：通过 `imports 和 exports` 数组来把多个 **NgModule** 放在一起，并让它们可用

## 模块化组织

每个 Angular 应用有至少一个模块，该模块称为根模块（AppModule）。Angular 应用的启动，便是由根模块开始的，可以参考后续的依赖注入的引导过程内容。

# Angular 框架解读--依赖注入的基本概念

1. 依赖注入

- 1.1. 依赖倒置原则、控制反转、依赖注入

  - 依赖倒置原则（DIP）：`模块间不应该直接依赖对方，应该依赖一个抽象的规则（接口或者时抽象类）`
  - 控制反转（IoC）: `模块间的依赖关系从程序内部提到外部来实例化管理`。即对象在被创建的时候，由一个调控系统内所有对象的外界实体控制，并将其所依赖的对象的引用传递(注入)给它。
    实现控制反转主要有两种方式：
    依赖注入：被动的接收依赖对象
    依赖查找：主动索取依赖的对象
  - 依赖注入（DI）：`是控制反转的最为常见的一种技术`
    依赖倒置和控制反转两者相辅相成，常常可以一起使用，可有效地降低模块间的耦合。

2.  Angular 中的依赖注入
    DI 框架会在实例化某个类时，向其提供这个类所声明的依赖项（依赖项：指当类需要执行其功能时，所需要的服务或对象）。

- 2.1. Injector 注入器
  Injector 注入器用于创建依赖，会`维护一个容器来管理这些依赖`，并尽可能地复用它们。注入器会提供依赖的一个单例，并把这个单例对象注入到多个组件中。
  我们可以将需要共享的依赖实例添加到注入器中，并通过 Token 查询和检索注入器来获取相应的依赖实例。

  ```ts
  export abstract class Injector {
    // 找不到依赖
    static THROW_IF_NOT_FOUND = THROW_IF_NOT_FOUND
    // NullInjector 是树的顶部
    // 如果你在树中向上走了很远，以至于要在 NullInjector 中寻找服务，那么将收到错误消息，或者对于 @Optional()，返回 null
    static NULL: Injector = new NullInjector()

    // 通过 Token 查询和检索注入器来获取相应的依赖实例
    // 查找依赖的过程也是向上遍历注入器树的过程
    abstract get<T>(token: Type<T> | AbstractType<T> | InjectionToken<T>, notFoundValue?: T, flags?: InjectFlags): T

    // 创建一个新的 Injector 实例，该实例提供一个或多个依赖项
    // 创建一个新的Injector实例时，传入的参数包括Provider：Injector不会直接创建依赖，而是通过Provider来完成的
    // 如果指定的注入器无法解析某个依赖，它就会请求父注入器来解析它
    static create(options: { providers: StaticProvider[]; parent?: Injector; name?: string }): Injector

    // ɵɵdefineInjectable 用于构造一个 InjectableDef
    // 它定义 DI 系统将如何构造 Token，并且在哪些 Injector 中可用
    static ɵprov = ɵɵdefineInjectable({
      token: Injector,
      providedIn: 'any' as any,
      // ɵɵinject 生成的指令：从当前活动的 Injector 注入 Token
      factory: () => ɵɵinject(INJECTOR)
    })

    static __NG_ELEMENT_ID__ = InjectorMarkers.Injector
  }
  ```

- 2.2. Provider 提供者

Provider 提供者用来`告诉注入器应该如何获取或创建依赖`，要想让注入器能够创建服务（或提供其它类型的依赖），必须使用某个提供者配置好注入器。

- 2.3. Angular 中的依赖注入服务
  在 Angular 中，`服务就是一个带有@Injectable 装饰器的类`，它封装了可以在应用程序中复用的非 UI 逻辑和代码。Angular 把组件和服务分开，是为了增进模块化程度和可复用性。

3.  总结
    对于注入器、提供者和可注入服务，我们可以简单地这样理解：

    - 注入器用于创建依赖，会维护一个容器来管理这些依赖，并尽可能地复用它们。
    - 一个注入器中的依赖服务，只有一个实例。
    - 注入器需要使用提供者来管理依赖，并通过 token（DI 令牌）来进行关联。
    - 提供者用于告诉注入器应该如何获取或创建依赖。
    - 可注入服务类会根据元数据编译后，得到可注入对象，该对象可用于创建实例。

    注入器：老板，管人(依赖)的
    提供者：老板的助手们
    可注入服务：干活的

# Angular 框架解读--多级依赖注入设计

在 Angular 应用中，各个组件和模块间又是怎样共享依赖的，同样的服务是否可以多次实例化呢

- 1. 多级依赖注入
  - 1.1. 模块注入器
    - 1.1.1. 平台模块（PlatformModule）注入器
    - 1.1.2. 应用程序根模块（AppModule）注入器
    - 1.1.3. 模块注入器层级
  - 1.2. 元素注入器
    - 1.2.1. 元素注入器的引入
    - 1.2.2. 元素注入器（Element Injector）
    - 1.2.3. 元素注入器与模块注入器的设计
  - 1.3. Angular 解析依赖过程
    - 1.3.1. 合并注入器（Merge Injector）
    - 1.3.2. 解析过程
      在 Angular 应用中，各个组件和模块间又是怎样共享依赖的，同样的服务是否可以多次实例化呢
  - 1.4. 总结

# Angular 框架解读--Ivy 编译器整体设计

https://www.bilibili.com/video/BV12v4y1w7P5/?vd_source=e825037ab0c37711b6120bbbdabda89e

- 通过`增量编译`缩短构建时间
- 减少构建大小

1. ngtsc：TypeScript 编译器，将 Angular 装饰器化为静态属性
2. Ngcc：处理来自 npm 代码并生成等效的 Ivy 版本，就像使用 ngtsc 编译代码一样

标记模版 -> 将标记内容解析为 HTML AST -> 将 HTML AST 转换为 Angular 模版 AST -> 将模版 AST 转换为模版函数
parse、transform、generate

1. Ivy 编译器能力
   Angular 重构编译器并将之命名为 Ivy 编译器，这对于 Angular 框架来说有着非常重要的意义，有点类似于 React 重构 Fiber。
   1.1. Ivy 新特性
2. Ivy 架构设计
   2.1. 模板编译
   2.2. Typescript 解析器
   2.3. 编译器设计
   2.4. Ivy 编译模型
3. 总结
   3.1. 参考

# Angular 框架解读--Ivy 编译器的视图数据和依赖解析

# Angular 框架解读--Ivy 编译器之 CLI 编译器

# Angular 框架解读--Ivy 编译器之心智模型

**装饰器就是编译器**

# Angular 框架解读--Ivy 编译器之 AOT/JIT

在 Angular 中，提供了两种方式编译 Angular 应用：

即时编译 (JIT，Just in time)：它会在`运行期间在浏览器中编译你的应用`
预先编译（AOT，Ahead of Time）：它会在`构建时编译你的应用和库`

# Angular 框架解读--Ivy 编译器之增量 DOM

树状的结构会带来算法的简化以及性能的提升。
增量 DOM 的设计核心思想是：

在创建新的（虚拟）DOM 树时，沿着现有的树走，并在进行时找出更改。
如果没有变化，则不分配内存；
如果有，改变现有树（仅在绝对必要时分配内存）并将差异应用到物理 DOM。
https://github.com/google/incremental-dom
它是一个用于`表达和应用DOM树更新`的库

# Angular 框架解读--Ivy 编译器之变更检测

增量编译

在 Angular 中将被标记为 `CheckAlways 或者 Dirty` 的组件进行视图刷新，在每个变更周期中，会执行` template()模板函数中的更新模式下逻辑`。而在 template()模板函数中的具体指令逻辑中，还会根据原来的值和新的值进行比较，`有差异的时候才会进行更新`。

MobX 是一个通过**函数响应式编程 (TFRP)** 实现状态管理的库。它的核心哲学是：**任何源自应用状态的东西都应该自动地获得**。

为了深入讲解 MobX 的 API，我将它们分为四个核心类别：**可观察状态 (Observables)**、**动作 (Actions)**、**衍生 (Derivations)** 和 **工具/配置 (Utils/Config)**。

---

### 一、 核心概念：可观察状态 (Observables)

这是 MobX 的基石，用于定义“随时间变化的数据”。

#### 1. `makeObservable(target, annotations?, options?)`

这是 MobX 6+ 推荐的 API。它将现有的对象属性转换为可观察对象。通常在类的 `constructor` 中使用。

- **用法**：

  ```typescript
  class Store {
    count = 0
    get double() {
      return this.count * 2
    }

    constructor() {
      makeObservable(this, {
        count: observable, // 标记为状态
        double: computed, // 标记为计算属性
        increment: action // 标记为动作
      })
    }
    increment() {
      this.count++
    }
  }
  ```

- **特点**：显式、清晰，对构建工具（如 Webpack/Rollup）最友好，Tree-shaking 支持最好。

#### 2. `makeAutoObservable(target, overrides?, options?)`

这是 `makeObservable` 的“懒人版”。它会自动推断属性类型：

- 所有属性 -> `observable`
- 所有 getter -> `computed`
- 所有方法 -> `action`
- 生成器函数 -> `flow`

- **用法**：
  ```typescript
  constructor() {
    makeAutoObservable(this);
    // 排除某个属性：
    // makeAutoObservable(this, { nonObservableField: false });
  }
  ```
- **缺点**：不能用于继承类（super class），因为自动推断会覆盖父类行为。

#### 3. `observable(value)` / `observable.box(value)`

用于创建独立的可观察对象，不依赖于类。

- **对象/数组/Map/Set**：

  ```typescript
  const todos = observable([{ title: 'Learn MobX', done: false }])
  ```

  MobX 会自动将内部属性深度代理（Deep Proxy）。

- **原始值 (Boxed Observables)**：
  JS 的原始值（string, number）是不可变的，无法直接代理。需要用 `.box` 包装。
  ```typescript
  const name = observable.box('MobX')
  console.log(name.get()) // 读取
  name.set('React') // 修改
  ```

#### 4. `observable.ref` / `observable.shallow` / `observable.struct`

这些是**修饰符 (Modifiers)**，用于控制观察的深度和行为。

- **`observable.ref`**：**不**进行深度观察。只有当引用发生变化时才触发更新。
  - _场景_：存储巨大的不可变数据（如从 API 获取的大 JSON），或者 DOM 节点引用。
- **`observable.shallow`**：只观察第一层。
  - _场景_：数组本身的变化（push/pop）需要响应，但数组里面的对象不需要被转换成 Observable。
- **`observable.struct`**：进行结构比较（Deep Equal）。
  - _场景_：如果新值和旧值在结构上相等（`{x:1}` vs `{x:1}`），则不触发更新。默认是引用比较。

---

### 二、 改变状态：动作 (Actions)

MobX 强烈建议（在严格模式下强制）所有修改状态的代码都必须在 Action 中执行。

#### 1. `action(fn)` / `action(name, fn)`

将函数标记为动作。

- **作用**：
  1.  **批处理 (Batching)**：在 Action 结束前，不会通知观察者。避免多次修改导致多次渲染。
  2.  **调试**：在 DevTools 中显示 Action 名称。

#### 2. `runInAction(fn)`

用于在异步回调或临时代码块中修改状态。

- **场景**：你不想专门为此写一个类方法，只想快速改个值。
  ```typescript
  // 异步回调中
  fetchData().then(data => {
    runInAction(() => {
      this.data = data
      this.loading = false
    })
  })
  ```

#### 3. `flow(generatorFn)`

处理**异步 Action** 的终极方案。

- **痛点**：普通的 `async/await` 在 `await` 之后，代码就不在 `action` 上下文中了，修改状态会报错（在严格模式下）。
- **解决**：使用生成器函数 (`function*`) 代替 `async`。
  ```typescript
  fetchData = flow(function* (this: Store) {
    this.loading = true
    try {
      const res = yield api.getData() // yield 代替 await
      this.data = res // 这里依然在 action 上下文中，无需 runInAction
    } finally {
      this.loading = false
    }
  })
  ```

---

### 三、 响应变化：衍生 (Derivations)

衍生分为两类：**计算值 (Computed values)** 和 **副作用 (Reactions)**。

#### 1. `computed(getter)`

- **定义**：根据现有的 Observable 派生出新值。
- **核心特性**：
  1.  **缓存 (Memoized)**：只要依赖没变，多次访问直接返回缓存值。
  2.  **惰性 (Lazy)**：只有当有人使用它（如 UI 组件）时，它才会计算。
  3.  **纯函数**：不应有副作用。

#### 2. `autorun(effectFn)`

- **定义**：**自动运行**传入的函数。
- **机制**：立即执行一次 -> 追踪函数内访问的所有 Observable -> 依赖变化时重新执行。
- **场景**：打印日志、持久化数据、更新 UI（MobX-React 的核心原理）。
  ```typescript
  autorun(() => {
    console.log('Count is:', store.count)
  })
  ```

#### 3. `reaction(dataFn, effectFn, options?)`

- **定义**：比 `autorun` 更精细的控制。
- **机制**：
  1.  执行 `dataFn`（数据函数），追踪依赖，返回数据。
  2.  当 `dataFn` 返回的数据变化时，执行 `effectFn`（副作用函数）。
- **区别**：`autorun` 既追踪又执行；`reaction` 分离了追踪和执行。且 `reaction` **不会立即执行副作用**（除非配置 `fireImmediately: true`）。
  ```typescript
  reaction(
    () => store.todos.length, // 追踪：只关心长度变化
    (length, prevLength) => {
      // 执行：副作用
      console.log(`Todos changed from ${prevLength} to ${length}`)
    }
  )
  ```

#### 4. `when(predicate, effect?)`

- **定义**：当 `predicate` 返回 true 时，执行 `effect`，然后**自动销毁**。
- **场景**：一次性的等待逻辑。

  ```typescript
  // 等待 isLoading 变为 false
  when(
    () => !store.isLoading,
    () => console.log('Loaded!')
  )

  // 或者使用 Promise 写法
  await when(() => !store.isLoading)
  ```

---

### 四、 工具与配置 (Utils & Config)

#### 1. `toJS(observableData)`

- **作用**：将 Observable 对象递归转换为普通的 JS 对象。
- **场景**：在 console.log 调试时（避免看到 Proxy 对象），或者传参给不支持 Proxy 的第三方库时。

#### 2. `configure(options)`

全局配置 MobX 的行为。

- **常用配置**：
  ```typescript
  configure({
    enforceActions: 'always', // 强制所有状态修改必须在 Action 中（推荐开启）
    computedRequiresReaction: true, // 禁止在 Reaction 之外直接读取 computed（防止昂贵的计算被浪费）
    reactionScheduler: f => f() // 自定义调度器
  })
  ```

#### 3. `intercept(target, property?, handler)`

- **作用**：在变化**发生前**拦截。
- **能力**：可以修改变化的值，或者完全取消这次变化。
- **场景**：数据验证、数据归一化。
  ```typescript
  intercept(store, 'age', change => {
    if (change.newValue < 0) {
      change.newValue = 0 // 修正数据
    }
    return change
  })
  ```

#### 4. `observe(target, property?, handler)`

- **作用**：在变化**发生后**监听。
- **区别**：`autorun/reaction` 追踪的是“值”，`observe` 追踪的是“特定属性的赋值事件”。它是更底层的 API。

---

### 五、 MobX-React 专用 API (如果是 React 项目)

虽然不属于 MobX 核心库，但通常配合使用。

#### 1. `observer(Component)`

- **作用**：高阶组件 (HOC)。将 React 组件转变为响应式组件。
- **原理**：在组件的 `render` 函数外包裹了一层 `reaction`。当组件渲染时访问的 Observable 变化，组件自动重渲染。

#### 2. `Observer` 组件

- **作用**：局部渲染。
- **场景**：不想重渲染整个组件，只想重渲染组件里的一小块区域。
  ```tsx
  <Observer>{() => <div>{store.count}</div>}</Observer>
  ```

#### 3. `useLocalObservable(initializer)`

- **作用**：在函数组件中创建生命周期稳定的 Observable 对象（替代 `useState` + `useEffect` 的复杂逻辑）。

---

### 总结：API 选择指南

1.  **定义状态**：首选 `makeAutoObservable`（简单），次选 `makeObservable`（精细）。
2.  **修改状态**：必须用 `action`。异步操作首选 `flow`，次选 `runInAction`。
3.  **响应状态**：
    - UI 渲染 -> `observer` (React)
    - 派生新值 -> `computed`
    - 打印日志/持久化 -> `autorun`
    - 精细控制副作用 -> `reaction`
    - 一次性等待 -> `when`
4.  **调试**：用 `toJS` 查看原始数据。

---

好的，既然你已经掌握了 MobX 的核心 API 和基本原理，接下来我们将**深入 MobX 的源码架构和高级运行机制**。

结合你提供的 src 目录结构（虽然附件内容为空，但我熟悉 MobX 的源码结构），我们将从**源码视角**来剖析 MobX 是如何运作的。

---

### 一、 MobX 源码的核心模块划分

MobX 的源码结构非常清晰，主要分为以下几个核心层级：

1.  **Core (核心层)**: 实现 Observable, Computed, Reaction 的原子逻辑。
2.  **Types (类型层)**: 定义各种 Observable 数据结构（ObservableObject, ObservableArray, ObservableMap）。
3.  **API (接口层)**: 暴露给用户的 `action`, `autorun`, `makeObservable` 等函数。
4.  **Utils (工具层)**: 内部使用的工具函数。

我们将重点剖析 **Core** 和 **Types**。

---

### 二、 核心机制深度剖析

#### 1. 原子 (Atom) —— 状态的最小单位

在 MobX 源码中，所有的 Observable 数据结构（Object, Array, Map）底层都持有一个或多个 `Atom`。

- **源码位置**: `core/atom.ts`
- **作用**: `Atom` 是 MobX 依赖收集图中的**节点**。它不存储值，只负责维护“谁依赖了我”和“我依赖了谁”的关系。
- **核心属性**:
  - `observers`: Set<Derivation>。当前正在观察这个 Atom 的观察者（如 Computed 或 Reaction）。
  - `reportObserved()`: 当值被读取时调用。如果当前有全局的 `Derivation` 正在运行，就建立依赖关系。
  - `reportChanged()`: 当值被修改时调用。通知所有的 `observers` 状态已过期。

#### 2. 派生 (Derivation) —— 依赖图的消费者

`Derivation` 是所有“观察者”的基类接口，`ComputedValue` 和 `Reaction` 都实现了它。

- **源码位置**: `core/derivation.ts`
- **核心状态机**:
  MobX 使用一个非常巧妙的状态机来优化性能，避免不必要的计算。每个 Derivation 都有一个 `dependenciesState` 属性：

  - **`UP_TO_DATE` (0)**: 数据是最新的，无需计算。
  - **`POSSIBLY_STALE` (1)**: 依赖项中有一个 Computed 变脏了，但我不知道它的新值是否真的变了（可能计算后结果没变）。此时需要先去问问那个 Computed。
  - **`STALE` (2)**: 依赖项中有一个 Observable 变了，我必须重新运行。

  **优化流程**:
  当一个 Observable 变化时，它会将依赖它的 Derivation 标记为 `STALE`。
  当一个 Computed 变化时，它会将依赖它的 Derivation 标记为 `POSSIBLY_STALE`。
  当 Derivation 准备重新运行时，如果发现是 `POSSIBLY_STALE`，它会先触发依赖的 Computed 重新计算。如果 Computed 算出来的值和上次一样（引用相等），Derivation 就会回到 `UP_TO_DATE`，**完全跳过执行**。这就是 MobX 高性能的秘密。

#### 3. 全局状态 (GlobalState)

- **源码位置**: `core/globalstate.ts`
- **作用**: 维护 MobX 的运行时环境。
- **关键变量**:
  - `trackingDerivation`: 当前正在运行的 Derivation（依赖收集的上下文）。
  - `inBatch`: 当前是否处于批处理模式（Transaction）。
  - `pendingReactions`: 等待执行的 Reaction 队列。

---

### 三、 复杂数据结构的实现原理

#### 1. ObservableObject (可观察对象)

当你调用 `makeObservable(this, { count: observable })` 时，MobX 实际上做了什么？

- **源码位置**: `types/observableobject.ts`
- **原理**:
  1.  MobX 会在你的对象上挂载一个隐藏属性 `$mobx`，这是一个 `ObservableObjectAdministration` 实例。
  2.  它会将 `count` 属性替换为 `get/set` 访问器（Proxy 模式下使用 Proxy 拦截）。
  3.  **Get**: 调用 `$mobx.read(key)` -> 找到对应的 `ObservableValue` -> 调用 `reportObserved()`。
  4.  **Set**: 调用 `$mobx.write(key, value)` -> 找到对应的 `ObservableValue` -> 比较新旧值 -> 调用 `reportChanged()`。

#### 2. ObservableArray (可观察数组)

数组比较特殊，因为它的长度可变，且有大量变异方法（push, pop, splice）。

- **源码位置**: `types/observablearray.ts`
- **原理**:
  MobX 并没有直接修改原生数组，而是创建了一个**类数组对象**（Proxy）。
  - **索引访问 (`arr[0]`)**: 被 Proxy 拦截，映射到底层的 ObservableValue。
  - **长度变化**: 数组内部维护了一个特殊的 Atom 专门用来追踪 `length` 的变化。
  - **方法劫持**: `push`, `pop` 等方法都被重写。它们不仅修改数据，还会调用 `reportChanged()`，并且这些操作都在一个 `transaction` 中执行，确保一次 `push` 只触发一次更新。

---

### 四、 事务与批处理 (Transactions & Batching)

这是 MobX 性能优化的另一个关键。

- **问题**:
  ```typescript
  store.a = 1 // 触发更新
  store.b = 2 // 触发更新
  // 如果组件同时依赖 a 和 b，会渲染两次。
  ```
- **Action 的作用**:
  `action` 本质上是一个 `transaction`（事务）。
  ```typescript
  action(() => {
    store.a = 1
    store.b = 2
  })
  ```
  **执行流程**:
  1.  `startBatch()`: 全局 `inBatch` 计数器 +1。此时所有的 `reportChanged` 只会把 Reaction 放入 `pendingReactions` 队列，**不执行**。
  2.  执行函数体：修改 `a`，修改 `b`。
  3.  `endBatch()`: 全局 `inBatch` 计数器 -1。
  4.  如果计数器归零，遍历 `pendingReactions` 队列，执行所有 Reaction。

---

### 五、 常见的高级场景与陷阱

#### 1. 异步 Action 的陷阱

在 MobX 6 之前，异步 Action 是个大坑。

```typescript
// 错误示范
@action
async load() {
  this.loading = true;
  const data = await fetch(); // await 之后，上下文丢失！
  this.data = data; // 报错！因为这行代码不在 action 中执行
  this.loading = false;
}
```

**原因**: `await` 之后的代码是在微任务队列中执行的，此时 `action` 的 `startBatch/endBatch` 早就结束了。
**解决**: 使用 `runInAction` 包裹后续代码，或者使用 `flow`（生成器）。

#### 2. 细粒度更新 (Fine-grained Reactivity)

MobX 的更新粒度非常细。

```typescript
const store = observable({
  user: {
    name: 'Jack',
    address: { city: 'Beijing' }
  }
})

autorun(() => {
  console.log(store.user.address.city)
})
```

- 修改 `store.user.name = "Tom"` -> **不会**触发 autorun。
- 修改 `store.user.address = { city: "Shanghai" }` -> **会**触发。
- 修改 `store.user = ...` -> **会**触发。

**原理**: `autorun` 在运行时，访问了三个 Observable：

1.  `store` 的 `user` 属性。
2.  `user` 对象的 `address` 属性。
3.  `address` 对象的 `city` 属性。
    只有这三个属性对应的 Atom 发生变化，才会触发更新。

#### 3. 内存泄漏与清理

- **Reaction 泄漏**: 如果你手动创建了 `reaction` 或 `autorun`，并且它们依赖了长生命周期的 Observable，那么这个 Reaction 永远不会被 GC，因为它被 Observable 的 `observers` 集合引用着。
- **解决**: `autorun` 返回一个 `dispose` 函数。在组件卸载时必须调用它。
  ```typescript
  const dispose = autorun(...)
  // ...
  dispose(); // 切断依赖关系，允许 GC 回收
  ```

---

### 六、 总结：MobX 的设计哲学

MobX 的源码极其复杂，但设计哲学极其简单：

1.  **同步执行**: 依赖收集是同步的，计算值的更新也是同步的（惰性求值）。
2.  **最小化计算**: 利用状态机 (`UP_TO_DATE` / `STALE`) 避免无用计算。
3.  **透明性**: 通过 Proxy 和 getter/setter，让用户感觉不到“响应式系统”的存在。

如果你想阅读源码，建议的顺序是：

1.  `core/atom.ts` (理解节点)
2.  `core/observable.ts` (理解值)
3.  `core/derivation.ts` (理解状态机)
4.  `core/reaction.ts` (理解调度)

---

好的，既然你要求“一针见血”地讲解 MobX 源码，那我们就跳过那些辅助性的工具代码，直接切入 MobX 的**心脏**。

MobX 的核心其实就是一套**基于观察者模式的依赖追踪系统**。这套系统由三个核心类支撑：`Atom`（被观察者）、`Derivation`（观察者）、`Reaction`（调度器）。

我将结合你提供的 src 目录结构（虽然附件内容为空，但我会基于 MobX v6 的标准源码结构）进行剖析。

---

### 一、 心脏起搏器：`core/atom.ts`

**一针见血：`Atom` 是所有 Observable 的“身份证”和“通讯录”。**

任何一个 `observable`（无论是对象属性、数组元素，还是 `computed`），底层都对应一个 `Atom` 实例。它不存值，只存关系。

#### 核心源码逻辑（简化版）：

```typescript
class Atom {
  // 谁在观察我？（Set<Derivation>）
  observers = new Set()

  // 1. 读操作时调用：建立联系
  reportObserved() {
    // globalState.trackingDerivation 就是当前正在运行的那个 autorun/computed
    const derivation = globalState.trackingDerivation
    if (derivation) {
      // 双向绑定：我记录它，它也记录我
      this.observers.add(derivation)
      derivation.newObserving.push(this)
    }
  }

  // 2. 写操作时调用：通知变更
  reportChanged() {
    // 遍历通讯录，告诉所有观察者：我变了！
    for (const observer of this.observers) {
      observer.onBecomeStale() // 标记为“陈旧”
    }
    // 触发调度，让它们重新运行
    globalState.pendingReactions.push(...this.observers)
    runReactions() // 立即或在 batch 结束后执行
  }
}
```

**关键点**：

- **依赖收集是动态的**：每次 `reportObserved` 都会重新确认关系。如果 `autorun` 里有个 `if (bool) a else b`，当 `bool` 变了，依赖关系会从 `a` 切换到 `b`。这是 MobX 比 Vue 2 更灵活的地方。

---

### 二、 大脑：`core/derivation.ts`

**一针见血：`Derivation` 是状态机，负责决定“要不要重新计算”。**

`ComputedValue` 和 `Reaction` (autorun) 都是 `Derivation`。它们的核心任务是**懒**——能不跑就不跑。

#### 核心状态机（State Machine）：

MobX 用一个 `dependenciesState` 字段来控制计算，这是性能优化的核心。

1.  **`UP_TO_DATE` (0)**: 我是最新的，别烦我。
2.  **`POSSIBLY_STALE` (1)**: 我依赖的某个 `Computed` 变脏了，但我不知道它的新值是不是真的变了（比如 `a=1` -> `a=2` -> `computed=a>0`，结果还是 `true`）。
3.  **`STALE` (2)**: 我依赖的某个 `Observable` 变了，我必须重算。

#### 核心源码逻辑（`shouldCompute`）：

```typescript
function shouldCompute(derivation) {
  if (derivation.dependenciesState === UP_TO_DATE) return false

  if (derivation.dependenciesState === POSSIBLY_STALE) {
    // 关键优化：先去问问那些变脏的 Computed，你们算出来的新值变了吗？
    for (const dep of derivation.observing) {
      if (dep.isComputed) {
        // 强制 Computed 重算
        dep.get()
        // 如果 Computed 算完发现值没变，它会把自己设回 UP_TO_DATE
        if (derivation.dependenciesState === STALE) return true
      }
    }
  }

  // 如果还是 POSSIBLY_STALE，说明所有 Computed 算完都没变，那我也不用变
  return derivation.dependenciesState === STALE
}
```

**一针见血**：这个机制解决了**菱形依赖**和**无效更新**的问题。

---

### 三、 调度中心：`core/reaction.ts`

**一针见血：`Reaction` 是副作用的执行者，负责“防抖”和“批处理”。**

`autorun` 和 `observer` (React) 本质上都是 `Reaction`。

#### 核心源码逻辑（`runReaction`）：

```typescript
class Reaction {
  onBecomeStale() {
    // 当依赖变了，我不会立即跑，而是把自己扔进全局队列
    globalState.pendingReactions.push(this)
    // 尝试启动调度器
    scheduleReactions()
  }

  runReaction() {
    // 1. 开启一个新的追踪上下文
    const prev = globalState.trackingDerivation
    globalState.trackingDerivation = this

    // 2. 清空旧依赖（准备重新收集）
    this.newObserving = []

    // 3. 执行用户的函数（这会触发 Atom.reportObserved，填充 newObserving）
    this.onInvalidate()

    // 4. 绑定新依赖，解绑旧依赖（Diff 算法）
    bindDependencies(this)

    // 5. 恢复上下文
    globalState.trackingDerivation = prev
  }
}
```

**关键点**：

- **两阶段提交**：`onBecomeStale` (标记) -> `runReaction` (执行)。
- **自动清理**：每次执行都会重新收集依赖，之前用到但这次没用到的依赖会被自动解绑（`removeObserver`），防止内存泄漏。

---

### 四、 魔法外衣：`types/observableobject.ts`

**一针见血：`ObservableObject` 是 `Atom` 的代理人 (Proxy)。**

当你写 `makeAutoObservable(this)` 时，MobX 把你的对象变成了一个“傀儡”。

#### 核心源码逻辑（Proxy Handler）：

```typescript
const objectProxyTraps = {
  get(target, name) {
    // 1. 拿到幕后管理者 Administration
    const adm = target.$mobx
    // 2. 拿到属性对应的 ObservableValue (里面藏着 Atom)
    const observable = adm.values.get(name)
    // 3. 触发依赖收集
    observable.reportObserved()
    // 4. 返回值
    return observable.value
  },

  set(target, name, value) {
    const adm = target.$mobx
    const observable = adm.values.get(name)
    // 1. 拦截修改
    observable.setNewValue(value) // 这里面会调用 reportChanged
    return true
  }
}
```

**一针见血**：用户以为在操作对象，其实是在操作 `Atom`。

---

### 五、 总结：MobX 的运行闭环

把上面四个部分串起来，就是 MobX 的一次完整生命周期：

1.  **定义**：`makeObservable` 创建 `Atom` 和 `Proxy`。
2.  **挂载**：`autorun` 创建 `Reaction`，设置 `globalState.trackingDerivation` 为自己，开始执行函数。
3.  **读取**：函数读取 `proxy.prop` -> 触发 `Atom.reportObserved` -> `Atom` 把 `Reaction` 加入自己的 `observers`。
4.  **修改**：用户修改 `proxy.prop = 2` -> 触发 `Atom.reportChanged` -> 遍历 `observers` 找到 `Reaction`。
5.  **调度**：`Reaction` 被推入 `pendingReactions` 队列。
6.  **执行**：当前同步代码执行完（或 `action` 结束），MobX 遍历队列，执行 `Reaction.runReaction`。
7.  **重置**：`Reaction` 重新执行，清空旧依赖，重复步骤 3。

这就是 MobX。没有黑魔法，只有严谨的数据结构和状态机。

---

好的，既然你提供了 src 的目录结构（虽然附件内容为空，但我会基于 MobX v6 的标准源码结构），我们将继续深入挖掘 MobX 源码中那些**容易被忽视但至关重要**的细节。

我们将重点补充讲解以下三个方面：

1.  **`ObservableValue` 的值变更检测与拦截机制**（如何防止无限循环和无效更新）。
2.  **`ComputedValue` 的双向链表优化**（如何高效管理依赖）。
3.  **`Transaction` (事务) 的底层实现**（`action` 是如何工作的）。

---

### 一、 `ObservableValue` 的精细控制

在 `types/observablevalue.ts` 中，`ObservableValue` 是最基础的可观察对象封装。它不仅仅是一个存值的容器，还包含了一套复杂的**变更检测逻辑**。

#### 1. 值的预处理 (Enhancer)

当你定义 `observable.ref` 或 `observable.shallow` 时，MobX 是如何区分的？
答案在于 **Enhancer (增强器)**。

- **源码逻辑**：

  ```typescript
  class ObservableValue<T> extends Atom {
    constructor(
      value: T,
      public enhancer: IEnhancer<T>, // <--- 关键点
      public name = 'ObservableValue'
    ) {
      super(name)
      this.value = enhancer(value, undefined, name)
    }

    setNewValue(newValue: T) {
      const oldValue = this.value
      // 1. 使用 enhancer 处理新值
      newValue = this.enhancer(newValue, oldValue, this.name)

      // 2. 变更检测 (checkIfStateModificationsAreAllowed)
      if (newValue !== oldValue) {
        this.value = newValue
        this.reportChanged()
      }
    }
  }
  ```

- **Enhancer 的种类**：
  - `deepEnhancer` (默认): 递归地将新值也变成 Observable。
  - `referenceEnhancer` (`observable.ref`): 直接返回原值，不递归。
  - `shallowEnhancer` (`observable.shallow`): 只代理第一层。

#### 2. 拦截器 (Interceptor) 与 监听器 (Listener)

MobX 允许你在值改变**前**拦截 (`intercept`) 和改变**后**监听 (`observe`)。

- **源码逻辑**：
  `ObservableValue` 内部维护了 `interceptors` 和 `changeListeners` 两个数组。
  - 在 `setNewValue` 的第一步，会遍历 `interceptors`，允许它们修改 `newValue` 或抛出异常阻止更新。
  - 在 `reportChanged` 之后，会遍历 `changeListeners`，通知外部。

---

### 二、 `ComputedValue` 的双向链表优化

`ComputedValue` (`core/computedvalue.ts`) 是 MobX 中最复杂的类之一。为了在依赖树极其庞大时依然保持高性能，MobX 并没有使用简单的数组来存储依赖，而是使用了一种**侵入式的双向链表**结构。

#### 1. 为什么不用 Set 或 Array？

如果一个 Computed 依赖了 1000 个 Observable，每次重新计算时，都需要清空旧依赖、添加新依赖。如果用 `Set`，GC 压力很大。如果用 `Array`，查找和去重很慢。

#### 2. 源码中的优化

MobX 在 `Derivation` (Computed/Reaction) 和 `Observable` (Atom) 之间维护了一个**对象池**和**双向链接**。

- **核心结构**：
  每个 `Derivation` 都有一个 `observing` 数组，存储它依赖的 `Atom`。
  每个 `Atom` 都有一个 `observers` Set，存储依赖它的 `Derivation`。

  **Diff 算法优化**：
  当 Computed 重新运行时，它会生成一个新的 `newObserving` 数组。
  MobX 会对比 `observing` (旧) 和 `newObserving` (新)：

  1.  新旧都有：保持不变。
  2.  旧有新无：从 Atom 的 `observers` 中移除自己（解绑）。
  3.  旧无新有：将自己添加到 Atom 的 `observers` 中（绑定）。

  这个过程经过了极度优化，使用了 `diffValue` 标记位来避免 O(N\*M) 的复杂度，做到了 O(N+M)。

---

### 三、 `Transaction` (事务) 与 `Action`

`action` (`core/action.ts`) 的本质是 `startBatch` 和 `endBatch`。

#### 1. 嵌套事务的处理

MobX 支持 `action` 嵌套调用。

```typescript
action(() => {
  store.a = 1
  action(() => {
    store.b = 2
  })()
})
```

- **源码逻辑**：
  MobX 使用一个全局计数器 `globalState.inBatch`。

  - `startBatch()`: `inBatch++`。
  - `endBatch()`: `inBatch--`。
  - 只有当 `inBatch` 降为 0 时，才会触发 `runReactions()`。

  这意味着内部的 `action` 结束时，不会触发 UI 更新，只有最外层的 `action` 结束时才会触发。

#### 2. `untracked` 与 `allowStateChanges`

有时候我们需要在 `action` 或 `computed` 中打破规则。

- **`untracked(fn)`**:
  暂时将 `globalState.trackingDerivation` 设为 `null`。
  这样 `fn` 中访问的 Observable 就不会被收集为依赖。
  _场景_：在 computed 中打印日志，但不希望日志里的 observable 导致 computed 重新计算。

- **`allowStateChanges(allow, fn)`**:
  控制是否允许修改 Observable。
  - `Computed` 执行期间，默认 `allowStateChanges(false)`，防止副作用。
  - `Action` 执行期间，默认 `allowStateChanges(true)`。

---

### 四、 总结：MobX 源码的“一针见血”

如果把 MobX 源码比作一个精密的机械表：

1.  **Atom (发条)**: 动力的源头，负责通知。
2.  **Derivation (齿轮)**: 咬合在一起，传递动力。通过 `dependenciesState` 状态机（UP_TO_DATE / STALE）来决定齿轮是否转动。
3.  **Action (外壳)**: 保护内部结构，确保一次上发条（Batch）只带动一次齿轮转动，防止齿轮乱转（UI 闪烁）。
4.  **Proxy (表盘)**: 用户看到的界面，把底层的机械运动转化为直观的读写操作。

理解了这四点，你就彻底看透了 MobX。

===

好的，既然我们要彻底吃透 MobX 源码，前几轮我们已经拆解了核心的“心脏”（Atom）、“大脑”（Derivation）和“骨架”（ObservableValue/Computed）。

这一轮，我们将深入 MobX 中**最实用但实现最棘手**的三个部分：**数组的黑魔法 (`ObservableArray`)**、**异步的救星 (`flow`)** 以及 **调试的后门 (`Spy` 机制)**。

---

### 一、 数组的黑魔法：`types/observablearray.ts`

在 JavaScript 中，拦截数组的变化比拦截对象难得多。对象只有 `get/set`，而数组有索引访问、`length` 变化以及一堆变异方法 (`push`, `pop`, `splice` 等)。

#### 1. 核心设计：`StubArray` 与 `Administration`

MobX 的 `ObservableArray` 并不是直接继承自原生 `Array`，而是采用了一种**代理 + 管理器**的模式。

- **Legacy Mode (Proxy 之前)**: MobX 实际上创建了一个“伪数组”（StubArray），手动把 `0..1000` 的索引定义为 `get/set` 属性。这就是为什么旧版 MobX 数组不支持越界赋值（`arr[100] = 1` 当 `length=0` 时无效）的原因。
- **Modern Mode (Proxy)**: 使用 `Proxy` 拦截 `get/set`。

#### 2. 万法归宗：`spliceWithArray`

你可能会以为 `push`、`pop`、`shift` 各自都有复杂的实现。
**一针见血**：**所有**变异方法，最终都会调用同一个底层方法 —— `spliceWithArray`。

- **源码逻辑**：

  ```typescript
  // 伪代码简化
  push(...items) {
    const currentLength = this.length;
    // push 本质上就是在数组末尾 splice 插入元素
    this.spliceWithArray(currentLength, 0, items);
    return this.length;
  }
  ```

- **`spliceWithArray` 做了什么？**

  1.  **开启事务**：`startBatch()`。
  2.  **处理新值**：使用 `enhancer`（如 `deepEnhancer`）把插入的新元素变成 Observable。
  3.  **更新原生数组**：调用底层原生数组的 `splice`。
  4.  **更新 Length Atom**：如果数组长度变了，通知监听 `length` 的观察者。
  5.  **通知变更**：调用 `reportChanged()`，通知监听数组内容的观察者。
  6.  **结束事务**：`endBatch()`。

  这就是为什么你在 `action` 外面写 `arr.push(1)` 也是安全的，因为 `push` 内部自带了一个微型事务。

---

### 二、 异步的救星：`api/flow.ts`

在 MobX 6 之前，异步 Action 是新手的噩梦。`async/await` 会导致 `await` 之后的代码脱离 `action` 上下文。

**一针见血**：`flow` 本质上是一个**自动执行生成器 (Generator) 并强制在 Action 中恢复执行**的包装器。

#### 源码深度剖析：

`flow` 接收一个 Generator 函数，返回一个 Promise。它的魔法在于如何处理 `yield` 返回的 Promise。

```typescript
export function flow(generator) {
  return function (...args) {
    const iterator = generator.apply(this, args)

    return new Promise((resolve, reject) => {
      // 定义一个步进函数
      function step(nextValue) {
        let result
        try {
          // 关键点 1: 恢复 Generator 执行
          result = iterator.next(nextValue)
        } catch (e) {
          return reject(e)
        }

        if (result.done) {
          return resolve(result.value)
        }

        // 关键点 2: 处理 yield 出去的 Promise
        Promise.resolve(result.value).then(
          // 关键点 3: 在 Action 中执行下一步！
          // createAction 包装了 step，确保 step 执行时 globalState.inBatch > 0
          createAction('flow_step', step),
          createAction('flow_error', e => iterator.throw(e))
        )
      }

      // 启动
      step(undefined)
    })
  }
}
```

**核心差异**：

- 普通的 `async/await`: 浏览器引擎控制 Promise 的 `then` 回调，直接放入微任务队列执行，MobX 无法插手。
- `flow`: MobX 手动控制 `then` 回调，并用 `action` 包裹它。这样，当 Promise 回来时，代码依然运行在 `action` 的保护伞下。

---

### 三、 调试的后门：`core/spy.ts`

MobX DevTools 是怎么知道哪个 Action 触发了哪个 Reaction，耗时多少的？靠的是 **Spy (间谍)** 系统。

#### 1. 全局事件总线

MobX 内部有一个全局的 `globalState.spyListeners` 数组。

#### 2. 埋点 (Instrumentation)

MobX 在源码的关键位置埋下了大量的 `spyReport` 调用。

- **Action 开始时**：`spyReportStart({ type: 'action', name: ... })`
- **值修改时**：`spyReport({ type: 'update', newValue, oldValue ... })`
- **Reaction 运行时**：`spyReportStart({ type: 'reaction', ... })`

#### 3. 嵌套追踪

`spyReportStart` 和 `spyReportEnd` 是成对出现的。这使得 DevTools 可以构建出**调用栈树**。

```typescript
action A
  update X
  action B
    update Y
```

DevTools 监听这些事件，就能画出那个漂亮的依赖图和性能火焰图。

---

### 四、 总结：MobX 源码的全景图

至此，我们已经把 MobX 源码拆解得差不多了。让我们最后用一张逻辑图来串联所有知识点：

1.  **入口**：`makeObservable` -> 创建 `ObservableValue` (持有 `Atom`)。
2.  **触发**：`autorun` -> 创建 `Reaction` (持有 `Derivation`)。
3.  **读取**：`Reaction` 执行 -> 读 `ObservableValue` -> `Atom.reportObserved` -> **依赖收集** (双向链表)。
4.  **修改**：`Action` 执行 -> 改 `ObservableValue` -> `Atom.reportChanged` -> **标记脏状态** (STALE)。
5.  **调度**：`Transaction` 结束 -> 遍历 `pendingReactions` -> `Reaction.runReaction`。
6.  **优化**：
    - `ComputedValue` -> **状态机** (POSSIBLY_STALE) 避免无效计算。
    - `ObservableArray` -> **Stub & Splice** 统一管理变异。
    - `Flow` -> **Generator + Action** 解决异步上下文丢失。

这就是 MobX：**一个用同步代码模拟出的、自动化的、高性能的电子表格系统**。

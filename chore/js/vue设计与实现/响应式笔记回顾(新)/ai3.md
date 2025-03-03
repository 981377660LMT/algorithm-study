下面的内容将结合 **Vue 3.x** 的源码思路（基于 [@vue/reactivity](https://github.com/vuejs/core/tree/main/packages/reactivity)）以及部分 **Vue 2.x** 的实现原理（基于 `Object.defineProperty`）进行综合讲解，让你对 Vue 响应式系统的核心原理有一个清晰的认识。其中会穿插简化版的示例代码，以更直观地展示 Vue 内部是如何完成依赖收集与响应式更新的。

---

## 一、Vue 3 响应式系统概览

Vue 3 使用了 **Proxy** 来实现数据拦截与追踪。其核心模块被抽离到 [@vue/reactivity](https://github.com/vuejs/core/tree/main/packages/reactivity) 中，包含以下重要概念和方法：

1. **响应式创建函数**
   - `reactive()`：将普通对象转换为响应式对象
   - `ref()`：为基本数据类型创建响应式容器
2. **依赖收集与触发**
   - `track()`：依赖收集（在“读”数据时进行）
   - `trigger()`：派发更新（在“写”数据时进行）
3. **副作用函数 (effect)**
   - `effect(fn)`：将一个函数 `fn` 包装为副作用函数，并在其内部读取响应式数据时完成依赖收集。数据变更时会自动触发该函数再次执行。

当你在应用中使用了 `reactive` 或 `ref` 等 API 后，Vue 就会对相应的数据进行“代理”或“包装”，从而在访问（get）或修改（set）时执行必要的逻辑。具体流程可概括为：

1. 访问响应式数据时，通过 **Proxy get 拦截** 或 `Object.defineProperty` 的 getter，调用 `track()` 进行依赖收集。
2. 修改响应式数据时，通过 **Proxy set 拦截** 或 `Object.defineProperty` 的 setter，调用 `trigger()` 去通知所有收集到的副作用函数执行。

---

## 二、Vue 2.x 和 Vue 3.x 的区别

- **Vue 2.x**：基于 `Object.defineProperty` 实现数据劫持，对对象的每个属性都进行遍历并定义 getter/setter，进而在 getter 和 setter 中进行依赖收集和派发更新。对于数组的变更则进行了特殊处理（重写数组原型方法）。
- **Vue 3.x**：使用 **Proxy** 来统一拦截对象属性的读写、枚举、删除等操作，避免了 Vue 2.x 中的遍历开销、数组变异方法 hack 等问题。响应式系统更加通用和灵活。

由于 **Vue 3** 更具代表性，同时在《Vue.js 设计与实现》（霍春阳）中也主要是基于 Vue 3 的思想进行讲解，下面将重点对 **Vue 3** 的实现进行详细解析。

---

## 三、数据劫持：`reactive()` 与 `ref()`

### 3.1 `reactive()` 如何创建响应式对象

`reactive` 函数的核心逻辑大致如下（简化示例代码）：

```js
// reactive.js（简化示例）
const proxyMap = new WeakMap() // 缓存已经创建过的Proxy

function createReactiveObject(target) {
  // 如果已经存在缓存的Proxy，则直接返回
  const existingProxy = proxyMap.get(target)
  if (existingProxy) {
    return existingProxy
  }

  // 否则，创建一个新的Proxy
  const proxy = new Proxy(target, {
    get(target, key, receiver) {
      // 依赖收集
      track(target, key)
      // Reflect.get 优点：可保证 this 正确
      const res = Reflect.get(target, key, receiver)
      // 如果拿到的值还是一个对象，递归转换为响应式
      return typeof res === 'object' && res !== null ? reactive(res) : res
    },
    set(target, key, value, receiver) {
      const oldValue = target[key]
      const result = Reflect.set(target, key, value, receiver)
      // 如果值发生了变化，触发更新
      if (oldValue !== value) {
        trigger(target, key)
      }
      return result
    }
  })

  // 将创建的 Proxy 存入缓存
  proxyMap.set(target, proxy)

  return proxy
}

export function reactive(target) {
  if (typeof target !== 'object' || target === null) {
    return target // 非对象直接返回
  }
  return createReactiveObject(target)
}
```

> 1. 首次调用 `reactive(obj)` 时，会创建一个新 `Proxy` 并缓存在 `proxyMap` 中，如果多次对同一个 `obj` 调用 `reactive`，会复用同一个 `Proxy`。
> 2. 在 `get` 拦截中调用 `track` 函数进行依赖收集，并且对返回的值进行递归 `reactive` 化。
> 3. 在 `set` 拦截中调用 `trigger` 通知依赖更新。

### 3.2 `ref()` 如何为基础类型创建响应式

`ref` 主要解决的是**基本类型**（number、string、boolean 等）的响应式需求。它内部通常返回一个包裹对象 `{ value: ... }`，然后对这个包裹对象执行与 `reactive` 类似的操作。在 Vue 3 中大致会这样实现：

```js
// ref.js（简化示例）
function RefImpl(rawValue) {
  this._value = rawValue
  this.__v_isRef = true
}

Object.defineProperty(RefImpl.prototype, 'value', {
  get() {
    track(this, 'value') // 依赖收集
    return this._value
  },
  set(newVal) {
    if (newVal !== this._value) {
      this._value = newVal
      trigger(this, 'value') // 触发更新
    }
  }
})

export function ref(rawValue) {
  return new RefImpl(rawValue)
}
```

> Vue 3 在实现细节上会使用 `Object.defineProperty` 或 `Proxy` 混合方式，但核心是：**`ref` 通过对象形式包裹原始值，并为 `value` 属性做依赖收集和派发更新**。

---

## 四、依赖收集：`track()`

### 4.1 为什么需要依赖收集

当我们在一个 “副作用函数” 中读取响应式数据时，Vue 要记录：**哪个副作用函数读取了哪些响应式数据**。这样当这部分数据更新时，就可以**精确**地触发对应的副作用函数执行，而不用去重新执行所有函数。

Vue 内部使用一个全局变量（如 `activeEffect`）来标记当前正在运行的副作用函数。然后在读操作时，会将 `activeEffect` 收集到对应的数据依赖中。

### 4.2 `track()` 简化示例

Vue 会使用一个数据结构（通常是 `targetMap`）保存目标对象、其属性、以及依赖它们的副作用函数之间的对应关系。可视化如下：

```
targetMap = {
  [targetObj]: {
    [keyName]: Set of effects
  },
  ...
}
```

示例实现：

```js
// effect.js (或 track.js)
const targetMap = new WeakMap()
let activeEffect = null // 当前正在执行的副作用函数

export function track(target, key) {
  if (!activeEffect) return

  // 尝试从 targetMap 中取到 depsMap
  let depsMap = targetMap.get(target)
  if (!depsMap) {
    depsMap = new Map()
    targetMap.set(target, depsMap)
  }

  // 再根据 key 取到对应的依赖集合
  let dep = depsMap.get(key)
  if (!dep) {
    dep = new Set()
    depsMap.set(key, dep)
  }

  // 将当前激活的副作用函数加入到依赖集合
  dep.add(activeEffect)
}
```

> 当在副作用函数中访问到某个 `target` 对象的某个 `key` 属性时，就会执行 `track(target, key)`，将该副作用函数记录到 `dep` 集合中。

---

## 五、派发更新：`trigger()`

### 5.1 触发更新的原理

当某个数据变更时，需要通知所有依赖该数据的副作用函数，让它们重新执行以获取最新的值。这就是 `trigger` 的职责。

### 5.2 `trigger()` 简化示例

```js
export function trigger(target, key) {
  // 从 targetMap 中取到依赖集合
  const depsMap = targetMap.get(target)
  if (!depsMap) return

  const dep = depsMap.get(key)
  if (!dep) return

  // 逐个执行收集到的副作用函数
  dep.forEach(effectFn => {
    // 这里可根据不同情况做调度策略
    effectFn()
  })
}
```

---

## 六、副作用函数：`effect(fn)`

### 6.1 如何收集依赖

- 在执行 `effect(fn)` 时，Vue 会先将 `activeEffect` 指向 `fn`，然后执行 `fn`，从而在 `fn` 内部访问响应式数据时触发 `track`，完成依赖收集。
- 每次执行 `fn` 前后，Vue 会将 `activeEffect` 恢复或清除，防止干扰其他操作。

### 6.2 简化实现

```js
export function effect(fn) {
  // 将 effect 包装一下
  const effectFn = () => {
    activeEffect = effectFn // 设置激活的副作用函数
    fn() // 执行用户传入的函数
    activeEffect = null // 重置
  }

  // 默认执行一次
  effectFn()
  return effectFn
}
```

#### 使用示例

```js
import { reactive, effect } from './reactivity'

const state = reactive({ count: 0 })

effect(() => {
  console.log('count changed:', state.count)
})

// 改变数据，触发更新
state.count++
// console 输出：count changed: 1
```

当 `state.count++` 触发 `set` 拦截时，调用 `trigger(state, 'count')`，进而调用先前收集到的 effect 函数，从而打印出新的值。

---

## 七、调度策略与异步队列

在真实的 Vue 应用中，如果同步地立即执行副作用，会导致某些场景下触发多次重复计算或渲染。为提高性能，Vue 通常会采用 **调度器** 或 **异步队列** 的策略：

1. 当多次触发同一 effect 时，会将 effect 放入一个队列中去重，异步地执行一次批量更新。
2. 具体实现里常用 **微任务**（`Promise.resolve().then()`）或 **宏任务** 实现 “下一次事件循环” 中执行批量更新。
3. Vue 2.x 用 `nextTick`，Vue 3.x 也是类似原理，有一个内置的 scheduler，默认会异步执行。我们也可以在 `effect()` 时自定义调度器函数。

示例（简化）：

```js
let jobQueue = new Set()
let isFlushing = false

function flushJob() {
  if (isFlushing) return
  isFlushing = true
  Promise.resolve().then(() => {
    jobQueue.forEach(job => job())
    jobQueue.clear()
    isFlushing = false
  })
}

export function trigger(target, key) {
  // ...同前述trigger逻辑，拿到dep...
  dep.forEach(effectFn => {
    jobQueue.add(effectFn) // 放入队列
  })
  flushJob() // 尝试执行
}
```

> 这样当数据多次变化时，副作用函数只会在下一次事件循环中执行一次，避免重复渲染，极大提升了性能。

---

## 八、Vue 2.x 中的核心类：`Observer`、`Dep`、`Watcher`

虽然 Vue 3.x 的实现更简洁，但了解 Vue 2.x 也有助于理解它们的相同思路。Vue 2.x 在内部定义了若干类：

1. **Observer**：遍历对象的属性，通过 `Object.defineProperty` 定义 getter/setter 进行数据劫持，并将自己挂载到 `__ob__` 属性上。
2. **Dep**：收集依赖并派发更新，其内部包含 `subs` 数组，用来存放依赖该数据的所有 “订阅者”(Watcher)。
3. **Watcher**：即副作用逻辑的封装，每一个组件渲染函数会对应一个渲染 Watcher，每一个用户自定义的 watch 也会有一个 Watcher。

简要流程图：

```
data() -> Observer -> 对每个属性 defineProperty -> Dep
                                ↑       ↓
                              Watcher(渲染/用户watch)
```

- **Observer**：当访问数据时触发 `getter`，`getter` 内部执行 `dep.depend()`；当修改数据时触发 `setter`，`setter` 内部执行 `dep.notify()`。
- **Dep**：`depend()` 时，将当前活动的 `Watcher` 存进 `subs`；`notify()` 时，遍历 `subs` 调用每个 `Watcher` 的 `update()`。
- **Watcher**：`update()` 内部会进行重新求值等操作，如果是渲染 Watcher，就会再次执行渲染逻辑。

---

## 九、渲染更新与 Virtual DOM

在 Vue 的组件系统中，**渲染函数**（或模板编译生成的渲染函数）会读取组件 `data` 中的数据。当这些数据被转换成响应式之后：

1. 初次渲染时，渲染函数会访问响应式数据，从而收集到渲染 Watcher（Vue 2.x）或渲染 effect（Vue 3.x）。
2. 当数据改变时，触发依赖更新，再次执行渲染函数。
3. 在 Vue 2.x 中，模板编译生成的渲染函数会产生 Virtual DOM，并调用相应的 diff 算法更新真实 DOM。
4. 在 Vue 3.x 中，Virtual DOM 同样存在，diff 算法也有提升，但核心理念一致。

即使实现方式不同，本质都是**“数据 -> 渲染函数 -> Virtual DOM -> 真实DOM”的更新流程**，而响应式系统只是确保当数据发生变化时，可以再次触发这条链路。

---

## 十、总结

1. **核心思想：拦截数据的“读”和“写”**
   - “读”时收集依赖（哪些副作用函数用到了它？）
   - “写”时触发更新（让这些副作用函数重新执行）
2. **Vue 2.x** 基于 `Object.defineProperty` + “遍历”方式 + 对数组方法做特殊处理。
3. **Vue 3.x** 基于 `Proxy` 统一拦截对象操作，代码实现更简洁、更通用。
4. **依赖管理** 通过全局激活的副作用函数 `activeEffect` + 数据结构 `targetMap`（或 `Dep`），将两者关联起来。
5. **派发更新** 通过异步队列或调度器的方式，批量执行而非同步多次执行。
6. **渲染更新** 结合 Virtual DOM 做最小化 DOM 更新，提升性能。

理解了以上流程，就能知道为什么在 Vue 中，只要把数据定义为响应式，当我们改动它时，相关视图就会自动刷新。**这并不是“魔法”，而是使用拦截 + 依赖收集 + 批量更新的设计模式**。

---

### 参考链接

- [Vue.js 官方源码仓库 (Vue 3)](https://github.com/vuejs/core)
- [Vue.js 设计与实现（霍春阳）](https://book.douban.com/subject/35582943/)
- [@vue/reactivity 源码](https://github.com/vuejs/core/tree/main/packages/reactivity)

如果你想要进一步研究，可直接阅读 [@vue/reactivity](https://github.com/vuejs/core/tree/main/packages/reactivity) 的源码，从 `reactive.ts`、`ref.ts`、`effect.ts`、`baseHandlers.ts` 等文件入手，与上面介绍的简化逻辑对比，即可获得深入理解。祝学习愉快！

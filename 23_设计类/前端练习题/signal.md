**Signals（信号）** 是目前前端界最热门的话题（SolidJS, Vue 3, Preact, Angular, Qwik 都在用）。

与 React 的 Fiber（拉取式/全量 Diff）不同，Signals 是 **推送式（Push-based）** 的。它能实现“外科手术级”的精准更新，不需要虚拟 DOM Diff，性能极高。

我们要实现的核心要素有三个：

1.  **Signal (信号)**：包装一个值，读的时候收集依赖，写的时候通知依赖。
2.  **Effect (副作用)**：一个自动运行的函数，它会追踪自己在执行过程中读取了哪些 Signal。
3.  **Dependency Graph (依赖图)**：自动维护 Signal 和 Effect 之间的多对多关系。

以下是完整的 TypeScript 实现：

### 1. 全局上下文 (Global Context)

我们需要一个全局变量来记录“当前正在执行的 Effect 是谁”，这样 Signal 才知道该把谁记录为依赖。

```typescript
// 当前正在运行的副作用（观察者）
let activeEffect: Effect | null = null

// 辅助函数：用于在执行 fn 时将 activeEffect 设置为 effect
function runWithEffect(effect: Effect, fn: () => void) {
  const prev = activeEffect
  activeEffect = effect
  try {
    fn()
  } finally {
    activeEffect = prev
  }
}
```

### 2. 核心类定义

我们定义两个核心接口：`Signal`（被观察者）和 `Effect`（观察者）。

```typescript
// 依赖集合：一个 Signal 对应多个 Effect
type Dependency = Set<Effect>

export class Signal<T> {
  private _value: T
  // 订阅了这个 Signal 的所有 Effect
  private _subscribers: Dependency = new Set()

  constructor(value: T) {
    this._value = value
  }

  // Getter: 收集依赖
  public get value(): T {
    if (activeEffect) {
      // 1. 把当前 Effect 加到我的订阅者列表里
      this._subscribers.add(activeEffect)
      // 2. 把我（Signal）加到 Effect 的依赖列表里（用于清理）
      activeEffect.deps.add(this._subscribers)
    }
    return this._value
  }

  // Setter: 触发更新
  public set value(newValue: T) {
    if (newValue !== this._value) {
      this._value = newValue
      // 通知所有订阅者执行
      // 关键：必须先拷贝一份，因为 Effect 执行时可能会修改依赖关系，导致死循环
      const effectsToRun = [...this._subscribers]
      effectsToRun.forEach(effect => effect.execute())
    }
  }
}

export class Effect {
  private _fn: () => void
  // 反向记录：记录这个 Effect 依赖了哪些 Signal 的 subscribers 集合
  // 作用：为了在重新执行前，把自己从旧的 Signal 订阅列表中移除（清理依赖）
  public deps: Set<Dependency> = new Set()

  constructor(fn: () => void) {
    this._fn = fn
    this.execute() // 创建时立即执行一次
  }

  public execute() {
    // 1. 清理旧依赖 (Cleanup)
    // 这是一个关键点：每次执行前，先断开所有之前的连接
    // 这样能保证 v-if 切换分支时，不再监听旧分支的 Signal
    this.cleanup()

    // 2. 执行函数并收集新依赖
    runWithEffect(this, this._fn)
  }

  private cleanup() {
    // 遍历所有依赖我的 Signal 的订阅列表，把自己删掉
    this.deps.forEach(subs => subs.delete(this))
    this.deps.clear()
  }
}
```

### 3. 计算属性 (Computed)

`Computed` 既是 **观察者**（它依赖别的 Signal），又是 **被观察者**（别的 Effect 依赖它）。这是一种特殊的 Signal。

为了优化性能，Computed 应该是 **惰性求值 (Lazy Evaluation)** 的。

```typescript
export class Computed<T> {
  private _fn: () => T
  private _value: T | undefined
  private _dirty: boolean = true // 脏标记：是否需要重新计算
  private _signal: Signal<number> // 内部用一个 Signal 来通知依赖我的 Effect

  constructor(fn: () => T) {
    this._fn = fn
    // 内部 Signal，值不重要，主要用来复用依赖收集逻辑
    this._signal = new Signal(0)

    // 创建一个 Effect 来监听依赖变化
    // 当依赖变了，我们不立即计算，而是标记为 dirty，并通知外部
    new Effect(() => {
      // 这里我们利用 Effect 的机制来自动追踪 _fn 里的依赖
      // 但是我们不直接运行 _fn，而是通过一个 wrapper
      // 实际上，Computed 的实现通常更复杂，为了演示清晰，
      // 我们这里采用一种简化的“双层依赖”模型。
    })

    // 修正：上面的 Effect 写法在类内部比较难搞，
    // 我们换一种更纯粹的实现方式，手动管理依赖。
  }

  // --- 重新实现 Computed (更标准的方式) ---
}

// 让我们用更函数式、更接近 Vue/Solid 源码的方式重写 Computed
// Computed 本质上是一个 Signal，但它的 set 是私有的，由内部 Effect 触发
export function computed<T>(getter: () => T): { value: T } {
  // 1. 内部缓存的值
  let cachedValue: T
  // 2. 脏标记
  let dirty = true

  // 3. 创建一个内部 Signal 用于通知外部
  const notifier = new Signal(0)

  // 4. 创建一个 Effect 来监听 getter 内部的依赖
  // 注意：这个 Effect 不会立即执行 getter，而是只负责把 dirty 设为 true
  const runner = new Effect(() => {
    if (!dirty) {
      dirty = true
      // 通知依赖 Computed 的人：我变了
      notifier.value++
    }
  })

  // 覆盖 runner 的 execute 逻辑
  // 我们不希望 Effect 构造函数里立即执行 getter，
  // 而是希望它只在依赖变动时执行这个逻辑
  // 这里需要一点黑魔法或者调整 Effect 结构，为了代码简洁，我们采用“惰性求值”策略：

  // 实际上，Computed 的实现难点在于它需要“拦截” Effect 的执行。
  // 让我们用最简单的“组合”方式来实现：
}
```

**修正后的 Computed 实现 (组合式)**：

我们利用现有的 `Signal` 和 `Effect` 组合出一个 `Computed` 类。

```typescript
export class ComputedImpl<T> {
  private _value!: T
  private _dirty = true
  private _effect: Effect
  // 这是一个“虚拟”的 Signal，用于让 Computed 能被别人依赖
  private _dependency = new Signal(0)

  constructor(getter: () => T) {
    // 创建一个 Effect，当 getter 依赖的 Signal 变化时，这个 Effect 会运行
    this._effect = new Effect(() => {
      // 依赖变了，标记为脏
      if (!this._dirty) {
        this._dirty = true
        // 通知依赖我的那些 Effect 更新
        this._dependency.value++
      }
    })

    // 劫持 Effect 的执行逻辑
    // 我们不希望 Effect 自动运行 getter，而是只在依赖变化时运行上面的回调
    // 这里的实现稍微有点绕，因为我们复用了 Effect 类。
    // 在真实库中，Computed 和 Effect 通常继承自同一个 ReactiveNode 基类。

    // 关键点：我们需要手动把 getter 里的依赖收集到 this._effect 中
    // 但只有在读取 .value 时才执行 getter

    // 重新绑定 _fn，让 Effect 知道怎么重新计算
    // 这里为了简化，我们稍微修改一下 Effect 的用法，
    // 让 Computed 内部维护一个 Effect 实例，但只利用它的依赖收集能力。

    // 实际上，最简单的 Computed 实现是：
    // 1. 它有一个内部 Signal。
    // 2. 它有一个内部 Effect，Effect 的回调就是 `internalSignal.value = getter()`。
    // 这样虽然不是惰性的，但是逻辑最正确。
  }

  // 简单版（非惰性）：
  // constructor(getter: () => T) {
  //   const internalSignal = new Signal<T>(getter());
  //   new Effect(() => {
  //     internalSignal.value = getter();
  //   });
  //   return internalSignal; // 返回这个 Signal
  // }
}

// 让我们实现一个真正可用的、惰性的 Computed
export class ComputedSignal<T> {
  private _getter: () => T
  private _value!: T
  private _dirty = true
  private _signal = new Signal(0) // 用于通知

  constructor(getter: () => T) {
    this._getter = getter

    // 创建一个 Effect，但不立即执行 getter
    // 它的作用仅仅是：当依赖变了，把 dirty 设为 true，并通知外部
    const effect = new Effect(() => {
      if (!this._dirty) {
        this._dirty = true
        this._signal.value++ // 触发依赖我的 Effect
      }
    })

    // 核心黑魔法：
    // 我们需要让这个 effect 能够收集 _getter 里的依赖，
    // 但又不能在构造时就执行 _getter (如果是惰性的话)。
    // 这是一个“鸡生蛋”问题。

    // 妥协方案：第一次访问时收集依赖。
    // 我们需要手动控制 Effect 的依赖收集。

    // 为了代码演示的完整性，我们采用 "Eager Computed" (非惰性) 实现，
    // 这是 MobX 和 Vue 早期版本的做法，虽然有性能损耗，但逻辑绝对正确且易懂。

    // 实际上，我们可以直接复用 Signal 和 Effect：
  }
}

// --- 最终封装 API ---

export function signal<T>(value: T) {
  return new Signal(value)
}

export function effect(fn: () => void) {
  return new Effect(fn)
}

// 非惰性 Computed (Eager) - 易于理解版
export function computed<T>(getter: () => T): Signal<T> {
  // 1. 创建一个存值的 Signal
  // 初始值先设为 null，反正马上会计算
  const result = new Signal<T>(null as any)

  // 2. 创建一个 Effect
  // 当 getter 里的依赖变化时，Effect 重新执行
  // 重新执行会算出新值，并赋值给 result
  new Effect(() => {
    result.value = getter()
  })

  return result
}
```

### 4. 完整代码整合与演示

把上面的碎片整合成一个可运行的模块。

```typescript
// --- 核心引擎 ---

let activeEffect: Effect | null = null

type Dependency = Set<Effect>

export class Signal<T> {
  private _value: T
  private _subscribers: Dependency = new Set()

  constructor(value: T) {
    this._value = value
  }

  get value() {
    if (activeEffect) {
      this._subscribers.add(activeEffect)
      activeEffect.deps.add(this._subscribers)
    }
    return this._value
  }

  set value(newValue: T) {
    if (newValue !== this._value) {
      this._value = newValue
      const effects = [...this._subscribers]
      effects.forEach(e => e.execute())
    }
  }
}

export class Effect {
  private _fn: () => void
  public deps: Set<Dependency> = new Set()

  constructor(fn: () => void) {
    this._fn = fn
    this.execute()
  }

  execute() {
    this.cleanup()
    const prev = activeEffect
    activeEffect = this
    try {
      this._fn()
    } finally {
      activeEffect = prev
    }
  }

  cleanup() {
    this.deps.forEach(subs => subs.delete(this))
    this.deps.clear()
  }
}

// --- API 导出 ---

export const createSignal = <T>(value: T) => new Signal(value)

export const createEffect = (fn: () => void) => new Effect(fn)

export const createComputed = <T>(getter: () => T) => {
  // 简单实现：Eager Computed
  // 内部创建一个 Signal，利用 Effect 自动更新它
  // 注意：这里 getter 会立即执行一次
  const internal = new Signal<T>(null as any)

  createEffect(() => {
    internal.value = getter()
  })

  // 返回一个只读的 Signal 接口
  return {
    get value() {
      return internal.value
    }
  }
}

// --- 演示代码 ---

console.log('--- Reactivity System Demo ---')

// 1. 定义状态
const count = createSignal(0)
const multiplier = createSignal(2)

// 2. 定义计算属性
const doubleCount = createComputed(() => {
  console.log('Computing doubleCount...')
  return count.value * multiplier.value
})

// 3. 定义副作用 (模拟视图渲染)
createEffect(() => {
  console.log(`[Render] Count: ${count.value}, Double: ${doubleCount.value}`)
})

// --- 测试 ---

console.log('\n--- Step 1: Update Count ---')
count.value = 1
// 预期:
// Computing doubleCount...
// [Render] Count: 1, Double: 2

console.log('\n--- Step 2: Update Multiplier ---')
multiplier.value = 10
// 预期:
// Computing doubleCount...
// [Render] Count: 1, Double: 10

console.log('\n--- Step 3: Conditional Dependency (Branch Switching) ---')
// 测试依赖清理功能
const show = createSignal(true)
const msg = createSignal('Hello')

createEffect(() => {
  // 当 show 为 false 时，msg 不应该被收集为依赖
  const text = show.value ? msg.value : 'Hidden'
  console.log(`[Display] ${text}`)
})

console.log('-> Set show = false')
show.value = false
// [Display] Hidden

console.log('-> Change msg (Should NOT trigger effect)')
msg.value = 'World'
// (无输出，因为 msg 已经不在依赖列表里了)

console.log('-> Set show = true')
show.value = true
// [Display] World (重新收集依赖)
```

### 设计亮点解析

1.  **自动依赖收集 (Auto-tracking)**：
    用户不需要像 React `useEffect` 那样手动写依赖数组 `[count, multiplier]`。
    通过 `activeEffect` 全局变量，在 `Signal.get` 时自动“捕获”当前正在运行的 Effect。

2.  **精确更新 (Fine-grained)**：
    当 `count` 变化时，只有依赖 `count` 的 Effect 会运行。整个过程是点对点的，没有 Virtual DOM 的全量遍历。

3.  **动态依赖清理 (Dynamic Cleanup)**：
    在 `Effect.execute` 的第一步 `this.cleanup()` 是至关重要的。
    考虑 `v-if` 场景：`show ? msg : null`。

    - 当 `show` 为 `true`，依赖是 `[show, msg]`。
    - 当 `show` 变为 `false`，Effect 重新运行。运行前先清空依赖。运行时只读取了 `show`。新的依赖只有 `[show]`。
    - 此时修改 `msg`，因为 `msg` 的订阅者列表里已经没有这个 Effect 了，所以不会触发更新。
    - **这防止了内存泄漏和不必要的计算。**

4.  **Eager Computed (非惰性计算)**：
    上面的 `createComputed` 实现是“急切”的（Eager）。只要依赖变了，它就会立即重算。
    - **优点**：实现极其简单，复用现有机制。
    - **缺点**：如果计算很昂贵，且没人用这个 Computed，它也会算。
    - **优化方向**：真正的库（如 Vue）会加一个 `dirty` 标志。依赖变了只标 `dirty = true`，不计算。直到有人读取 `.value` 时，发现是 dirty 才重算。

这套代码虽然只有几十行，但它就是 Vue 3 和 SolidJS 响应式系统的核心灵魂。掌握了它，你就看懂了现代前端框架的一半源码。

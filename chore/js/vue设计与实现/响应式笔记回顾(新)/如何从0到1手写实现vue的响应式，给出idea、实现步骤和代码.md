# 从零实现Vue响应式系统

## 核心思路

Vue响应式系统的核心在于三个关键点：

1. **数据劫持**：拦截数据的读写操作
2. **依赖收集**：追踪谁依赖了这些数据
3. **派发更新**：数据变化时通知依赖者更新

## 实现步骤

### 第一步：创建依赖收集器

首先，我们需要一个全局存储依赖关系的结构和活跃依赖项的追踪机制。

```javascript
// 全局变量跟踪当前活跃的effect函数
let activeEffect = null

// 存储所有响应式对象的依赖关系，使用WeakMap避免内存泄漏
const targetMap = new WeakMap()

// 收集依赖
function track(target, key) {
  if (!activeEffect) return

  let depsMap = targetMap.get(target)
  if (!depsMap) {
    targetMap.set(target, (depsMap = new Map()))
  }

  let dep = depsMap.get(key)
  if (!dep) {
    depsMap.set(key, (dep = new Set()))
  }

  // 添加到依赖集合
  dep.add(activeEffect)
  // 反向收集，让effect知道自己依赖了哪些集合
  activeEffect.deps.push(dep)
}

// 触发更新
function trigger(target, key) {
  const depsMap = targetMap.get(target)
  if (!depsMap) return

  const dep = depsMap.get(key)
  if (!dep) return

  // 创建一个新的Set进行遍历，防止在遍历过程中Set发生变化
  const effectsToRun = new Set()

  dep.forEach(effectFn => {
    // 如果effect正在运行，不要再次触发它
    if (effectFn !== activeEffect && !effectFn.running) {
      effectsToRun.add(effectFn)
    }
  })

  effectsToRun.forEach(effectFn => {
    if (effectFn.scheduler) {
      effectFn.scheduler()
    } else {
      effectFn()
    }
  })
}
```

### 第二步：实现响应式数据

使用ES6的Proxy创建响应式对象：

```javascript
// 创建响应式对象
function reactive(target) {
  if (typeof target !== 'object' || target === null) {
    return target
  }

  const handler = {
    get(target, key, receiver) {
      const result = Reflect.get(target, key, receiver)

      // 依赖收集
      track(target, key)

      // 如果获取的值是对象，继续将其变成响应式
      return typeof result === 'object' && result !== null ? reactive(result) : result
    },
    set(target, key, value, receiver) {
      const oldValue = target[key]
      const result = Reflect.set(target, key, value, receiver)

      // 只有值变化时才触发更新
      if (result && oldValue !== value) {
        // 触发依赖更新
        trigger(target, key)
      }
      return result
    },
    deleteProperty(target, key) {
      const hadKey = Object.prototype.hasOwnProperty.call(target, key)
      const result = Reflect.deleteProperty(target, key)

      // 如果删除成功且对象上确实有该属性，触发更新
      if (hadKey && result) {
        trigger(target, key)
      }
      return result
    }
  }

  return new Proxy(target, handler)
}
```

### 第三步：实现effect副作用函数

```javascript
function effect(fn, options = {}) {
  const effectFn = () => {
    try {
      activeEffect = effectFn
      // 执行传入的函数，期间会触发依赖收集
      return fn()
    } finally {
      // 执行完毕后清除当前活跃effect
      activeEffect = null
    }
  }

  if (!options.lazy) {
    // 默认立即执行一次
    effectFn()
  }

  return effectFn // 返回effect函数以便后续调用
}
```

### 第四步：实现ref

对于简单值类型，需要用对象包裹后才能实现响应式：

```javascript
function ref(value) {
  // 创建包装对象
  const refObject = {
    get value() {
      track(refObject, 'value')
      return value
    },
    set value(newValue) {
      if (value !== newValue) {
        value = newValue
        trigger(refObject, 'value')
      }
    }
  }

  return refObject
}
```

### 第五步：实现computed计算属性

```javascript
function computed(getter) {
  // 用于缓存计算结果
  let value
  // 脏标记，表示是否需要重新计算
  let dirty = true

  // 创建effect
  const effectFn = effect(() => getter(), {
    lazy: true,
    // 当依赖变化时触发的回调
    scheduler: () => {
      if (!dirty) {
        dirty = true
        // 通知依赖此计算属性的effect更新
        trigger(computedRef, 'value')
      }
    }
  })

  const computedRef = {
    get value() {
      if (dirty) {
        value = effectFn()
        dirty = false
      }
      // 当有人读取计算属性时，建立依赖关系
      track(computedRef, 'value')
      return value
    }
  }

  return computedRef
}
```

### 第六步：实现watch监听器

```javascript
function watch(source, callback, options = {}) {
  let getter

  // 处理不同类型的监听源
  if (typeof source === 'function') {
    getter = source
  } else {
    getter = () => traverse(source)
  }

  let oldValue
  let cleanup

  // onInvalidate函数，用于注册清理回调
  function onInvalidate(fn) {
    cleanup = fn
  }

  const job = () => {
    const newValue = effectFn()

    // 调用清理函数
    if (cleanup) {
      cleanup()
    }

    // 调用回调函数
    callback(newValue, oldValue, onInvalidate)
    oldValue = newValue
  }

  const effectFn = effect(() => getter(), {
    lazy: true,
    scheduler: job
  })

  // 处理立即执行选项
  if (options.immediate) {
    job()
  } else {
    oldValue = effectFn()
  }
}

// 递归遍历对象的所有属性，用于深度监听
function traverse(value, seen = new Set()) {
  // 如果是原始值或已经遍历过，则返回
  if (typeof value !== 'object' || value === null || seen.has(value)) {
    return value
  }

  seen.add(value)

  // 递归遍历所有属性
  if (Array.isArray(value)) {
    for (let i = 0; i < value.length; i++) {
      traverse(value[i], seen)
    }
  } else {
    for (const key of Object.keys(value)) {
      traverse(value[key], seen)
    }
  }

  return value
}
```

### 第七步：添加批量更新（类似Vue的nextTick）

```javascript
const queue = new Set()
let isFlushing = false
const p = Promise.resolve()

function queueJob(job) {
  queue.add(job)

  if (!isFlushing) {
    isFlushing = true

    p.then(() => {
      try {
        queue.forEach(job => job())
      } finally {
        isFlushing = false
        queue.clear()
      }
    })
  }
}

// 修改trigger函数支持批量更新
function trigger(target, key) {
  const depsMap = targetMap.get(target)
  if (!depsMap) return

  const dep = depsMap.get(key)
  if (dep) {
    dep.forEach(effectFn => {
      if (effectFn.scheduler) {
        effectFn.scheduler()
      } else {
        queueJob(effectFn)
      }
    })
  }
}

// 修改effect函数支持调度器
function effect(fn, options = {}) {
  const effectFn = () => {
    // 清除此effect对之前依赖的引用，防止不必要的更新
    cleanup(effectFn)

    try {
      activeEffect = effectFn
      // 添加标识，表示此effect正在执行
      effectFn.running = true
      return fn()
    } finally {
      activeEffect = null
      effectFn.running = false
    }
  }

  // 用于存储此effect依赖的所有集合
  effectFn.deps = []
  effectFn.scheduler = options.scheduler

  if (!options.lazy) {
    effectFn()
  }

  return effectFn
}

// 清理effect依赖
function cleanup(effectFn) {
  // 从所有依赖集合中移除此effect
  for (let i = 0; i < effectFn.deps.length; i++) {
    const dep = effectFn.deps[i]
    dep.delete(effectFn)
  }
  effectFn.deps.length = 0
}
```

### 完整使用示例

```javascript
// 创建响应式对象
const state = reactive({
  count: 0,
  message: 'Hello'
})

// 使用effect跟踪变化
effect(() => {
  console.log(`Count is: ${state.count}`)
})

// 创建计算属性
const doubledCount = computed(() => state.count * 2)

// 创建监听器
watch(
  () => state.count,
  (newVal, oldVal) => {
    console.log(`Count changed from ${oldVal} to ${newVal}`)
  }
)

// 修改数据，自动触发更新
state.count++ // 会触发effect、computed和watch
```

以上就是一个简化但功能完整的Vue响应式系统实现。实际Vue的源码会有更多边缘情况处理、性能优化和API扩展，但核心原理就是这些。

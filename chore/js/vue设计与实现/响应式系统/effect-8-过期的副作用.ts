// effect-8-过期的副作用
// 竞态问题，连续发送两次请求，后面的先返回，导致先发送的返回结果覆盖了后面的请求。而一般的需求是，保留最后一次请求的结果。
// 类似取消Promise的操作，我们可以通过一个标识来判断是否需要执行副作用函数。
// !假设在 watch 我们会发送一个异步请求，可以通过在 watch 的回调函数中新增一个 onInvalidate 参数解决。
//
// !先返回的请求可以使用onInvalidate函数作用到后返回的请求上

interface IEffectFn<R = any> {
  (): R
  deps: Set<() => void>[]
  options?: IEffectOptions
}

interface IEffectOptions {
  scheduler?: (fn: () => void) => void
  lazy?: boolean
}

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: IEffectFn | undefined

const effectStack: IEffectFn[] = []

/** 用于注册副作用函数. */
function effect<R = any>(fn: () => R, options: IEffectOptions = {}): IEffectFn<R> {
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    effectStack.push(effectFn)
    const res = fn()
    effectStack.pop()
    activeEffect = effectStack.length > 0 ? effectStack[effectStack.length - 1] : undefined
    return res
  }

  /** 把 options 挂在effectFn上. */
  effectFn.options = options
  /** 用来存储与该副作用函数关联的依赖集合. */
  effectFn.deps = [] as IEffectFn['deps']

  if (!options.lazy) {
    effectFn()
  }
  return effectFn
}

function cleanup(effectFn: IEffectFn): void {
  /** 在每个依赖集合中把该函数删除. */
  for (let i = 0; i < effectFn.deps.length; i++) {
    const depsSet = effectFn.deps[i]
    depsSet.delete(effectFn)
  }
  effectFn.deps.length = 0
}

const bucket = new WeakMap<object, Map<PropertyKey, Set<IEffectFn>>>() // target -> target key -> Set<副作用函数>
const data = { age: 1 }

/** reactive. */
const reactiveData = new Proxy(data, {
  get(target, key: string) {
    track(target, key)
    // @ts-ignore
    return target[key]
  },
  set(target, key, newValue) {
    // @ts-ignore
    target[key] = newValue
    trigger(target, key)
    return true
  }
})

/** 在track中记录deps. */
function track(target: object, key: PropertyKey): void {
  if (!activeEffect) return
  let depsMap = bucket.get(target)
  if (!depsMap) {
    depsMap = new Map()
    bucket.set(target, depsMap)
  }
  let depsSet = depsMap.get(key)
  if (!depsSet) {
    depsSet = new Set()
    depsMap.set(key, depsSet)
  }
  depsSet.add(activeEffect)

  /** 当前副作用函数也记录下关联的依赖. */
  activeEffect.deps.push(depsSet)
}

function trigger(target: object, key: PropertyKey): void {
  const depsMap = bucket.get(target)
  if (!depsMap) return
  const depsSet = depsMap.get(key)
  if (!depsSet) return

  const effectsToRun = new Set<IEffectFn>()
  depsSet.forEach(effect => {
    // 在 trigger 中执行副作用函数的时候，不执行当前正在处理的副作用函数，即 activeEffect
    if (effect !== activeEffect) {
      effectsToRun.add(effect)
    }
  })

  /** 如果一个副作用函数存在调度器 就用调度器执行副作用函数. */
  effectsToRun.forEach(effect => {
    if (effect.options?.scheduler) {
      effect.options.scheduler(effect)
    } else {
      effect()
    }
  })
}

function tranverse(value: unknown, visited = new Set()): any {
  if (typeof value !== 'object' || value === null || visited.has(value)) {
    return
  }
  visited.add(value)
  if (Array.isArray(value)) {
    value.forEach(item => tranverse(item, visited))
  }
  if (value instanceof Map) {
    value.forEach((_, key) => tranverse(key, visited))
  }
  if (value instanceof Set) {
    value.forEach(item => tranverse(item, visited))
  }
  for (const key in value) {
    tranverse((value as Record<any, any>)[key], visited)
  }
  return value
}

export {}

interface IWathcOptions {
  /**
   * 在第一次执行获取初始值之后也会立即执行回调函数，不过第一次的 oldValue 是 undefined
   */
  immediate?: boolean

  /**
   * 除了立即执行，我们也可以通过其他方式指定回调函数执行时机.
   * 这样的话会导致连续修改数据的时候，执行结果有点问题，毕竟相当是数据全部改变后统一执行回调函数了.
   */
  flush?: 'post' | 'sync'
}

function watch<T>(
  source: () => T,
  cb: (newValue: T, oldValue: T | undefined, onInvaliate: (fn: () => void) => void) => void,
  options: IWathcOptions = {}
) {
  let getter: () => T

  // 判断用户传入的参数是否是函数
  if (typeof source === 'function') {
    getter = source
  } else {
    getter = () => tranverse(source)
  }

  let oldValue: T
  let newValue: T

  /** 用于存储用户注册的过期回调. */
  let cleanup: () => void
  function onInvaliate(fn: () => void) {
    cleanup = fn
  }

  // 把scheduler调度函数提取为job函数
  // 在 scheduler 手动调用副作用函数，获取最新的值并缓存，然后在回调时传入。
  // 这里使用了 lazy，是为了手动调用第一次副作用函数以获取 oldValue
  const job = () => {
    newValue = effectFn()

    // !在执行回调函数之前 先执行过期函数
    // 我们在回调函数中会调用失效函数 会把过期函数绑在cleanup上
    // !我们先调用的回调会先把失效函数绑定
    // 而如果在上一次回调函数执行之前 就触发了下一次的执行 就会调用失效函数
    // 也就是上一次的回调函数对应的失效函数 则上一次的结果会被取消
    if (cleanup) {
      cleanup()
    }
    cb(newValue, oldValue, onInvaliate)
    oldValue = newValue
  }

  const effectFn = effect(() => getter(), {
    lazy: true,
    scheduler: () => {
      if (options.flush === 'post') {
        const p = Promise.resolve()
        p.then(job)
      } else {
        job() // 同步执行
      }
    }
  })

  if (options.immediate) {
    job()
  } else {
    oldValue = effectFn()
  }
}

if (require.main === module) {
  let t = 0
  const mock = () => {
    return new Promise(resolve => {
      if (++t <= 2) {
        // 模拟一下 前两次需要1s返回 第3次立即返回
        t = 1000
      }
      setTimeout(() => {
        resolve(1)
      }, t)
    })
  }

  let finalData: any

  watch(
    () => reactiveData.age,
    async (newValue, oldValue, onInvalidate) => {
      console.log('foo 的值变了', newValue, oldValue)
      let expired = false
      onInvalidate(() => {
        expired = true
      })
      const res = await mock()
      console.log(expired ? '过期了' : '未过期', 'newValue=' + newValue)
      if (!expired) {
        finalData = res
      }
    }
  )
}

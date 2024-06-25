// https://www.cnblogs.com/wenruo/p/17050995.html

interface IEffectFn {
  (): void
  deps: Set<() => void>[]
  options?: IEffectOptions
}

interface IEffectOptions {
  scheduler?: (fn: () => void) => void
}

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: IEffectFn | undefined

const effectStack: IEffectFn[] = []

/** 用于注册副作用函数. */
function effect(fn: () => void, options: IEffectOptions = {}) {
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    effectStack.push(effectFn)
    fn()
    effectStack.pop()
    activeEffect = effectStack.length > 0 ? effectStack[effectStack.length - 1] : undefined
  }

  /** 把 options 挂在effectFn上. */
  effectFn.options = options
  /** 用来存储与该副作用函数关联的依赖集合. */
  effectFn.deps = [] as IEffectFn['deps']

  effectFn()
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

export {}

if (require.main === module) {
  effect(
    () => {
      console.log({ age: reactiveData.age })
    },
    { scheduler: fn => setTimeout(fn, 1000) }
  )

  reactiveData.age++
  console.log('结束')
}

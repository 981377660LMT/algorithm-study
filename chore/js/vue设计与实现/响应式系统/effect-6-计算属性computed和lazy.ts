// https://www.cnblogs.com/wenruo/p/17050995.html
// 如果我们有时希望副作用函数不要立即执行，则需要提供一个选项，lazy 来决定是否立即执行
// !函数只有在不指定 lazy 的时候才立即执行，同时把函数返回，这样就可以把函数在外部获取并随时手动执行。
// 同时 effectFn 函数返回了函数执行的结果。
// 只要依赖不改变，计算属性的值是不会变的

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

export {}

// 每一次读取值，都会触发 getter 的执行，但其实只要依赖不改变，计算属性的值是不会变的。
// !这时我们一方面可以设置一个 flag 标识是否需要重新计算，同时在依赖修改时，
// !只需要更新这个 flag 即可，并不需要重新计算，只需要在读取时计算即可
function computed<R>(getter: () => R) {
  let value: R
  let dirty = true // 是否需要重新计算
  const effectFn = effect(getter, {
    lazy: true,
    // 神来之笔啊！如果依赖修改了 并不需要重新计算 getter 但是需要更新 dirty
    // 只需要在 scheduler 中指定调度方式即可
    scheduler() {
      if (!dirty) {
        dirty = true
        // 在另一个 effect 中读取计算属性的时候，需要让计算属性收集依赖
        // !因为当 computed 的依赖改变时 computed 的值应该被重新计算
        // 这个时候需要手动触发依赖
        // 如果有依赖的话 就会重新计算 computed 的值了
        trigger(obj, 'value')
      }
    }
  })
  const obj = {
    get value() {
      if (dirty) {
        value = effectFn()
        dirty = false
      }
      // 读取value时手动触发依赖收集
      track(obj, 'value')
      return value
    }
  }
  return obj
}

// !在读取的时候做了依赖收集，同时在 computed 的依赖改变时，对收集的依赖触发执行

if (require.main === module) {
  const sumRes = computed(() => {
    console.log('computed recalced')
    return reactiveData.age + 1
  })

  effect(() => {
    console.log(sumRes.value)
  })

  console.log(sumRes.value)
  console.log(sumRes.value)
  reactiveData.age = 2
  console.log(sumRes.value)
  console.log(sumRes.value)
}

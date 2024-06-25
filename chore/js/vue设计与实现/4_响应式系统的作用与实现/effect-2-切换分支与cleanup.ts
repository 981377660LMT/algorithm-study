// https://www.cnblogs.com/wenruo/p/17050995.html

interface IEffectFn {
  (): void
  deps: Set<() => void>[]
}

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: IEffectFn | undefined

/** 用于注册副作用函数. */
function effect(fn: () => void) {
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    fn()
    activeEffect = undefined
  }

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
const data = { name: 'hello world', ok: true }

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

  /**
   * 不能直接执行effects，因为执行 effects 会先 cleanup 清除 bucket 中的依赖集合.
   * 但是再次执行后会再在集合中 push 副作用函数，导致 forEach 死循环.
   */
  // depsSet.forEach(effect => effect())

  const effectsToRun = new Set(depsSet)
  effectsToRun.forEach(effect => effect())
}

export {}

if (require.main === module) {
  let text: any = ''

  effect(() => {
    // effect 在 reactiveData.ok 取不同值的时候，执行代码会发生变化，叫做分支切换
    // 很明显在 reactiveData.ok = true 时 effect 依赖 reactiveData.name，
    // 但是当 reactiveData.ok = false 的时候，就和 reactiveData.name 无关，
    // reactiveData.name 再修改时也不应该触发函数重新执行了.
    // !为了实现这个效果，代码需要修改，每次副作用函数重新执行的时候，我们要先把它从所有与之关联的依赖集合中删除。执行后会建立新的关联。
    text = reactiveData.ok ? reactiveData.name : 'not'
    console.log('run', { text })
  })

  // setTimeout(() => {
  reactiveData.ok = false
  reactiveData.name = 'changed'
  // }, 1000)
}

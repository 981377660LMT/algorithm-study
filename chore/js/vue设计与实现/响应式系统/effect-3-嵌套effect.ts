// https://www.cnblogs.com/wenruo/p/17050995.html
// 在 Vue 中如果我们使用组件嵌套组件，就会有 effect 嵌套执行。
//
// !如果有嵌套的 effect 执行，我们就需要在保存当前 effect 函数的同时，记录之前的 effect 函数，并在当前的函数之前完之后，
// 把上一层的 effect 赋值为 activeEffect。很简单的会想到用栈来实现这个功能

interface IEffectFn {
  (): void
  deps: Set<() => void>[]
}

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: IEffectFn | undefined

const effectStack: IEffectFn[] = []

/** 用于注册副作用函数. */
function effect(fn: () => void) {
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    effectStack.push(effectFn)
    fn()
    effectStack.pop()
    activeEffect = effectStack.length > 0 ? effectStack[effectStack.length - 1] : undefined
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
    text = reactiveData.ok ? reactiveData.name : 'not'

    effect(() => {
      console.log('inner', { text })
    })
  })

  reactiveData.ok = false
  reactiveData.name = 'changed'
}

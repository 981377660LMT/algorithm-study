// https://www.cnblogs.com/wenruo/p/17050995.html

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: (() => void) | undefined

/** 用于注册副作用函数. */
function effect(fn: () => void) {
  activeEffect = fn
  fn()
  activeEffect = undefined
}

const bucket = new WeakMap<object, Map<PropertyKey, Set<() => void>>>() // target -> target key -> Set<副作用函数>
const data: Record<PropertyKey, any> = { text: 'hello world', text1: 'before' }

/** reactive. */
const reactiveData = new Proxy(data, {
  get(target, key) {
    track(target, key)
    return target[key]
  },
  set(target, key, newValue) {
    target[key] = newValue
    trigger(target, key)
    return true
  }
})

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
}

function trigger(target: object, key: PropertyKey): void {
  let depsMap = bucket.get(target)
  if (depsMap) {
    let depsSet = depsMap.get(key)
    if (depsSet) {
      depsSet.forEach(effect => effect())
    }
  }
}

export {}

if (require.main === module) {
  let text = 'foo'

  // 触发读取.
  effect(() => {
    text = reactiveData.text
  })

  setTimeout(() => {
    reactiveData.text = 'bar'
    console.log({ text })
  }, 1000)
}

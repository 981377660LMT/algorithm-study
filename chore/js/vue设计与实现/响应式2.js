// 全局变量，用于存储被注册的副作用函数
let activeEffect

// effect 函数用于`注册`副作用函数
function effect(fn) {
  activeEffect = fn
  fn()
  activeEffect = undefined // 不过书上没有这一句
}
const bucket = new WeakMap()

const data = { text: 'hello world', text1: 'before' }
const obj = new Proxy(data, {
  get(target, key) {
    // 将副作用函数activeEffect加到对应的桶中
    track(target, key)
    return target[key]
  },
  set(target, key, newVal) {
    target[key] = newVal
    trigger(target, key)
    // 返回true代表设置操作成功
    return true
  }
})

function track(target, key) {
  if (!activeEffect) return
  let depsMap = bucket.get(target)
  if (!depsMap) {
    bucket.set(target, (depsMap = new Map()))
  }
  let deps = depsMap.get(key)
  if (!deps) {
    depsMap.set(key, (deps = new Set()))
  }
  deps.add(activeEffect)
}

function trigger(target, key) {
  let depsMap = bucket.get(target)
  if (depsMap) {
    let effects = depsMap.get(key)
    effects && effects.forEach(fn => fn())
  }
}

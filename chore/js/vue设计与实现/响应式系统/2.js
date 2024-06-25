let activeEffect

function effect(fn) {
  const effectFn = () => {
    // 执行前先清除依赖
    cleanup(effectFn)
    activeEffect = effectFn
    fn()
    activeEffect = undefined
  }
  // 用来存储与该副作用函数相关联的依赖集合
  effectFn.deps = []
  effectFn()
}
function cleanup(effectFn) {
  // 很简单 就是在每个依赖集合中把该函数删除
  for (let i = 0; i < effectFn.deps.length; i++) {
    const deps = effectFn.deps[i]
    deps.delete(effectFn)
  }
  effectFn.deps.length = 0
}
const bucket = new WeakMap()
// 在 track 中记录 deps
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
  // 当前副作用函数也记录下关联的依赖
  activeEffect.deps.push(deps)
}

function trigger(target, key) {
  let depsMap = bucket.get(target)
  if (depsMap) {
    let effects = depsMap.get(key)
    // 不能直接执行effects 因为执行 effects 会先 cleanup 清除 bucket 中的依赖集合
    // 但是再次执行后会再在集合中添加副作用函数
    // 这样会导致死循环
    const effectsToRun = new Set(effects)
    effectsToRun.forEach(fn => fn())
  }
}

const data = { ok: true, text: 'hello world' }
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

let innerText = ''
effect(() => {
  innerText = obj.ok ? obj.text : 'not'
  console.log('run!')
  console.log({ innerText })
})
setTimeout(() => {
  obj.ok = false
  obj.text = 'changed'
  obj.text = 'changed'
  obj.text = 'changed'
}, 1000)

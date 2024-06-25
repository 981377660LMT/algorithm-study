// 全局变量，用于存储被注册的副作用函数
let activeEffect

// effect 函数用于注册副作用函数
function effect(fn) {
  activeEffect = fn
  fn()
}

const data = { text: 'hello world' }
const bucket = new Set()
const obj = new Proxy(data, {
  get(target, key) {
    if (activeEffect) {
      bucket.add(activeEffect)
    }
    return target[key]
  },
  set(target, key, newVal) {
    target[key] = newVal
    bucket.forEach(fn => fn())
    // 返回true代表设置操作成功
    return true
  }
})

// 触发读取
effect(
  // 一个匿名的副作用函数
  () => {
    document.body.innerText = obj.text
  }
)

// 修改响应式数据
setTimeout(() => {
  obj.text = 'hello vue3'
}, 1000)

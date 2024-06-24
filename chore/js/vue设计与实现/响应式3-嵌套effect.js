let activeEffect
const effectStack = []
function effect(fn) {
  const effectFn = () => {
    // 执行前先清除依赖
    cleanup(effectFn)
    // 执行前先压入栈中
    activeEffect = effectFn
    effectStack.push(effect)
    fn()
    // 执行后弹出
    effectStack.pop()
    activeEffect = effectStack[effectStack.length - 1]
  }
  // 用来存储与该副作用函数相关联的依赖集合
  effectFn.deps = []
  effectFn()
}

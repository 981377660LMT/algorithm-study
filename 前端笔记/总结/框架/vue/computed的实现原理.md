1. 当组件初始化的时候，computed 和 data 会分别建立各自的响应系统，Observer 遍历 data 中每个属性设置 get/set 数据拦截
2. 初始化 computed 会调用 initComputed 函数

   1. 注册一个 watcher 实例，并在内实例化一个 Dep 消息订阅器用作后续收集依赖（比如渲染函数的 watcher 或者其他观察该计算属性变化的 watcher ）
   2. 调用计算属性时会触发其 Object.defineProperty 的 get 访问器函数
   3. 调用 watcher.depend() 方法向自身的消息订阅器 dep 的 subs 中添加其他属性的 watcher
   4. 调用 watcher 的 evaluate 方法（进而调用 watcher 的 get 方法）让自身成为其他 watcher 的消息订阅器的订阅者，首先将 watcher 赋给 Dep.target，然后执行 getter 求值函数，当访问求值函数里面的属性（比如来自 data、props 或其他 computed）时，会同样触发它们的 get 访问器函数从而将该计算属性的 watcher 添加到求值函数中属性的 watcher 的消息订阅器 dep 中，当这些操作完成，最后关闭 Dep.target 赋为 null 并返回求值函数结果。

3. 当某个属性发生变化，触发 set 拦截函数，然后调用自身消息订阅器 dep 的 notify 方法，遍历当前 dep 中保存着所有订阅者 wathcer 的 subs 数组，并逐个调用 watcher 的 update 方法，完成响应更新。

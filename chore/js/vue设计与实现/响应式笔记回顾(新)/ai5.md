Vue.js 的响应式系统是其核心特性之一，通过数据劫持和依赖追踪实现数据变化到视图更新的自动同步。以下结合源码和设计原理，分 Vue 2 与 Vue 3 两代框架的实现方式，深入解析其核心机制：

---

### 一、Vue 2 响应式原理

Vue 2 基于 `Object.defineProperty` 实现数据劫持，通过 **Observer、Dep、Watcher** 三个核心类管理依赖关系。

#### 1. **数据劫持与依赖收集**

- **递归遍历对象属性**：初始化时，`Observer` 类递归遍历对象的每个属性，通过 `defineReactive` 方法将其转换为 `getter/setter`。若属性值是对象或数组，则继续递归监听。
  ```javascript
  function defineReactive(obj, key, val) {
    const dep = new Dep()
    Object.defineProperty(obj, key, {
      get() {
        if (Dep.target) {
          // 当前活跃的 Watcher
          dep.depend() // 收集依赖到 Dep
        }
        return val
      },
      set(newVal) {
        if (newVal === val) return
        val = newVal
        dep.notify() // 触发依赖更新
      }
    })
  }
  ```
- **数组的特殊处理**：通过重写数组的 `push`、`pop` 等 7 种方法，在调用这些方法时手动触发更新（因 `Object.defineProperty` 无法监听数组索引变化）。

#### 2. **依赖管理机制**

- **Dep（依赖管理器）**：每个属性对应一个 `Dep` 实例，用于存储依赖该属性的 `Watcher`。当属性被访问时（触发 `getter`），当前 `Watcher` 会被添加到 `Dep.subs` 中。
- **Watcher（观察者）**：作为连接数据与视图的桥梁，每个组件实例对应一个 `Watcher`。在组件渲染时，`Watcher` 会订阅所有被访问的数据属性，形成依赖关系。数据变化时，`Dep` 通知所有关联的 `Watcher` 执行更新。

#### 3. **局限性**

- **对象新增/删除属性无法监听**：需通过 `Vue.set` 或 `Vue.delete` 手动触发更新。
- **数组监听需特殊处理**：直接通过索引修改数组元素或修改 `length` 属性不会触发更新。

---

### 二、Vue 3 响应式原理

Vue 3 改用 `Proxy` 和 `Reflect` 重构响应式系统，解决了 Vue 2 的局限性，并优化了性能与功能扩展性。

#### 1. **Proxy 代理与反射**

- **Proxy 拦截操作**：通过 `Proxy` 代理目标对象，拦截 `get`、`set`、`deleteProperty` 等 13 种操作，无需递归初始化所有属性，按需触发依赖收集。
  ```javascript
  const objProxy = new Proxy(obj, {
    get(target, key, receiver) {
      track(target, key) // 依赖收集
      return Reflect.get(target, key, receiver)
    },
    set(target, key, value, receiver) {
      const result = Reflect.set(target, key, value, receiver)
      trigger(target, key) // 触发更新
      return result
    }
  })
  ```
- **Reflect 保证行为一致性**：使用 `Reflect` 方法操作目标对象，确保 `this` 指向正确，避免因 `Proxy` 代理导致的副作用。

#### 2. **依赖追踪优化**

- **WeakMap 结构管理依赖**：采用 `WeakMap(target → Map(key → Set(effects))` 三层结构存储依赖关系，避免内存泄漏，同时支持动态属性。
- **副作用函数（Effect）**：通过 `effect` API 注册副作用函数，执行时触发 `track` 收集依赖。函数内部访问的每个属性均会建立精确的依赖关系，避免无效更新。

#### 3. **优势与扩展**

- **动态属性监听**：支持对象新增/删除属性、数组索引修改等场景的自动响应。
- **惰性依赖收集**：仅在副作用函数实际访问属性时收集依赖，减少初始化开销。

---

### 三、核心设计思想对比

| **特性**         | **Vue 2（Object.defineProperty）**      | **Vue 3（Proxy）**                |
| ---------------- | --------------------------------------- | --------------------------------- |
| **初始化性能**   | 递归遍历所有属性，性能较差              | 按需代理，惰性收集依赖，性能更优  |
| **动态属性支持** | 需手动调用 `Vue.set`/`Vue.delete`       | 自动支持新增/删除属性             |
| **数组监听**     | 需重写数组方法，无法监听索引和 `length` | 直接监听数组索引和 `length` 变化  |
| **嵌套对象处理** | 递归初始化所有子属性                    | 按需代理，仅在访问时初始化        |
| **内存管理**     | 依赖关系通过闭包存储，可能泄漏          | 使用 `WeakMap` 自动释放无引用依赖 |

---

### 四、源码实现关键点

1. **Vue 2 的 `Observer` 类**：

   - 递归调用 `defineReactive` 转换属性为响应式。
   - 数组通过重写原型方法触发更新（如 `arrayMethods.push`）。

2. **Vue 3 的 `reactive` 函数**：
   - 使用 `Proxy` 创建代理对象，结合 `track` 和 `trigger` 函数管理依赖。
   - `effect` 函数包裹副作用逻辑，通过全局变量 `activeEffect` 标记当前执行的副作用。

---

### 总结

Vue 的响应式系统通过数据劫持与发布-订阅模式实现，Vue 2 受限于 `Object.defineProperty` 的设计，需手动处理动态属性和数组；Vue 3 利用 `Proxy` 的灵活性，实现了更高效、更全面的响应式机制。理解这一原理有助于优化代码结构，避免响应式失效问题（如直接修改数组索引），并更好地使用 `reactive`、`ref` 等 API。

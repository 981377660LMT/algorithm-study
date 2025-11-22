**沙箱 (Sandbox)** 的核心目标是：**在一个受控的环境中运行不可信的代码，防止其污染全局作用域、窃取敏感数据或破坏宿主环境。**

在前端/JS 领域，沙箱主要用于：微前端（Micro-frontends）、在线代码编辑器（CodePen）、插件系统（VS Code Extensions）等。

我们将实现一套基于 **Proxy** 的现代沙箱机制（类似微前端框架 `qiankun` 的实现原理）。

### 核心特性

1.  **全局变量隔离**：沙箱内的 `window` 修改不影响真实的 `window`。
2.  **快照与恢复**：沙箱激活时加载状态，失活时还原环境。
3.  **白名单机制**：控制沙箱能访问哪些全局 API（如 `localStorage`, `fetch`）。
4.  **多实例共存**：支持同时运行多个沙箱。

---

### 1. 核心接口定义

```typescript
export interface Sandbox {
  name: string
  proxy: Window // 沙箱暴露给代码的 window 对象
  active(): void // 激活沙箱
  inactive(): void // 失活沙箱
}
```

### 2. 基于 Proxy 的沙箱 (ProxySandbox)

这是现代浏览器中最主流的沙箱实现方案。它不需要对全局对象进行 Diff（快照），而是直接代理对 `window` 的访问。

```typescript
// 模拟一个假的 Window 类型
type FakeWindow = Window & Record<string, any>

export class ProxySandbox implements Sandbox {
  public name: string
  public proxy: Window
  public running = false

  // 沙箱内新增/修改的属性，存储在这里
  private updatedValueSet = new Map<string | symbol, any>()

  // 真实的 window 对象
  private globalContext: Window

  constructor(name: string, globalContext = window) {
    this.name = name
    this.globalContext = globalContext

    const { updatedValueSet } = this

    // 创建 Proxy
    // 核心原理：
    // 1. 读属性：先从 updatedValueSet 找，找不到去 globalContext 找。
    // 2. 写属性：直接写到 updatedValueSet，不污染 globalContext。
    this.proxy = new Proxy(this.globalContext, {
      // --- 拦截读取 (Get) ---
      get: (target: Window, p: string | symbol): any => {
        // 1. 如果沙箱处于非激活状态，理论上不应该允许访问（视策略而定）
        // 这里为了简单，允许访问，但通常应该报错或返回 undefined

        // 2. 优先返回沙箱内修改过的值
        if (updatedValueSet.has(p)) {
          return updatedValueSet.get(p)
        }

        // 3. 兜底：返回真实 window 上的值
        const value = (target as any)[p]

        // 4. 关键点：处理函数绑定的 this 问题
        // 比如 console.log, alert, fetch，它们的 this 必须指向 window
        // 如果直接返回 value，this 可能会指向 proxy，导致 Illegal invocation 错误
        if (typeof value === 'function' && isNativeGlobalFunction(p)) {
          return value.bind(target)
        }

        return value
      },

      // --- 拦截设置 (Set) ---
      set: (target: Window, p: string | symbol, value: any): boolean => {
        if (!this.running) {
          console.warn(
            `[Sandbox: ${this.name}] Cannot set property '${String(p)}' when sandbox is inactive.`
          )
          return false
        }

        // 所有的修改都记录在 Map 中，绝不触碰真实 window
        updatedValueSet.set(p, value)
        return true
      },

      // --- 拦截属性判断 (Has) ---
      // 用于拦截 'prop' in window
      has: (target: Window, p: string | symbol): boolean => {
        return updatedValueSet.has(p) || p in target
      },

      // --- 拦截属性描述符 (GetOwnPropertyDescriptor) ---
      // 用于 Object.getOwnPropertyDescriptor(window, 'prop')
      getOwnPropertyDescriptor: (target: Window, p: string | symbol) => {
        if (updatedValueSet.has(p)) {
          return {
            value: updatedValueSet.get(p),
            writable: true,
            enumerable: true,
            configurable: true
          }
        }
        return Object.getOwnPropertyDescriptor(target, p)
      }
    })
  }

  active() {
    this.running = true
    console.log(`[Sandbox: ${this.name}] Activated`)
  }

  inactive() {
    this.running = false
    console.log(`[Sandbox: ${this.name}] Deactivated`)
    // 注意：ProxySandbox 不需要还原 window，因为我们根本没改过 window
    // 它的状态一直保存在 updatedValueSet 里，下次 active 时还在
  }
}

// 辅助函数：判断是否是原生全局函数（需要绑定 this）
function isNativeGlobalFunction(key: string | symbol): boolean {
  const nativeFunctions = ['alert', 'setTimeout', 'setInterval', 'fetch', 'console']
  return nativeFunctions.includes(String(key))
}
```

### 3. 基于快照的沙箱 (SnapshotSandbox)

为了兼容不支持 Proxy 的老浏览器（IE），我们需要一种降级方案。
原理：激活时把当前 window 拍照，失活时把 window 还原回照片的样子。

```typescript
export class SnapshotSandbox implements Sandbox {
  public name: string
  public proxy: Window // 在快照模式下，proxy 就是 window 本身
  public running = false

  private windowSnapshot: Record<string, any> = {}
  private modifyPropsMap: Record<string, any> = {} // 记录沙箱期间修改了哪些属性

  constructor(name: string) {
    this.name = name
    this.proxy = window
  }

  active() {
    this.windowSnapshot = {}

    // 1. 拍照：记录当前 window 的状态
    for (const prop in window) {
      if (window.hasOwnProperty(prop)) {
        this.windowSnapshot[prop] = (window as any)[prop]
      }
    }

    // 2. 恢复：把上次沙箱运行期间修改的属性，重新应用到 window 上
    Object.keys(this.modifyPropsMap).forEach(prop => {
      ;(window as any)[prop] = this.modifyPropsMap[prop]
    })

    this.running = true
  }

  inactive() {
    this.modifyPropsMap = {}

    // 1. 记录差异：遍历当前 window，看看哪些属性变了
    for (const prop in window) {
      if (window.hasOwnProperty(prop)) {
        if ((window as any)[prop] !== this.windowSnapshot[prop]) {
          // 记录下来，下次激活时要用
          this.modifyPropsMap[prop] = (window as any)[prop]

          // 2. 还原：把 window 变回拍照时的样子
          ;(window as any)[prop] = this.windowSnapshot[prop]
        }
      }
    }

    this.running = false
  }
}
```

### 4. 沙箱执行器 (Executor)

有了沙箱环境（Proxy），我们还需要一个机制来执行代码。通常使用 `new Function` 或 `eval`，配合 `with` 语句来限制作用域。

```typescript
export function evaluateCode(code: string, sandbox: Sandbox) {
  // 核心黑魔法：with (sandbox)
  // 在 with 块内部，访问变量会优先去 sandbox.proxy 上找
  // 如果找不到，才会去上一层作用域（但我们通过 Proxy 拦截了，通常不会漏出去）

  // 包装代码：
  // (function(window, self, globalThis) {
  //    with(window) {
  //       ... 用户代码 ...
  //    }
  // }).call(sandbox.proxy, sandbox.proxy, sandbox.proxy, sandbox.proxy)

  const wrappedCode = `
    ;(function(window, self, globalThis){
      with(window) {
        ${code}
      }
    }).call(window, window, window, window);
  `

  try {
    // 使用 new Function 执行
    // 这里的 window 参数实际上是我们传进去的 sandbox.proxy
    const fn = new Function('window', wrappedCode)
    fn(sandbox.proxy)
  } catch (e) {
    console.error(`[Sandbox: ${sandbox.name}] Execution Error:`, e)
  }
}
```

### 5. 实战演示

让我们看看这两个沙箱是如何工作的，以及它们是否真的隔离了环境。

```typescript
// --- 场景演示 ---

// 1. 创建两个沙箱
const sandboxA = new ProxySandbox('App-A')
const sandboxB = new ProxySandbox('App-B')

// 2. 激活沙箱 A
sandboxA.active()

// 3. 在沙箱 A 中执行代码
const codeA = `
  window.appName = 'I am App A';
  window.globalVar = 100;
  console.log('In App A:', window.appName);
`
evaluateCode(codeA, sandboxA)

// 4. 验证隔离性
console.log('Real Window appName:', (window as any).appName) // 应该是 undefined
console.log('Sandbox A appName:', (sandboxA.proxy as any).appName) // 'I am App A'

// 5. 切换到沙箱 B
sandboxA.inactive() // A 失活
sandboxB.active() // B 激活

// 6. 在沙箱 B 中执行代码
const codeB = `
  window.appName = 'I am App B';
  console.log('In App B:', window.appName);
  console.log('Can I see A\'s var?', window.globalVar); // 应该是 undefined
`
evaluateCode(codeB, sandboxB)

// 7. 再次切回 A
sandboxB.inactive()
sandboxA.active()

console.log('Back to A, check globalVar:', (sandboxA.proxy as any).globalVar) // 应该是 100 (状态保留)

/**
 * 预期输出:
 * [Sandbox: App-A] Activated
 * In App A: I am App A
 * Real Window appName: undefined  <-- 成功隔离！
 * Sandbox A appName: I am App A
 * [Sandbox: App-A] Deactivated
 * [Sandbox: App-B] Activated
 * In App B: I am App B
 * Can I see A's var? undefined    <-- 成功隔离！
 * [Sandbox: App-B] Deactivated
 * [Sandbox: App-A] Activated
 * Back to A, check globalVar: 100 <-- 状态恢复！
 */
```

### 设计要点解析

1.  **`with` 语句的作用**：
    这是沙箱实现的灵魂。虽然 `with` 在严格模式下被禁用，但在沙箱实现中它是不可替代的。

    ```javascript
    with (proxy) {
      appName = 'A' // 等价于 proxy.appName = 'A'
    }
    ```

    如果没有 `with`，用户代码里的 `appName = 'A'` 会直接变成全局变量赋值，Proxy 拦截不到。

2.  **`isNativeGlobalFunction` 的处理**：
    这是一个常见的坑。如果你在沙箱里调用 `alert('hi')`，Proxy 会拦截 `get('alert')`。
    如果你直接返回 `window.alert`，那么执行时 `this` 指向的是 `Proxy` 对象。
    浏览器原生的 `alert`、`fetch` 等方法会检查 `this` 是否真的是 `window`，如果不是就会报错 `TypeError: Illegal invocation`。
    所以必须用 `.bind(window)` 修正 `this`。

3.  **Proxy vs Snapshot**：

    - **ProxySandbox**：支持多实例（因为每个 Proxy 都是独立的对象），性能好（只拦截修改），是微前端的首选。
    - **SnapshotSandbox**：不支持多实例（因为同一时间只能有一个真实的 window 被修改），性能差（每次切换都要遍历 window 属性），仅作为兼容方案。

4.  **逃逸风险**：
    这套机制能防住大部分无意的污染，但防不住恶意的攻击。
    例如：可以通过 `({}).constructor.prototype` 修改原生对象的原型，这会影响到沙箱外部。
    要解决这个问题，需要更高级的 **ShadowRealm** (TC39 提案) 或 **iframe** 隔离。

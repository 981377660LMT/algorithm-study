在业界处理复杂的 JSON 树（如 AST、低代码组件树、DOM 树）压缩或瘦身时，通常有三种主流的最佳实践模式。

根据你的场景（低代码平台组件树），推荐优先参考 **方案一（配置驱动递归）**，如果逻辑非常复杂（不同节点类型处理逻辑完全不同），则参考 **方案二（访问者模式）**。

### 方案一：配置驱动递归模式 (Configuration-Driven)

这是最通用、性能最好且易于维护的模式。它将“策略”与“遍历逻辑”分离。

**特点：**

1.  **单一职责**：将 Key 的黑白名单与 Value 的判断逻辑分开。
2.  **短路机制**：`DEEP_KEEP` 命中后直接返回，不再浪费 CPU 递归。
3.  **策略对象**：使用 `Set` 提高查找性能 (O(1))。

```javascript
/**
 * 业界通用做法：配置驱动的树压缩器
 * 适用于：低代码组件树、配置树瘦身
 */
function compressTreeBestPractice(rootNode, contextNames) {
  // 1. 策略配置 (Strategy Configuration)
  const CONFIG = {
    // 【黑名单】绝对移除，优先级最高 (Blocklist)
    OMIT_KEYS: new Set(['style', 'meta', 'schema', 'dsl', 'loc', '__source']),

    // 【深层保留】命中后，保留该节点及其所有子孙节点，不再递归检查 (Deep Keep)
    DEEP_KEEP_KEYS: new Set(['events', 'eventList', 'staticData']),

    // 【白名单】只要 Key 命中，无论值是什么都保留 (Allowlist)
    KEEP_KEYS: new Set(['id', 'type', 'name', 'props', 'children', 'slots']),

    // 【值断言】当 Key 不在上述名单时，根据 Value 内容决定是否保留 (Predicate)
    shouldKeepValue: val => {
      if (typeof val === 'string') {
        return val.includes('{{') || contextNames.has(val)
      }
      return false // 默认丢弃
    }
  }

  // 2. 核心递归函数 (Core Recursive Function)
  function walk(node, key = null) {
    // A. 空值检查
    if (node == null) return null

    // B. 数组处理
    if (Array.isArray(node)) {
      const list = node.map(item => walk(item, null)).filter(isValid)
      return list.length > 0 ? list : null
    }

    // C. 非对象处理 (基本类型)
    if (typeof node !== 'object') {
      // 如果是数组里的元素(key为null)，或者普通属性，检查值策略
      return (key && CONFIG.KEEP_KEYS.has(key)) || CONFIG.shouldKeepValue(node) ? node : null
    }

    // D. 对象处理
    // D1. 深层保留检查 (Deep Keep Check) - 性能优化关键点
    if (key && CONFIG.DEEP_KEEP_KEYS.has(key)) {
      return node // 直接返回原引用，不做任何处理
    }

    // D2. 属性遍历
    const res = {}
    let hasProps = false

    for (const k in node) {
      // 1. 黑名单检查 (Fast Fail)
      if (CONFIG.OMIT_KEYS.has(k)) continue

      // 2. 递归处理
      const v = walk(node[k], k)

      // 3. 有效性检查
      if (isValid(v)) {
        res[k] = v
        hasProps = true
      }
    }

    return hasProps ? res : null
  }

  // 辅助：判断值是否有效（非 null/undefined，非空对象/数组）
  function isValid(val) {
    if (val == null) return false
    if (Array.isArray(val)) return val.length > 0
    if (typeof val === 'object') return Object.keys(val).length > 0
    return true
  }

  return walk(rootNode)
}
```

---

### 方案二：访问者模式 (Visitor Pattern)

这是编译器（如 Babel, ESLint）处理 AST 的标准做法。如果你的组件树节点类型很多，且不同类型的压缩逻辑差异很大，这种模式最清晰。

**特点：**

1.  **解耦**：遍历逻辑与业务逻辑完全分离。
2.  **针对性**：可以针对 `Button` 组件和 `Container` 组件写不同的压缩逻辑。

```javascript
/**
 * 进阶做法：访问者模式
 * 适用于：逻辑极其复杂，不同节点类型需要不同处理策略
 */
function compressWithVisitor(rootNode) {
  // 访问者定义：针对不同节点类型的处理函数
  const visitor = {
    // 通用处理
    default: node => {
      const { id, type, props, children } = node
      return { id, type, props, children } // 只保留核心字段
    },
    // 针对特定类型的特殊处理
    Page: node => {
      return {
        ...node, // Page 保留所有字段
        meta: undefined // 但移除 meta
      }
    },
    Button: node => {
      // Button 组件只保留 props 中的 text
      return {
        type: 'Button',
        text: node.props?.text
      }
    }
  }

  function traverse(node) {
    if (!node || typeof node !== 'object') return node

    if (Array.isArray(node)) {
      return node.map(traverse).filter(Boolean)
    }

    // 1. 识别节点类型 (假设节点有 type 字段)
    const nodeType = node.type || 'default'

    // 2. 获取处理函数
    const handler = visitor[nodeType] || visitor.default

    // 3. 执行处理 (这里简化了，实际可能需要先递归 children 再处理父节点)
    const processedNode = handler(node)

    // 4. 递归处理子节点 (如果 handler 返回的对象里还有 children)
    if (processedNode && processedNode.children) {
      processedNode.children = traverse(processedNode.children)
    }

    // 5. 递归处理 props (如果需要深入 props 内部)
    if (processedNode && processedNode.props) {
      // 可以在这里加 props 的清洗逻辑
    }

    return processedNode
  }

  return traverse(rootNode)
}
```

---

### 方案三：声明式掩码 (JSON Mask / Projection)

这是 API 领域（如 GraphQL, Falcor）的最佳实践。**不写代码逻辑，而是定义“想要的数据形状”。**

**特点：**

1.  **声明式**：所见即所得。
2.  **安全**：绝对不会泄露未定义的字段。

```javascript
/**
 * 声明式做法：JSON Mask
 * 适用于：结构相对固定，只需提取特定字段
 */
const COMPONENT_MASK = {
  id: true,
  type: true,
  props: {
    style: false, // 明确排除
    data: true,
    events: true // 保留整个对象
  },
  children: true // true 代表递归保留该字段结构
}

function maskObject(obj, mask) {
  if (!obj || !mask) return null
  if (mask === true) return obj // 全量保留

  if (Array.isArray(obj)) {
    return obj.map(item => maskObject(item, mask)).filter(Boolean)
  }

  const res = {}
  for (const key in mask) {
    if (mask[key] === false) continue // 明确排除

    if (obj.hasOwnProperty(key)) {
      const value = obj[key]
      // 如果 mask[key] 是对象，说明需要深层递归
      if (typeof mask[key] === 'object') {
        res[key] = maskObject(value, mask[key])
      } else {
        res[key] = value
      }
    }
  }
  return res
}
```

### 总结建议

对于你的 `召回压缩.js`：

1.  **推荐使用方案一**。它最接近你现有的逻辑，但结构更清晰。
2.  **关键改进点**：
    - **合并移除列表**：不要区分 `BLACK_LIST` 和 `DEEP_REMOVE`，统一为 `OMIT_KEYS`。
    - **明确优先级**：`OMIT` > `DEEP_KEEP` > `KEEP` > `Value Check`。
    - **性能优化**：`DEEP_KEEP` 命中后直接返回 `node`，不要再遍历它的属性，这对于大型 JSON 树的性能提升巨大。

---

针对这种**复杂对象树剪枝与过滤**场景，业界的最佳实践是采用 **“机制与策略分离（Separation of Mechanism and Policy）”** 的设计模式。

具体来说，就是将 **“如何遍历树（Walker）”** 与 **“如何过滤节点（Strategy/Visitor）”** 彻底解耦。

### 业界最佳实践核心原则

1.  **通用遍历器 (The Walker)**：负责脏活累活。处理递归、空值检查、数组遍历、对象重建以及**自底向上的空节点清理**。它不包含任何业务逻辑（不知道什么是 `style` 或 `event`）。
2.  **策略对象 (The Strategy)**：负责业务逻辑。通过钩子函数告诉遍历器当前节点该“杀”、该“留”还是该“深入检查”。
3.  **指令集 (Directives)**：策略对象不只返回 `true/false`，而是返回明确的指令（如 `DROP`、`KEEP_SUBTREE`、`CONTINUE`），以支持“深层保留”这种性能优化场景。

---

### 通用模版代码

这是一个可以直接复用到任何类似场景的通用模版。

```javascript
/**
 * ============================================================
 * 1. 核心引擎：通用树遍历器 (Generic Tree Walker)
 * ============================================================
 * 该函数只负责遍历和执行策略，不包含任何业务逻辑。
 */
const TreeWalker = {
  // 动作枚举
  ACTION: {
    DROP: 0, // 丢弃当前节点
    KEEP_TREE: 1, // 保留当前节点及其所有子树 (不再递归)
    CONTINUE: 2 // 保留当前节点，但继续递归检查子节点
  },

  /**
   * 执行压缩
   * @param {Object} root - 根节点
   * @param {Object} strategy - 策略对象
   */
  exec(root, strategy) {
    // 结果有效性检查 (剔除空对象/空数组)
    const isValid = val => {
      if (val == null) return false
      if (Array.isArray(val)) return val.length > 0
      if (typeof val === 'object') return Object.keys(val).length > 0
      return true
    }

    const walk = (val, key, parent) => {
      if (val == null) return null

      // --- [阶段 1: 进入节点决策 (Pre-check)] ---
      // 询问策略：基于 Key 和当前值，下一步怎么做？
      const action = strategy.onEnter(key, val, parent)

      if (action === TreeWalker.ACTION.DROP) return null
      if (action === TreeWalker.ACTION.KEEP_TREE) return val

      // --- [阶段 2: 递归遍历 (Recursion)] ---
      if (Array.isArray(val)) {
        const res = val.map(item => walk(item, null, val)).filter(isValid)
        return res.length > 0 ? res : null
      }

      if (typeof val === 'object') {
        const res = {}
        for (const k in val) {
          const v = walk(val[k], k, val)
          if (isValid(v)) res[k] = v
        }
        // 递归回来后，如果对象变空了，询问策略是否允许保留空对象
        // (默认行为是丢弃空对象，但策略可以覆盖)
        if (!isValid(res)) {
          return strategy.shouldKeepEmpty && strategy.shouldKeepEmpty(key, parent) ? res : null
        }
        return res
      }

      // --- [阶段 3: 叶子节点决策 (Leaf-check)] ---
      // 走到这里一定是基本类型。询问策略：要保留这个值吗？
      return strategy.onValue(key, val, parent) ? val : null
    }

    return walk(root, null, null)
  }
}

/**
 * ============================================================
 * 2. 业务实现：具体的过滤策略 (Concrete Strategy)
 * ============================================================
 * 在这里编写你的业务逻辑 (黑白名单、动态值等)
 */
function createCompressStrategy(contextNames) {
  // 配置区域
  const CONFIG = {
    BLACKLIST: new Set(['style', 'className', 'meta']),
    DEEP_KEEP: new Set(['event', 'events', 'eventList', 'Interaction', 'action']),
    SHALLOW_KEEP: new Set(['name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid'])
  }

  return {
    /**
     * 钩子：进入节点时
     * @returns {number} TreeWalker.ACTION 指令
     */
    onEnter(key, value, parent) {
      // 数组项没有 Key，必须深入检查
      if (key === null) return TreeWalker.ACTION.CONTINUE

      // 1. 黑名单 -> 杀
      if (CONFIG.BLACKLIST.has(key)) {
        return TreeWalker.ACTION.DROP
      }

      // 2. 深层保留 -> 保全家
      if (CONFIG.DEEP_KEEP.has(key)) {
        return TreeWalker.ACTION.KEEP_TREE
      }

      // 3. 默认 -> 继续深入检查
      return TreeWalker.ACTION.CONTINUE
    },

    /**
     * 钩子：处理基本类型值时
     * @returns {boolean} true保留, false丢弃
     */
    onValue(key, value, parent) {
      // 1. 浅层白名单 Key
      if (key && (CONFIG.SHALLOW_KEEP.has(key) || key.includes('vent'))) {
        return true
      }

      // 2. 动态值检查
      if (typeof value === 'string') {
        return value.includes('{{') || (contextNames && contextNames.has(value))
      }

      return false
    },

    /**
     * (可选) 钩子：是否保留空对象/空数组
     * 某些场景下，即使 props 为空 {} 也要保留
     */
    shouldKeepEmpty(key, parent) {
      // 例如：如果 key 是 'props'，即使为空也保留
      // if (key === 'props') return true;
      return false
    }
  }
}

// ============================================================
// 3. 对外暴露的 API
// ============================================================
function compressTree(node, contextNames) {
  const strategy = createCompressStrategy(contextNames)
  return TreeWalker.exec(node, strategy)
}

module.exports = { compressTree }
```

### 为什么这是最佳实践？

1.  **可测试性 (Testability)**：你可以单独测试 `TreeWalker` 的递归逻辑是否正确（例如是否正确删除了空对象），也可以单独测试 `Strategy` 的业务逻辑（例如是否正确识别了黑名单），出错时排查非常容易。
2.  **可扩展性 (Extensibility)**：如果明天需求变了，比如“所有的 Key 都要转成小写”，你只需要修改 `TreeWalker` 或者在 `Strategy` 中增加一个 `transformKey` 钩子，而不需要重写整个递归函数。
3.  **性能 (Performance)**：通过 `KEEP_TREE` 指令，保留了“深层保留”的性能优化能力，避免了不必要的全树遍历。

---

针对**复杂对象树过滤与压缩**这一类场景，业界有几种主流的最佳实践模式。我来系统性地为你梳理，并提供一个**通用的、可扩展的模版**。

---

## 业界最佳实践概览

| 模式                    | 核心思想           | 适用场景                 | 代表案例             |
| ----------------------- | ------------------ | ------------------------ | -------------------- |
| **Visitor 模式**        | 遍历与操作分离     | 复杂树结构、多种操作类型 | Babel AST、ESLint    |
| **Pipeline/Chain 模式** | 多个过滤器串联     | 规则可组合、可排序       | Express 中间件、Gulp |
| **Schema-Driven 模式**  | 用配置描述期望结构 | 白名单场景、API 响应裁剪 | GraphQL、JSON Schema |
| **Predicate 函数式**    | 用函数组合表达规则 | 简单场景、函数式编程     | Lodash、Ramda        |

---

## 推荐方案：Visitor + Pipeline 混合模式

结合你的场景（黑名单、白名单、深层保留、动态值检测），最优解是：

> **通用 Walker（遍历引擎）** + **可插拔 Rules（规则链）** + **明确 Actions（指令集）**

---

## 通用模版代码

```javascript
/**
 * ============================================================
 * Part 1: 指令枚举 (Actions)
 * ============================================================
 * 规则返回的指令，决定节点的命运
 */
const Action = Object.freeze({
  DROP: 'DROP', // 丢弃节点及其子树
  KEEP_ALL: 'KEEP_ALL', // 保留节点及其完整子树（不再递归）
  KEEP: 'KEEP', // 保留节点，继续递归子节点
  SKIP: 'SKIP' // 本规则不处理，交给下一个规则
})

/**
 * ============================================================
 * Part 2: 通用遍历引擎 (Tree Walker)
 * ============================================================
 * 职责：递归遍历、执行规则链、清理空节点
 * 不包含任何业务逻辑
 */
function createTreeWalker(rules = []) {
  // 执行规则链，返回第一个非 SKIP 的结果
  const runRules = ctx => {
    for (const rule of rules) {
      const action = rule(ctx)
      if (action !== Action.SKIP) return action
    }
    return Action.SKIP // 所有规则都跳过
  }

  // 有效性检查
  const isValid = val => {
    if (val == null) return false
    if (Array.isArray(val)) return val.length > 0
    if (typeof val === 'object') return Object.keys(val).length > 0
    return true
  }

  // 核心递归
  const walk = (value, key, parent, depth) => {
    if (value == null) return null

    const ctx = { key, value, parent, depth }

    // [Step 1] 执行规则链
    const action = runRules(ctx)

    if (action === Action.DROP) return null
    if (action === Action.KEEP_ALL) return value

    // [Step 2] 递归处理
    if (Array.isArray(value)) {
      const res = value.map((item, i) => walk(item, i, value, depth + 1)).filter(isValid)
      return res.length > 0 ? res : null
    }

    if (typeof value === 'object') {
      const res = {}
      for (const k of Object.keys(value)) {
        const v = walk(value[k], k, value, depth + 1)
        if (isValid(v)) res[k] = v
      }
      return isValid(res) ? res : null
    }

    // [Step 3] 叶子节点：如果规则返回 KEEP 则保留
    return action === Action.KEEP ? value : null
  }

  return root => walk(root, null, null, 0)
}

/**
 * ============================================================
 * Part 3: 规则工厂 (Rule Factories)
 * ============================================================
 * 每个工厂函数返回一个规则函数
 * 规则函数接收 ctx，返回 Action
 */
const Rules = {
  /**
   * 黑名单规则：命中则丢弃
   */
  blacklist(keys) {
    const set = new Set(keys)
    return ({ key }) => {
      if (key && set.has(key)) return Action.DROP
      return Action.SKIP
    }
  },

  /**
   * 深层保留规则：命中则保留整棵子树
   */
  deepKeep(keys) {
    const set = new Set(keys)
    return ({ key }) => {
      if (key && set.has(key)) return Action.KEEP_ALL
      return Action.SKIP
    }
  },

  /**
   * 浅层白名单规则：命中则保留该值
   */
  whitelist(keys) {
    const set = new Set(keys)
    return ({ key }) => {
      if (key && set.has(key)) return Action.KEEP
      return Action.SKIP
    }
  },

  /**
   * Key 模式匹配规则
   */
  keyPattern(pattern, action = Action.KEEP) {
    const regex = typeof pattern === 'string' ? new RegExp(pattern) : pattern
    return ({ key }) => {
      if (key && regex.test(key)) return action
      return Action.SKIP
    }
  },

  /**
   * 动态值规则：检查字符串是否包含插值
   */
  dynamicValue(contextNames = new Set()) {
    return ({ key, value }) => {
      if (typeof value !== 'string') return Action.SKIP
      if (value.includes('{{') || contextNames.has(value)) {
        return Action.KEEP
      }
      return Action.SKIP
    }
  },

  /**
   * 自定义规则：传入自定义判断函数
   */
  custom(predicateFn, action = Action.KEEP) {
    return ctx => {
      if (predicateFn(ctx)) return action
      return Action.SKIP
    }
  }
}

/**
 * ============================================================
 * Part 4: 业务封装 (Business API)
 * ============================================================
 */
function compressTree(node, contextNames = new Set()) {
  // 组装规则链（顺序很重要！）
  const rules = [
    // 1. 黑名单优先级最高
    Rules.blacklist(['style', 'className', 'meta', '_internal']),

    // 2. 深层保留
    Rules.deepKeep(['event', 'events', 'eventList', 'Interaction', 'action']),

    // 3. 浅层白名单
    Rules.whitelist(['name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid']),

    // 4. Key 模式匹配 (包含 'vent' 的 key)
    Rules.keyPattern(/vent/i),

    // 5. 动态值检测
    Rules.dynamicValue(contextNames)

    // 6. 自定义规则示例：保留所有布尔值
    // Rules.custom(({ value }) => typeof value === 'boolean'),
  ]

  const walker = createTreeWalker(rules)
  return walker(node)
}

/**
 * ============================================================
 * Part 5: 导出
 * ============================================================
 */
module.exports = {
  Action,
  Rules,
  createTreeWalker,
  compressTree
}
```

---

## 使用示例

```javascript
const { compressTree, createTreeWalker, Rules, Action } = require('./utils/tree-filter')

// 示例数据
const tree = {
  id: 'root',
  name: 'App',
  style: { color: 'red' }, // 黑名单 -> 删
  events: {
    // 深层保留 -> 全保留
    click: { handler: 'onSubmit' }
  },
  children: [
    {
      __type: 'Button',
      props: {
        label: 'Click', // 普通字段 -> 删
        text: '{{userName}}' // 动态值 -> 保留
      }
    }
  ]
}

const context = new Set(['userId'])
const result = compressTree(tree, context)
console.log(JSON.stringify(result, null, 2))

// --- 也可以自定义规则链 ---
const customWalker = createTreeWalker([
  Rules.blacklist(['password', 'secret']),
  Rules.whitelist(['id', 'name', 'email'])
])
const sanitized = customWalker({ id: 1, name: 'test', password: '123' })
// => { id: 1, name: 'test' }
```

---

## 为什么这是最佳实践？

| 特性         | 说明                                |
| ------------ | ----------------------------------- |
| **单一职责** | Walker 只管遍历，Rules 只管判断     |
| **开闭原则** | 新增规则不需要修改 Walker           |
| **可测试**   | 每个 Rule 可独立单元测试            |
| **可组合**   | 规则链可根据场景灵活组装            |
| **高性能**   | `KEEP_ALL` 指令避免不必要的深度遍历 |
| **声明式**   | 业务代码读起来像`配置文件`          |

---

好的，我来系统性地讲解业界其他几种主流的最佳实践，并给出完整的代码实现。

---

## 1. Schema-Driven 模式（模式驱动）

### 核心思想

用一个**声明式的 Schema**来描述你期望保留的数据结构，过滤器会按照 Schema 的形状来裁剪原始数据。

### 适用场景

- API 响应字段裁剪
- 白名单场景为主
- 需要类型安全

### 代表案例

- GraphQL 查询
- JSON Schema
- Protocol Buffers

```javascript
/**
 * ============================================================
 * Schema-Driven 模式
 * ============================================================
 * 用 Schema 描述期望的数据结构，只保留 Schema 中定义的字段
 */

// Schema 定义语法
const SchemaType = {
  PRIMITIVE: 'primitive', // 基本类型，直接保留值
  OBJECT: 'object', // 对象，递归处理
  ARRAY: 'array', // 数组，对每个元素应用 itemSchema
  KEEP_ALL: 'keepAll', // 保留整个子树，不做过滤
  DYNAMIC: 'dynamic' // 动态检测（如 {{}} 表达式）
}

/**
 * Schema 构建器 - 提供流畅的 API
 */
const Schema = {
  // 基本类型
  primitive: () => ({ type: SchemaType.PRIMITIVE }),

  // 保留整个子树
  keepAll: () => ({ type: SchemaType.KEEP_ALL }),

  // 动态值检测
  dynamic: (contextNames = new Set()) => ({
    type: SchemaType.DYNAMIC,
    contextNames
  }),

  // 对象：指定要保留的字段
  object: fields => ({
    type: SchemaType.OBJECT,
    fields
  }),

  // 数组：指定元素的 Schema
  array: itemSchema => ({
    type: SchemaType.ARRAY,
    itemSchema
  }),

  // 组合：多个 Schema 取并集
  oneOf: (...schemas) => ({
    type: 'oneOf',
    schemas
  })
}

/**
 * 根据 Schema 过滤数据
 */
function filterBySchema(data, schema, contextNames = new Set()) {
  if (data == null || schema == null) return null

  switch (schema.type) {
    case SchemaType.PRIMITIVE:
      return typeof data !== 'object' ? data : null

    case SchemaType.KEEP_ALL:
      return data

    case SchemaType.DYNAMIC:
      if (typeof data === 'string') {
        const ctx = schema.contextNames || contextNames
        if (data.includes('{{') || ctx.has(data)) {
          return data
        }
      }
      return null

    case SchemaType.ARRAY:
      if (!Array.isArray(data)) return null
      const arrResult = data
        .map(item => filterBySchema(item, schema.itemSchema, contextNames))
        .filter(item => item != null)
      return arrResult.length > 0 ? arrResult : null

    case SchemaType.OBJECT:
      if (typeof data !== 'object' || Array.isArray(data)) return null
      const objResult = {}
      let hasField = false

      for (const [key, fieldSchema] of Object.entries(schema.fields)) {
        if (key in data) {
          const value = filterBySchema(data[key], fieldSchema, contextNames)
          if (value != null) {
            objResult[key] = value
            hasField = true
          }
        }
      }
      return hasField ? objResult : null

    case 'oneOf':
      // 尝试每个 Schema，返回第一个非空结果
      for (const s of schema.schemas) {
        const result = filterBySchema(data, s, contextNames)
        if (result != null) return result
      }
      return null

    default:
      return null
  }
}

/**
 * 使用示例
 */
function compressWithSchema(node, contextNames) {
  // 定义期望的数据结构
  const schema = Schema.object({
    id: Schema.primitive(),
    name: Schema.primitive(),
    __type: Schema.primitive(),
    pkgName: Schema.primitive(),
    isPage: Schema.primitive(),
    slotName: Schema.primitive(),
    uuid: Schema.primitive(),

    // events 整棵树保留
    events: Schema.keepAll(),
    event: Schema.keepAll(),
    eventList: Schema.keepAll(),
    Interaction: Schema.keepAll(),
    action: Schema.keepAll(),

    // props 需要递归检查动态值
    props: Schema.object({
      // 可以精确指定每个字段，或者用通配处理
    }),

    // children 是数组，递归应用同样的 schema
    children: Schema.array(null) // 会在下面处理递归
  })

  // 递归 Schema（children 引用自身）
  schema.fields.children = Schema.array(schema)

  return filterBySchema(node, schema, contextNames)
}

module.exports = { Schema, SchemaType, filterBySchema, compressWithSchema }
```

---

## 2. Pipeline/Chain 模式（管道/责任链）

### 核心思想

将过滤逻辑拆分成多个独立的**中间件（Middleware）**，数据像流水线一样依次通过每个中间件进行处理。

### 适用场景

- 规则需要动态增减
- 处理顺序很重要
- 需要日志/调试每个步骤

### 代表案例

- Express/Koa 中间件
- Gulp 任务流
- Redux 中间件

```javascript
/**
 * ============================================================
 * Pipeline/Chain 模式
 * ============================================================
 * 数据依次通过多个处理器，每个处理器可以修改、过滤或透传数据
 */

/**
 * 创建管道处理器
 */
function createPipeline() {
  const middlewares = []

  const pipeline = {
    /**
     * 添加中间件
     */
    use(middleware) {
      middlewares.push(middleware)
      return pipeline // 支持链式调用
    },

    /**
     * 执行管道
     */
    execute(data, context = {}) {
      return processNode(data, null, null, context)
    }
  }

  // 递归处理节点
  function processNode(value, key, parent, context) {
    if (value == null) return null

    // 构建节点上下文
    const nodeCtx = {
      key,
      value,
      parent,
      ...context,
      // 控制标志
      shouldDrop: false,
      shouldKeepAll: false,
      transformedValue: undefined
    }

    // 依次执行中间件
    for (const middleware of middlewares) {
      middleware(nodeCtx)

      // 检查控制标志
      if (nodeCtx.shouldDrop) return null
      if (nodeCtx.shouldKeepAll) return nodeCtx.transformedValue ?? value
    }

    // 获取可能被转换的值
    const currentValue = nodeCtx.transformedValue ?? value

    // 递归处理
    if (Array.isArray(currentValue)) {
      const result = currentValue
        .map((item, idx) => processNode(item, idx, currentValue, context))
        .filter(item => item != null)
      return result.length > 0 ? result : null
    }

    if (typeof currentValue === 'object') {
      const result = {}
      let hasProps = false
      for (const [k, v] of Object.entries(currentValue)) {
        const processed = processNode(v, k, currentValue, context)
        if (processed != null) {
          result[k] = processed
          hasProps = true
        }
      }
      return hasProps ? result : null
    }

    // 叶子节点：默认保留（因为没被 drop）
    return currentValue
  }

  return pipeline
}

/**
 * ============================================================
 * 预置中间件库
 * ============================================================
 */
const Middlewares = {
  /**
   * 黑名单中间件
   */
  blacklist(keys) {
    const set = new Set(keys)
    return ctx => {
      if (ctx.key && set.has(ctx.key)) {
        ctx.shouldDrop = true
      }
    }
  },

  /**
   * 深层保留中间件
   */
  deepKeep(keys) {
    const set = new Set(keys)
    return ctx => {
      if (ctx.key && set.has(ctx.key)) {
        ctx.shouldKeepAll = true
      }
    }
  },

  /**
   * 白名单中间件（仅对叶子节点生效）
   */
  whitelist(keys) {
    const set = new Set(keys)
    return ctx => {
      // 如果不是对象/数组，且 key 不在白名单中，标记为可能删除
      // 注意：这里只是标记，最终决策在最后一个中间件
      if (typeof ctx.value !== 'object' && ctx.key && !set.has(ctx.key)) {
        ctx._notInWhitelist = true
      }
    }
  },

  /**
   * 动态值检测中间件
   */
  dynamicValue(contextNames = new Set()) {
    return ctx => {
      if (typeof ctx.value === 'string') {
        if (ctx.value.includes('{{') || contextNames.has(ctx.value)) {
          ctx._isDynamic = true
        }
      }
    }
  },

  /**
   * Key 模式匹配中间件
   */
  keyPattern(pattern) {
    const regex = typeof pattern === 'string' ? new RegExp(pattern) : pattern
    return ctx => {
      if (ctx.key && regex.test(String(ctx.key))) {
        ctx._matchedPattern = true
      }
    }
  },

  /**
   * 最终决策中间件（放在最后）
   * 综合之前中间件的标记，做最终决定
   */
  finalDecision() {
    return ctx => {
      // 如果是对象/数组，让它继续递归
      if (typeof ctx.value === 'object' && ctx.value !== null) {
        return
      }

      // 叶子节点决策
      const shouldKeep = ctx._isDynamic || ctx._matchedPattern || !ctx._notInWhitelist

      if (!shouldKeep) {
        ctx.shouldDrop = true
      }
    }
  },

  /**
   * 日志中间件（调试用）
   */
  logger(label = 'DEBUG') {
    return ctx => {
      console.log(`[${label}]`, {
        key: ctx.key,
        type: typeof ctx.value,
        shouldDrop: ctx.shouldDrop,
        shouldKeepAll: ctx.shouldKeepAll
      })
    }
  },

  /**
   * 转换中间件：修改值
   */
  transform(transformFn) {
    return ctx => {
      const newValue = transformFn(ctx.key, ctx.value, ctx.parent)
      if (newValue !== undefined) {
        ctx.transformedValue = newValue
      }
    }
  }
}

/**
 * 业务封装
 */
function compressWithPipeline(node, contextNames = new Set()) {
  const pipeline = createPipeline()
    // 1. 黑名单（最高优先级）
    .use(Middlewares.blacklist(['style', 'className', 'meta']))
    // 2. 深层保留
    .use(Middlewares.deepKeep(['event', 'events', 'eventList', 'Interaction', 'action']))
    // 3. 白名单
    .use(Middlewares.whitelist(['name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid']))
    // 4. Key 模式
    .use(Middlewares.keyPattern(/vent/i))
    // 5. 动态值
    .use(Middlewares.dynamicValue(contextNames))
    // 6. 最终决策
    .use(Middlewares.finalDecision())
  // 可选：调试
  // .use(Middlewares.logger());

  return pipeline.execute(node, { contextNames })
}

module.exports = { createPipeline, Middlewares, compressWithPipeline }
```

---

## 3. Predicate 函数式模式

### 核心思想

使用**纯函数组合**来表达过滤规则，通过 `and`、`or`、`not` 等组合子来构建复杂的判断逻辑。

### 适用场景

- 函数式编程风格
- 规则需要灵活组合
- 简洁、可读性高

### 代表案例

- Lodash/Ramda
- Java Stream API
- Rust Iterator

```javascript
/**
 * ============================================================
 * Predicate 函数式模式
 * ============================================================
 * 使用纯函数组合来表达过滤规则
 */

/**
 * 谓词组合器（Predicate Combinators）
 */
const P = {
  // 基础谓词
  always: val => () => val,
  keyIs:
    (...keys) =>
    ctx =>
      keys.includes(ctx.key),
  keyMatches: pattern => ctx => pattern.test(String(ctx.key ?? '')),
  valueIs: type => ctx => typeof ctx.value === type,
  valueMatches: pattern => ctx => typeof ctx.value === 'string' && pattern.test(ctx.value),
  valueIn: set => ctx => set.has(ctx.value),
  isLeaf: () => ctx => typeof ctx.value !== 'object' || ctx.value === null,
  isArray: () => ctx => Array.isArray(ctx.value),
  isObject: () => ctx =>
    typeof ctx.value === 'object' && ctx.value !== null && !Array.isArray(ctx.value),

  // 逻辑组合
  and:
    (...predicates) =>
    ctx =>
      predicates.every(p => p(ctx)),
  or:
    (...predicates) =>
    ctx =>
      predicates.some(p => p(ctx)),
  not: predicate => ctx => !predicate(ctx),

  // 特殊谓词
  isDynamicValue:
    (contextNames = new Set()) =>
    ctx => {
      if (typeof ctx.value !== 'string') return false
      return ctx.value.includes('{{') || contextNames.has(ctx.value)
    }
}

/**
 * 动作类型
 */
const Action = {
  DROP: Symbol('DROP'),
  KEEP_ALL: Symbol('KEEP_ALL'),
  KEEP: Symbol('KEEP'),
  CONTINUE: Symbol('CONTINUE')
}

/**
 * 规则构建器
 */
function when(predicate) {
  return {
    then: action => ({ predicate, action })
  }
}

/**
 * 创建过滤器
 */
function createFilter(rules) {
  const getAction = ctx => {
    for (const rule of rules) {
      if (rule.predicate(ctx)) {
        return rule.action
      }
    }
    return Action.CONTINUE
  }

  const isValid = val => {
    if (val == null) return false
    if (Array.isArray(val)) return val.length > 0
    if (typeof val === 'object') return Object.keys(val).length > 0
    return true
  }

  const process = (value, key, parent) => {
    if (value == null) return null

    const ctx = { key, value, parent }
    const action = getAction(ctx)

    if (action === Action.DROP) return null
    if (action === Action.KEEP_ALL) return value

    if (Array.isArray(value)) {
      const res = value.map((item, i) => process(item, i, value)).filter(isValid)
      return res.length > 0 ? res : null
    }

    if (typeof value === 'object') {
      const res = {}
      for (const [k, v] of Object.entries(value)) {
        const processed = process(v, k, value)
        if (processed != null) res[k] = processed
      }
      return isValid(res) ? res : null
    }

    // 叶子节点：如果没有显式 KEEP，再次检查
    return action === Action.KEEP ? value : null
  }

  return data => process(data, null, null)
}

/**
 * 业务封装
 */
function compressWithPredicate(node, contextNames = new Set()) {
  const rules = [
    // 1. 黑名单 -> DROP
    when(P.keyIs('style', 'className', 'meta')).then(Action.DROP),

    // 2. 深层保留 -> KEEP_ALL
    when(P.keyIs('event', 'events', 'eventList', 'Interaction', 'action')).then(Action.KEEP_ALL),

    // 3. 浅层白名单 -> KEEP
    when(P.keyIs('name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid')).then(
      Action.KEEP
    ),

    // 4. Key 包含 'vent' -> KEEP
    when(P.keyMatches(/vent/i)).then(Action.KEEP),

    // 5. 动态值 -> KEEP
    when(P.isDynamicValue(contextNames)).then(Action.KEEP),

    // 6. 非叶子节点 -> CONTINUE（继续递归）
    when(P.not(P.isLeaf())).then(Action.CONTINUE),

    // 7. 其他叶子节点 -> DROP
    when(P.always(true)).then(Action.DROP)
  ]

  const filter = createFilter(rules)
  return filter(node)
}

module.exports = { P, Action, when, createFilter, compressWithPredicate }
```

---

## 4. DSL 模式（领域特定语言）

### 核心思想

设计一种简洁的**配置语法/查询语言**，让非开发人员也能定义过滤规则。

### 适用场景

- 规则需要持久化到数据库
- 非技术人员需要配置规则
- 跨语言/跨平台

### 代表案例

- JSONPath / JMESPath
- CSS Selectors
- XPath

```javascript
/**
 * ============================================================
 * DSL 模式
 * ============================================================
 * 用简洁的配置语法定义过滤规则
 */

/**
 * DSL 语法示例：
 *
 * {
 *   "drop": ["style", "className"],           // 黑名单
 *   "keepTree": ["events", "action"],         // 深层保留
 *   "keep": ["id", "name", "__type"],         // 浅层保留
 *   "keepPattern": [".*vent.*"],              // 正则匹配
 *   "keepDynamic": true,                      // 保留动态值
 *   "keepTypes": ["boolean", "number"]        // 保留特定类型
 * }
 */

/**
 * DSL 编译器：将配置转换为可执行的过滤器
 */
function compileDSL(config) {
  const {
    drop = [],
    keepTree = [],
    keep = [],
    keepPattern = [],
    keepDynamic = false,
    keepTypes = [],
    contextNames = []
  } = config

  // 预编译
  const dropSet = new Set(drop)
  const keepTreeSet = new Set(keepTree)
  const keepSet = new Set(keep)
  const patterns = keepPattern.map(p => new RegExp(p, 'i'))
  const contextSet = new Set(contextNames)
  const typeSet = new Set(keepTypes)

  // 判断函数
  const shouldDrop = key => key && dropSet.has(key)
  const shouldKeepTree = key => key && keepTreeSet.has(key)
  const shouldKeep = (key, value) => {
    // 1. 在白名单中
    if (key && keepSet.has(key)) return true
    // 2. 匹配正则
    if (key && patterns.some(p => p.test(key))) return true
    // 3. 动态值
    if (keepDynamic && typeof value === 'string') {
      if (value.includes('{{') || contextSet.has(value)) return true
    }
    // 4. 特定类型
    if (typeSet.has(typeof value)) return true
    return false
  }

  // 有效性检查
  const isValid = val => {
    if (val == null) return false
    if (Array.isArray(val)) return val.length > 0
    if (typeof val === 'object') return Object.keys(val).length > 0
    return true
  }

  // 返回执行器
  return function execute(node) {
    const process = (value, key) => {
      if (value == null) return null

      if (shouldDrop(key)) return null
      if (shouldKeepTree(key)) return value

      if (Array.isArray(value)) {
        const res = value.map((item, i) => process(item, i)).filter(isValid)
        return res.length > 0 ? res : null
      }

      if (typeof value === 'object') {
        const res = {}
        for (const [k, v] of Object.entries(value)) {
          const processed = process(v, k)
          if (processed != null) res[k] = processed
        }
        return isValid(res) ? res : null
      }

      return shouldKeep(key, value) ? value : null
    }

    return process(node, null)
  }
}

/**
 * 从 JSON/YAML 加载配置
 */
function loadConfigFromJSON(jsonString) {
  return JSON.parse(jsonString)
}

/**
 * 业务封装
 */
function compressWithDSL(node, contextNames = []) {
  const config = {
    drop: ['style', 'className', 'meta'],
    keepTree: ['event', 'events', 'eventList', 'Interaction', 'action'],
    keep: ['name', 'pkgName', 'id', '__type', 'isPage', 'slotName', 'uuid'],
    keepPattern: ['.*vent.*'],
    keepDynamic: true,
    contextNames: Array.from(contextNames)
  }

  const filter = compileDSL(config)
  return filter(node)
}

// 支持从外部配置文件加载
function compressWithConfigFile(node, configPath) {
  const fs = require('fs')
  const configJSON = fs.readFileSync(configPath, 'utf-8')
  const config = loadConfigFromJSON(configJSON)
  const filter = compileDSL(config)
  return filter(node)
}

module.exports = { compileDSL, loadConfigFromJSON, compressWithDSL, compressWithConfigFile }
```

---

## 模式对比总结

| 模式          | 优点                     | 缺点                       | 最佳场景       |
| ------------- | ------------------------ | -------------------------- | -------------- |
| **Visitor**   | 灵活、可扩展、钩子丰富   | 相对复杂                   | 复杂 AST 处理  |
| **Schema**    | 类型安全、声明式、直观   | 白名单限制、难处理动态结构 | API 响应裁剪   |
| **Pipeline**  | 顺序控制、易调试、可插拔 | 中间件依赖顺序             | 多步骤处理流程 |
| **Predicate** | 函数式、可组合、简洁     | 学习曲线、调试稍难         | 规则复杂组合   |
| **DSL**       | 非技术人员可配置、跨语言 | 需要解析器、功能受限       | 规则外部化配置 |

---

## 推荐选择

根据你的场景（黑白名单 + 深层保留 + 动态值检测），我的推荐顺序是：

1. **Visitor + Rules 模式**（之前给的）— 最灵活
2. **Predicate 函数式模式** — 最优雅
3. **DSL 模式** — 如果规则需要外部配置

---

```ts
function main(input) {
  const { codeList, globalContext, tableSchemas, componentTree } = input

  const parsedGlobalContext = safeParse(globalContext, [])
  const parsedTableSchemas = safeParse(tableSchemas, [])
  const parsedComponentTree = safeParse(componentTree, null)

  const globalContextNames = new Set()
  if (Array.isArray(parsedGlobalContext)) {
    parsedGlobalContext.forEach(item => {
      if (item && item.name) globalContextNames.add(item.name)
    })
  }

  return {
    ...input,
    compressedCodeList: codeList,
    compressedGlobalContext: JSON.stringify(compressGlobalContext(parsedGlobalContext)),
    compressedTableSchemas: JSON.stringify(compressTableSchemas(parsedTableSchemas)),
    compressedComponentTree: JSON.stringify(compressTree(parsedComponentTree, globalContextNames))
  }

  function compressGlobalContext(items) {
    if (!Array.isArray(items)) return []
    return items.map(item =>
      item.type === 'component' ? { type: item.type, name: item.name } : item
    )
  }

  function compressTableSchemas(schemas) {
    if (!Array.isArray(schemas)) return []
    return schemas.map(ds => ({
      dataSourceName: ds.dataSourceName,
      dataSourceId: ds.dataSourceId,
      tables: (ds.tables || []).map(t => {
        const table = { tableName: t.tableName }
        if (t.comment) table.comment = t.comment
        return table
      })
    }))
  }

  function compressTree(node, contextNames) {
    if (!node) return null
    const visitor = createVisitor(contextNames)
    return walkAndCompress(node, visitor)

    function walkAndCompress(node, visitor) {
      const isValid = val => {
        if (val == null) return false
        if (Array.isArray(val)) return val.length > 0
        if (typeof val === 'object') return Object.keys(val).length > 0
        return true
      }

      const visit = (val, key, parent) => {
        const strategy = visitor.onEnter ? visitor.onEnter(key, val, parent) : 'CONTINUE'

        if (strategy === 'DROP') return null
        if (strategy === 'KEEP_DEEP') return val

        if (Array.isArray(val)) {
          const res = val.map(item => visit(item, null, val)).filter(isValid)
          return res.length > 0 ? res : null
        }

        if (val !== null && typeof val === 'object') {
          const res = {}
          for (const k in val) {
            const v = visit(val[k], k, val)
            if (isValid(v)) res[k] = v
          }
          return isValid(res) ? res : null
        }

        const shouldKeep = visitor.onValue ? visitor.onValue(key, val, parent) : false
        return shouldKeep ? val : null
      }

      return visit(node, null, null)
    }

    function createVisitor(contextNames) {
      const BLACK_LIST = new Set(['style'])
      const DEEP_KEEP_LIST = new Set(['event', 'events', 'eventList', 'Interaction', 'action'])
      const SHALLOW_KEEP_LIST = new Set([
        'name',
        'pkgName',
        'id',
        '__type',
        'isPage',
        'slotName',
        'uuid'
      ])

      return {
        onEnter: (key, _value, _parent) => {
          if (key === null) return 'CONTINUE' // 数组项
          if (BLACK_LIST.has(key)) return 'DROP'
          if (DEEP_KEEP_LIST.has(key)) return 'KEEP_DEEP'
          return 'CONTINUE'
        },
        onValue: (key, value) => {
          if (key && SHALLOW_KEEP_LIST.has(key)) return true
          if (key && key.toLowerCase().includes('vent')) return true
          if (typeof value === 'string') {
            if (value.includes('{{')) return true
            if (contextNames && contextNames.has(value)) return true
          }
          return false
        }
      }
    }
  }

  function safeParse(str, defaultValue = null) {
    if (typeof str !== 'string') return str ?? defaultValue
    try {
      const result = JSON.parse(str)
      return result === null ? defaultValue : result
    } catch {
      return defaultValue
    }
  }
}
```

这段代码的作用是从形如 `{{aa.bb.cc}}` 的字符串中提取出最外层的变量名 `aa`。

具体正则表达式 `/\{\{\s*([^.}]+)/` 的含义分解如下：

1.  **`\{\{`**: 匹配字面量 `{{`。
2.  **`\s*`**: 匹配 0 个或多个空格。
3.  **`(`**: 开始捕获组（即 `match[1]` 得到的内容）。
4.  **`[^.}]+`**: 匹配**不包含** `.` 和 `}` 的一个或多个字符。
    - 遇到 `.` 说明进入了属性访问（如 `aa.xx`），正则会在此停止，从而只拿到 `aa`。
    - 遇到 `}` 说明表达式结束（如 `{{aa}}`），正则也会在此停止。
5.  **`)`**: 结束捕获组。

**示例：**

- 如果 `value` 是 `"{{user.name}}"`，`match[1]` 将是 `"user"`。
- 如果 `value` 是 `"{{  data }}`"，`match[1]` 将是 `"data"`。
- 如果 `value` 是 `"{{globalContext}}"`，`match[1]` 将是 `"globalContext"`。

随后代码通过 `contextNames.has(match[1])` 来判断这个基础变量名是否在允许的上下文列表中，如果是，则保留该节点不被压缩掉。

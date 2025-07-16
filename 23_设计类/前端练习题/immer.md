Immer 是一个非常优雅的不可变数据结构库，它通过 **Proxy** 和 **Copy-on-Write (COW)** 策略来实现高效的不可变更新。让我详细分析其核心实现原理：

## 1. Immer 的核心思想

Immer 允许你像操作可变数据一样编写代码，但实际产生不可变的结果：

```javascript
import { produce } from 'immer'

const baseState = {
  users: [
    { id: 1, name: 'Alice', todos: [] },
    { id: 2, name: 'Bob', todos: [] }
  ]
}

const nextState = produce(baseState, draft => {
  draft.users[0].name = 'Alice Updated' // 看起来是可变操作
  draft.users.push({ id: 3, name: 'Charlie', todos: [] })
})

console.log(baseState === nextState) // false
console.log(baseState.users === nextState.users) // false
console.log(baseState.users[1] === nextState.users[1]) // true (共享未修改的部分)
```

## 2. 核心实现机制

### 2.1 Proxy 拦截操作

Immer 的核心是为每个对象创建一个 Proxy，拦截所有的读写操作：

```javascript
// 简化版的 Immer 实现
function createProxy(base, parent, prop) {
  const state = {
    type: 'object',
    scope: getCurrentScope(),
    modified: false,
    assigned: false,
    parent,
    base,
    draft: null,
    copy: null,
    revoked: false
  }

  const proxy = new Proxy(state, {
    get(state, prop) {
      if (prop === DRAFT_STATE) return state

      // 如果已经被修改，从 copy 中读取
      if (state.copy) {
        return state.copy[prop]
      }

      // 如果未修改，从原始对象读取
      const value = state.base[prop]

      // 如果值是对象，需要为其创建代理
      if (isObject(value)) {
        return createProxy(value, state, prop)
      }

      return value
    },

    set(state, prop, value) {
      // 标记为已修改
      if (!state.modified) {
        state.modified = true
        state.copy = Array.isArray(state.base) ? state.base.slice() : Object.assign({}, state.base)
      }

      state.copy[prop] = value
      return true
    }
  })

  return proxy
}
```

### 2.2 Copy-on-Write 策略

只有在第一次写操作时才创建副本：

```javascript
function prepareCopy(state) {
  if (!state.copy) {
    state.copy = Array.isArray(state.base)
      ? state.base.slice()  // 数组浅拷贝
      : Object.assign({}, state.base)  // 对象浅拷贝
  }
}

// 在 set trap 中调用
set(state, prop, value) {
  if (!state.modified) {
    if (state.base[prop] === value) {
      return true  // 值相同，不需要修改
    }

    markChanged(state)  // 标记修改并创建副本
  }

  state.copy[prop] = value
  return true
}
```

### 2.3 父子关系追踪

当子对象被修改时，需要通知父对象：

```javascript
function markChanged(state) {
  if (!state.modified) {
    state.modified = true
    prepareCopy(state)

    // 递归标记父对象
    if (state.parent) {
      markChanged(state.parent)
    }
  }
}
```

## 3. 完整的生产流程

```javascript
// 简化版的 produce 函数
function produce(base, recipe) {
  const scope = {
    patches: [],
    inversePatches: [],
    canAutoFreeze: true,
    drafts: [],
    parent: null
  }

  // 设置当前作用域
  setCurrentScope(scope)

  try {
    // 创建根代理
    const proxy = createProxy(base, null, null)
    scope.drafts.push(proxy)

    // 执行用户的修改函数
    const result = recipe(proxy)

    // 如果有返回值，直接使用
    if (result !== undefined && result !== proxy) {
      return result
    }

    // 否则从代理中提取最终结果
    return finalize(proxy, [])
  } finally {
    // 清理作用域
    setCurrentScope(scope.parent)
  }
}
```

### 3.4 最终化过程

```javascript
function finalize(draft, path) {
  const state = draft[DRAFT_STATE]

  if (!state) {
    return draft // 不是代理对象
  }

  if (!state.modified) {
    return state.base // 未修改，返回原始对象
  }

  // 递归处理子对象
  const copy = state.copy
  Object.keys(copy).forEach(prop => {
    const value = copy[prop]
    if (isProxy(value)) {
      copy[prop] = finalize(value, [...path, prop])
    }
  })

  // 冻结对象（在生产环境中）
  if (shouldFreeze) {
    Object.freeze(copy)
  }

  return copy
}
```

## 4. 性能优化要点

### 4.1 结构共享

```javascript
// 未修改的部分会被共享
const state = {
  a: { value: 1 },
  b: { value: 2 }
}

const newState = produce(state, draft => {
  draft.a.value = 10 // 只修改 a
})

// newState.b === state.b  (共享未修改的部分)
```

### 4.2 惰性创建代理

```javascript
// 只有在访问时才创建子对象的代理
get(state, prop) {
  const value = getFromCurrentState(state, prop)

  if (isObject(value) && !isProxy(value)) {
    return createProxy(value, state, prop)  // 惰性创建
  }

  return value
}
```

## 5. 关键优势

1. **API 简洁**：可以像操作普通对象一样编写代码
2. **性能高效**：只在需要时创建副本，最大化结构共享
3. **类型安全**：TypeScript 支持良好
4. **内存友好**：避免不必要的深拷贝

这就是 Immer 实现不可变数据的核心机制 - 通过 Proxy 拦截操作，Copy-on-Write 策略最小化复制，以及智能的结构共享来实现高性能的不可变更新。

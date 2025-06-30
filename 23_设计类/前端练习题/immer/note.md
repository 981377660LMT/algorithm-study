# Immer API 详解

Immer 是一个小型库，让您能以更便捷的方式处理不可变状态。它基于**草稿(draft)**的概念，让您可以直接修改数据，而底层实现会确保数据的不可变性。本文详细介绍所有 Immer API 的功能和用法。

## 核心 API

### 1. produce

`produce` 是 Immer 的核心 API，接受一个初始状态和一个函数（称为 recipe），允许您对草稿进行修改。

#### 基本用法：

```typescript
import { produce } from 'immer'

// 基本形式：初始状态 + recipe
const nextState = produce(currentState, draft => {
  // 可以直接修改 draft
  draft.x = 123
})
```

#### 三种调用方式：

**标准方式**：直接传入初始状态和 recipe

```typescript
const baseState = { count: 0, list: [1, 2] }
const nextState = produce(baseState, draft => {
  draft.count++
  draft.list.push(3)
})

console.log(baseState) // { count: 0, list: [1, 2] }
console.log(nextState) // { count: 1, list: [1, 2, 3] }
console.log(baseState === nextState) // false
console.log(baseState.list === nextState.list) // false
```

**柯里化方式**：仅传入 recipe，返回一个接受状态的函数

```typescript
// 创建一个可重用的更新函数
const increment = produce(draft => {
  draft.count++
})

const state1 = { count: 0 }
const state2 = increment(state1) // { count: 1 }
const state3 = increment(state2) // { count: 2 }
```

**带参数的柯里化方式**：recipe 可以接收额外参数

```typescript
const addValue = produce((draft, value) => {
  draft.count += value
})

const state1 = { count: 10 }
const state2 = addValue(state1, 5) // { count: 15 }
```

#### 不对状态做更改的情况:

如果在 recipe 中没有对草稿进行任何修改，produce 将返回原始状态对象，保持引用相等：

```typescript
const baseState = { count: 0 }
const nextState = produce(baseState, draft => {
  // 什么也不做
})

console.log(baseState === nextState) // true
```

### 2. produceWithPatches

与 `produce` 类似，但返回一个包含三个元素的数组：[下一个状态, 补丁数组, 逆补丁数组]。

```typescript
import { produceWithPatches } from 'immer'

const baseState = { users: [{ name: 'John' }] }
const [nextState, patches, inversePatches] = produceWithPatches(baseState, draft => {
  draft.users.push({ name: 'Mike' })
  draft.users[0].name = 'John Doe'
})

console.log(patches)
/* 输出:
[
  { op: 'replace', path: ['users', 0, 'name'], value: 'John Doe' },
  { op: 'add', path: ['users', 1], value: { name: 'Mike' } }
]
*/
```

### 3. createDraft 和 finishDraft

允许手动创建和完成草稿，适用于不能立即完成状态修改的场景。

```typescript
import { createDraft, finishDraft } from 'immer'

const baseState = { count: 0 }
// 创建草稿
const draft = createDraft(baseState)

// 在任何时候修改草稿
draft.count++
draft.newField = 'test'

// 完成草稿并生成最终状态
const nextState = finishDraft(draft)

console.log(nextState) // { count: 1, newField: 'test' }
```

也可以使用 patch 监听器：

```typescript
let patches = []
let inversePatches = []

const draft = createDraft(baseState)
draft.count += 10
const nextState = finishDraft(draft, (p, ip) => {
  patches = p
  inversePatches = ip
})
```

### 4. applyPatches

应用由 `produceWithPatches` 或带补丁监听器的 `produce` 生成的补丁。

```typescript
import { applyPatches, produce } from 'immer'

const baseState = { users: [{ name: 'John' }] }
let patches = []

const nextState = produce(
  baseState,
  draft => {
    draft.users.push({ name: 'Mike' })
  },
  p => {
    patches = p
  }
)

// 将原始状态回退到 nextState
const recreatedState = applyPatches(baseState, patches)
console.log(recreatedState === nextState) // 结构相同但不是同一对象

// 在未来应用更多补丁
const futureState = applyPatches(nextState, [
  { op: 'replace', path: ['users', 0, 'name'], value: 'John Doe' }
])
```

## 工具函数

### 5. current

创建草稿的当前状态的快照（不冻结），用于调试或将草稿值安全地传递给外部。

```typescript
import { current, produce } from 'immer'

const nextState = produce({ count: 0 }, draft => {
  draft.count++
  // 查看当前修改后的草稿状态
  console.log(current(draft)) // { count: 1 }

  // 很适合在调试中使用
  // console.log(draft); // 输出 Proxy 对象，不便于查看
})
```

### 6. original

获取草稿对应的原始对象。

```typescript
import { original, produce } from 'immer'

const baseState = { users: [{ name: 'John' }] }
const nextState = produce(baseState, draft => {
  // 从草稿中获取原始对象
  const originalUser = original(draft.users[0])
  console.log(originalUser === baseState.users[0]) // true

  // 修改草稿不会影响原始对象
  draft.users[0].name = 'John Doe'
  console.log(originalUser.name) // 'John'
})
```

### 7. isDraft

检查值是否为 Immer 草稿。

```typescript
import { isDraft, produce } from 'immer'

const baseState = { count: 0 }

produce(baseState, draft => {
  console.log(isDraft(draft)) // true
  console.log(isDraft(baseState)) // false
})
```

### 8. isDraftable

检查值是否可以被制作成草稿。

```typescript
import { isDraftable } from 'immer'

console.log(isDraftable({})) // true
console.log(isDraftable([])) // true
console.log(isDraftable(new Map())) // true (启用 enableMapSet 后)
console.log(isDraftable(3)) // false
console.log(isDraftable(new Date())) // false
```

### 9. freeze

冻结对象，使其不可变。浅冻结或深冻结取决于第二个参数。

```typescript
import { freeze } from 'immer'

const obj = { a: 1, b: { c: 2 } }
// 浅冻结（默认）
const frozen = freeze(obj)
frozen.a = 2 // 抛出错误
frozen.b.c = 3 // 可以修改（浅冻结）

// 深冻结
const deepFrozen = freeze(obj, true)
deepFrozen.b.c = 3 // 抛出错误
```

## 特殊常量

### 10. nothing

特殊的哨兵值，从 recipe 返回它可以将状态设置为 undefined。

```typescript
import { produce, nothing } from 'immer'

const state = { a: 1, b: 2 }

const nextState1 = produce(state, draft => {
  // 移除 b 属性
  delete draft.b
})
console.log(nextState1) // { a: 1 }

const nextState2 = produce(state, draft => {
  // 使用 nothing 返回 undefined
  return nothing
})
console.log(nextState2) // undefined
```

### 11. immerable

符号常量，用于在类上标记，使 Immer 能够处理类实例。

```typescript
import { immerable, produce } from 'immer'

class Person {
  [immerable] = true // 标记此类为可草稿

  constructor(public name: string, public age: number) {}
}

const john = new Person('John', 30)
const olderJohn = produce(john, draft => {
  draft.age += 1
})

console.log(john.age) // 30
console.log(olderJohn.age) // 31
console.log(olderJohn instanceof Person) // true
```

## 配置 API

### 12. setAutoFreeze

控制 Immer 是否自动冻结通过 produce 产生的对象。

```typescript
import { setAutoFreeze, produce } from 'immer'

// 默认为 true
setAutoFreeze(false)

const baseState = { count: 0 }
const nextState = produce(baseState, draft => {
  draft.count++
})

// 不会被冻结
nextState.count = 100 // 不会抛出错误
```

### 13. setUseStrictShallowCopy

设置是否使用严格的浅拷贝（复制非枚举属性、getter/setter）。

```typescript
import { setUseStrictShallowCopy, produce } from 'immer'

// 启用严格模式拷贝
setUseStrictShallowCopy(true)

const baseState = {}
Object.defineProperty(baseState, 'hidden', {
  enumerable: false,
  value: 'secret'
})

const nextState = produce(baseState, draft => {
  // 空操作
})

// 严格模式下，非枚举属性也会被拷贝
console.log(Object.getOwnPropertyDescriptor(nextState, 'hidden')) // 值被保留
```

### 14. enableMapSet

启用 Map 和 Set 的支持。

```typescript
import { enableMapSet, produce } from 'immer'

// 启用对 Map 和 Set 的支持
enableMapSet()

const baseState = {
  map: new Map([['key', 'value']]),
  set: new Set(['item'])
}

const nextState = produce(baseState, draft => {
  draft.map.set('key2', 'value2')
  draft.set.add('item2')
})

console.log(baseState.map.has('key2')) // false
console.log(nextState.map.has('key2')) // true
console.log(nextState.set.has('item2')) // true
```

### 15. enablePatches

启用补丁功能 (produceWithPatches, applyPatches)。

```typescript
import { enablePatches, produceWithPatches } from 'immer'

// 启用补丁
enablePatches()

// 现在可以使用 produceWithPatches 和 applyPatches
```

## 类型工具

### 16. castDraft

TypeScript 工具函数，将不可变类型转换为草稿类型。

```typescript
import { castDraft, produce } from 'immer'

interface User {
  readonly name: string
  readonly age: number
}

const addYear = produce((draft: User) => {
  // 错误，name 是只读的
  // draft.name = "New name";

  // 解决方法
  const mutableDraft = castDraft(draft)
  mutableDraft.name = 'New name'
})
```

### 17. castImmutable

TypeScript 工具函数，将可变类型转换为不可变类型。

```typescript
import { castImmutable } from 'immer'

interface MutableUser {
  name: string
  age: number
}

// 转换为只读类型
const user: MutableUser = { name: 'John', age: 30 }
const immutableUser = castImmutable(user)

// TypeScript 错误，但实际运行时不会有影响
// immutableUser.name = 'New name';
```

## Immer 类

创建自定义配置的 Immer 实例。

```typescript
import { Immer } from 'immer'

// 创建自定义 Immer 实例
const myImmer = new Immer({
  autoFreeze: false,
  useStrictShallowCopy: true
})

// 使用自定义实例的 produce
const nextState = myImmer.produce({ count: 0 }, draft => {
  draft.count++
})
```

## 使用模式与最佳实践

### 1. 递归更新深层对象

```typescript
import { produce } from 'immer'

const state = {
  nested: {
    deeply: {
      objects: {
        count: 0
      }
    }
  }
}

const nextState = produce(state, draft => {
  draft.nested.deeply.objects.count++
})
```

### 2. 条件更新

```typescript
import { produce } from 'immer'

const toggle = produce((draft, id: number) => {
  const item = draft.items.find(item => item.id === id)
  if (item) {
    item.active = !item.active
  }
})

const state = {
  items: [
    { id: 1, active: false },
    { id: 2, active: false }
  ]
}
const nextState = toggle(state, 1)
```

### 3. 批量更新和复合更新

```typescript
import { produce } from 'immer'

const state = { items: [1, 2, 3, 4, 5] }

// 组合多个简单更新
const add = n =>
  produce(draft => {
    draft.items.push(n)
  })
const removeFirst = produce(draft => {
  draft.items.shift()
})

// 链式调用
const nextState = removeFirst(add(6)(state))
console.log(nextState) // { items: [2, 3, 4, 5, 6] }
```

### 4. 使用 patches 进行撤销/重做

```typescript
import { produceWithPatches, applyPatches } from 'immer'

let state = { count: 0 }
const history = []
const undoStack = []

function updateState(recipe) {
  const [nextState, patches, inversePatches] = produceWithPatches(state, recipe)
  history.push(patches)
  undoStack.push(inversePatches)
  state = nextState
  return state
}

function undo() {
  if (undoStack.length === 0) return state
  const inversePatches = undoStack.pop()
  history.pop()
  state = applyPatches(state, inversePatches)
  return state
}

// 使用示例
updateState(draft => {
  draft.count += 5
}) // count: 5
updateState(draft => {
  draft.count *= 2
}) // count: 10
console.log(state) // { count: 10 }
undo() // 撤销 count *= 2
console.log(state) // { count: 5 }
```

## 性能与注意事项

1. **结构共享**: Immer 只复制被修改的部分，未更改的部分保持引用不变
2. **避免在循环中使用 produce**: 每次调用 produce 都会创建新的代理和对象
3. **不可变库的调用开销**: 对于非常频繁的简单更新，直接使用对象扩展可能更高效
4. **自动冻结**: 开发环境中建议启用，生产环境可关闭提高性能

```typescript
import { setAutoFreeze } from 'immer'

if (process.env.NODE_ENV !== 'production') {
  setAutoFreeze(true) // 开发环境冻结对象，帮助发现错误
} else {
  setAutoFreeze(false) // 生产环境关闭冻结，提高性能
}
```

---

通过以上 API 和示例，你可以充分利用 Immer 的强大功能，以更简洁、直观的方式处理不可变数据。Immer 特别适合与 React、Redux 等状态管理库配合使用，大大简化状态更新逻辑。

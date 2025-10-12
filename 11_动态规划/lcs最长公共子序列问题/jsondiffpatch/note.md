好的，我们来详细讲解 `jsondiffpatch` 及其使用方式。

### `jsondiffpatch` 是什么？

`jsondiffpatch` 是一个 JavaScript 库，用于对 JSON 对象进行差异比较（diff）、应用补丁（patch）和撤销补丁（unpatch）。

它的核心思想是，当你有两个版本的 JSON 对象时，`jsondiffpatch` 可以生成一个非常紧凑的 "delta" 对象，这个 delta 对象描述了从第一个版本变动到第二个版本所需的所有更改。然后，你可以将这个 delta 应用到原始对象上，得到新版本的对象，或者从新版本对象上“撤销”这个 delta，从而恢复到原始版本。

这在很多场景下都非常有用，例如：

- **数据同步**：在客户端和服务器之间只传输数据的变动部分，而不是整个对象，从而节省带宽。
- **版本控制/历史记录**：存储对象的每次变更记录，可以轻松地回溯到任意历史版本。
- **协同编辑**：实时合并多个用户对同一个数据模型的修改。

### 主要特点

- **支持复杂结构**：可以处理嵌套对象、数组、字符串、数字等各种 JSON 数据类型。
- **高效的数组差异算法**：使用 [LCS (最长公共子序列)](https://en.wikipedia.org/wiki/Longest_common_subsequence_problem) 算法来智能地检测数组中的元素移动、添加和删除，而不仅仅是按索引比较。
- **可逆操作**：生成的 delta 可以用于 `patch`（正向应用）和 `unpatch`（反向应用）。
- **紧凑的 Delta 格式**：生成的差异描述文件非常小。
- **支持文本差异**：可以选择对长字符串进行行级的文本差异比较。

---

### 安装

你可以通过 npm 或 yarn 将其安装到你的项目中。

```bash
npm install jsondiffpatch
```

---

### 基本用法

`jsondiffpatch` 的 API 非常直观，主要包含 `diff`, `patch`, `unpatch` 三个核心方法。

#### 1. 比较差异 (diff)

`diff` 方法接收两个 JSON 对象，并返回一个描述它们之间差异的 delta 对象。

**示例代码：**

```javascript
// 引入库
import * as jsondiffpatch from 'jsondiffpatch'

// 创建一个 jsondiffpatch 实例
const diffpatcher = jsondiffpatch.create()

// 原始对象
const left = {
  name: 'John Doe',
  age: 30,
  tasks: ['write code', 'fix bugs'],
  contact: {
    email: 'john.doe@example.com'
  }
}

// 修改后的对象
const right = {
  name: 'Johnathan Doe', // 修改
  age: 31, // 修改
  tasks: ['write code', 'review code', 'fix bugs'], // 数组中添加元素
  // contact 属性被删除
  status: 'active' // 新增属性
}

// 生成 delta
const delta = diffpatcher.diff(left, right)

// 打印 delta
console.log(JSON.stringify(delta, null, 2))
```

**输出的 Delta：**

```json
{
  "name": ["John Doe", "Johnathan Doe"],
  "age": [30, 31],
  "tasks": {
    "_t": "a",
    "1": ["review code"]
  },
  "contact": [
    {
      "email": "john.doe@example.com"
    },
    0,
    0
  ],
  "status": ["active"]
}
```

**Delta 格式解释：**

- **修改 (Modification)**: `"name": ["John Doe", "Johnathan Doe"]` 表示 `name` 属性的值从 `"John Doe"` 变成了 `"Johnathan Doe"`。
- **新增 (Addition)**: `"status": ["active"]` 表示新增了 `status` 属性，其值为 `"active"`。
- **删除 (Deletion)**: `"contact": [ { ... }, 0, 0 ]` 表示 `contact` 属性被删除了。`[oldValue, 0, 0]` 是删除的固定格式。
- **数组差异 (Array Diff)**:
  - `"_t": "a"` 表示这是一个数组（array）的差异。
  - `"1": ["review code"]` 表示在索引 `1` 的位置插入了新元素 `"review code"`。其他未变动的元素（'write code', 'fix bugs'）通过 LCS 算法被识别为“移动”或“保留”，因此不会出现在 delta 中，从而使 delta 更小。

#### 2. 应用补丁 (patch)

`patch` 方法将 `diff` 生成的 delta 应用到原始对象上，从而得到新版本的对象。

**示例代码：**

```javascript
// ... 接上面的代码

// 使用 delta 来更新 left 对象
// 注意：patch 会直接修改原始对象。如果不想修改，可以先克隆。
const patchedLeft = JSON.parse(JSON.stringify(left)) // 创建一个克隆
diffpatcher.patch(patchedLeft, delta)

// 验证 patchedLeft 是否和 right 对象深度相等
console.log(patchedLeft)
// 输出会和 right 对象的内容完全一样
/*
{
  name: 'Johnathan Doe',
  age: 31,
  tasks: [ 'write code', 'review code', 'fix bugs' ],
  status: 'active'
}
*/
```

#### 3. 撤销补丁 (unpatch)

`unpatch` 方法是 `patch` 的逆操作。它将 delta 应用到新版本的对象上，使其恢复到原始状态。

**示例代码：**

```javascript
// ... 接上面的代码

// 使用 delta 来从 right 对象恢复到 left 对象的状态
const unpatchedRight = JSON.parse(JSON.stringify(right)) // 创建一个克隆
diffpatcher.unpatch(unpatchedRight, delta)

// 验证 unpatchedRight 是否和 left 对象深度相等
console.log(unpatchedRight)
// 输出会和 left 对象的内容完全一样
/*
{
  name: 'John Doe',
  age: 30,
  tasks: [ 'write code', 'fix bugs' ],
  contact: { email: 'john.doe@example.com' }
}
*/
```

---

### 高级用法与配置

在创建 `jsondiffpatch` 实例时，可以传入一个配置对象来自定义其行为。

```javascript
const diffpatcher = jsondiffpatch.create({
  // 用于识别数组中的对象，而不是按索引。
  // 当数组中的对象顺序改变时非常有用。
  // 'id' 是对象中作为唯一标识的属性名。
  objectHash: (obj, index) => obj.id || `$$index:${index}`,

  // 对长字符串使用文本差异（需要额外引入 'jsondiffpatch/dist/formatters/annotated'）
  textDiff: {
    minLength: 60 // 只对长度超过60的字符串进行文本差异比较
  }
})
```

**`objectHash` 示例：**

如果你的数组包含带 `id` 的对象，使用 `objectHash` 会让 diff 结果更准确。

```javascript
const diffpatcherWithObjectHash = jsondiffpatch.create({
  objectHash: obj => obj.id
})

const listBefore = [
  { id: 'a', value: 1 },
  { id: 'b', value: 2 }
]

const listAfter = [
  { id: 'b', value: 3 }, // value 改变
  { id: 'a', value: 1 } // 顺序改变
]

const delta = diffpatcherWithObjectHash.diff(listBefore, listAfter)

console.log(JSON.stringify(delta, null, 2))
```

**输出：**

```json
{
  "_t": "a",
  "_0": [
    {
      "id": "b",
      "value": 2
    },
    {
      "id": "b",
      "value": 3
    },
    1
  ]
}
```

**解释：**

- `_t: "a"`: 数组差异。
- `_0`: 表示对 ID 为 `b` 的对象进行操作（因为 `objectHash` 返回了 `b`）。
- `[ { ...old }, { ...new }, 1 ]`: 这是一个移动+修改的组合。它识别出 ID 为 `b` 的对象从索引 `1` 移动到了索引 `0`，并且其 `value` 属性发生了变化。

如果没有 `objectHash`，`jsondiffpatch` 会认为索引 `0` 和 `1` 的元素都完全改变了，生成的 delta 会更大、更不直观。

---

## 类型

好的，我们来详细讲解 `jsondiffpatch` 的核心 TypeScript 类型定义。这些类型精确地描述了配置选项和差异（delta）对象的结构。

### `Options` 接口

这个接口定义了你在创建 `jsondiffpatch` 实例时可以传入的配置对象，用于自定义比较和打补丁的行为。

```typescript
export interface Options {
  // 用于在数组中唯一标识对象。
  // 如果提供，jsondiffpatch 会根据这个函数返回的哈希值（如对象的 id）来匹配元素，
  // 而不是仅仅根据它们在数组中的位置。这对于检测数组中对象的移动非常重要。
  objectHash?: (item: object, index?: number) => string | undefined

  // 仅在数组差异计算中使用。如果为 true，则强制按位置匹配，即使提供了 objectHash。
  matchByPosition?: boolean

  // 针对数组差异的特定配置
  arrays?: {
    // 是否检测元素的移动。默认为 true。如果为 false，元素顺序的改变会被视为删除和添加。
    detectMove?: boolean
    // 在移动的 delta 中是否包含元素的值。默认为 false，以减小 delta 体积。
    includeValueOnMove?: boolean
  }

  // 针对长字符串的文本差异配置
  textDiff?: {
    // 必须传入 diff-match-patch 库的实例。
    diffMatchPatch: typeof diff_match_patch
    // 字符串达到此最小长度时，才启用文本差异算法。
    minLength?: number
  }

  // 一个函数，用于在 diff 过程中过滤掉某些属性。
  // 如果函数返回 false，则该属性将被忽略，不会出现在 delta 中。
  propertyFilter?: (name: string, context: DiffContext) => boolean

  // 在将值（如新值或旧值）放入 delta 时，是否克隆它们。
  // 这可以防止原始对象和 delta 之间的意外引用共享。
  // 可以是一个布尔值，也可以是一个自定义的克隆函数。
  cloneDiffValues?: boolean | ((value: unknown) => unknown)

  // 如果为 true，对于被删除的属性，delta 中将不包含其原始值。
  // 这会使 DeletedDelta 从 [oldValue, 0, 0] 变为 [undefined, 0, 0]，从而减小 delta 体积，
  // 但也意味着你将无法从 delta 中恢复被删除的值。
  omitRemovedValues?: boolean
}
```

---

### Delta 类型

这些类型是 `jsondiffpatch` 的核心，它们定义了差异（delta）对象的格式。`Delta` 是一个联合类型，可以是以下任何一种具体形式。

- `export type Delta = AddedDelta | ModifiedDelta | DeletedDelta | ObjectDelta | ArrayDelta | MovedDelta | TextDiffDelta | undefined;`
  这是一个总的联合类型，表示一个差异可以是任何一种定义的 Delta 格式，或者 `undefined`（表示没有差异）。

#### 基础 Delta 类型

- `export type AddedDelta = [unknown];`
  **新增**：一个只包含一个元素的数组，该元素就是新增的值。

  - 示例：`{ "status": ["active"] }` 表示新增了 `status` 属性，值为 `'active'`。

- `export type ModifiedDelta = [unknown, unknown];`
  **修改**：一个包含两个元素的数组，第一个是旧值，第二个是新值。

  - 示例：`{ "age": [30, 31] }` 表示 `age` 属性的值从 `30` 变为了 `31`。

- `export type DeletedDelta = [unknown, 0, 0];`
  **删除**：一个包含三个元素的数组，第一个是旧值，后两个是固定的标记 `0, 0`。
  - 示例：`{ "lastLogin": ["2025-10-11", 0, 0] }` 表示 `lastLogin` 属性被删除了。

#### 结构化 Delta 类型

- `export interface ObjectDelta { [property: string]: Delta; }`
  **对象差异**：一个普通对象，其键是被改变的属性名，其值是描述该属性变化的 `Delta` 对象。

  - 示例：`{ "contact": { "email": ["a@a.com", "b@b.com"] } }`

- `export interface ArrayDelta { _t: "a"; ... }`
  **数组差异**：一个特殊的对象，有固定的 `_t: "a"` 属性来标识它是一个数组 delta。
  - `[index: number | \`${number}\`]: Delta;`
    - 用于**插入**或**修改**。键是数组索引，值是一个 `AddedDelta` (插入) 或 `ModifiedDelta` (修改)。
  - `[index: \`\_${number}\`]: DeletedDelta | MovedDelta;`
    - 用于**删除**或**移动**。键是带下划线前缀的**原始索引**。值是一个 `DeletedDelta` (删除) 或 `MovedDelta` (移动)。
  - 示例：`{ "_t": "a", "1": ["review code"], "_2": [ "fix bugs", 0, 0 ] }` 表示在索引 `1` 处插入了 `"review code"`，并删除了原始索引为 `2` 的元素 `"fix bugs"`。

#### 特殊 Delta 类型

- `export type MovedDelta = [unknown, number, 3];`
  **移动**：仅在 `ArrayDelta` 中使用。一个三元素数组，`[value, newIndex, 3]`。`value` 是被移动的元素值（如果 `includeValueOnMove` 为 `false` 则为 `undefined`），`newIndex` 是它移动到的新位置，`3` 是移动的固定标记。

  - 示例：在 `ArrayDelta` 中，`"_3": ["", 1, 3]` 表示原始索引为 `3` 的元素被移动到了新的索引 `1`。

- `export type TextDiffDelta = [string, 0, 2];`
  **文本差异**：用于长字符串。一个三元素数组，`[diff, 0, 2]`。`diff` 是由 `diff-match-patch` 库生成的特殊格式字符串，`0, 2` 是文本差异的固定标记。

---

好的，我们来对 `jsondiffpatch` 的 types.d.ts 文件进行一次全面而深入的讲解。这个文件是理解 `jsondiffpatch` 工作原理和如何在 TypeScript 项目中正确使用它的关键。

我们将分块解析这个文件：

1.  **依赖导入 (Imports)**
2.  **配置项 (`Options` interface)**
3.  **核心：Delta 类型定义**
4.  **内部机制 (`Filter` interface)**
5.  **工具函数：类型守卫 (Type Guards)**

---

### 1. 依赖导入 (Imports)

```typescript
import type { diff_match_patch } from '@dmsnell/diff-match-patch'
import type Context from './contexts/context.js'
import type DiffContext from './contexts/diff.js'
```

- `diff_match_patch`: 这是一个强大的文本比较库（由 Google 的 Neil Fraser 开发）。`jsondiffpatch` 并不直接包含它，而是将其作为一个可选的对等依赖（peer dependency）。如果你想对长字符串进行高效的文本级别差异比较（而不是简单地将整个字符串视为已修改），你就需要安装并在这里提供这个库的实例。
- `Context` / `DiffContext`: 这些是 `jsondiffpatch` 内部使用的类型。它们代表了在执行 `diff`、`patch` 或 `unpatch` 操作期间的“上下文”对象。这个上下文对象会携带状态，例如当前正在比较的左右两个值、完整的根对象等。`propertyFilter` 函数的第二个参数就是 `DiffContext` 类型，允许你基于当前的比较环境来决定是否过滤某个属性。

---

### 2. 配置项 (`Options` interface)

这是 `jsondiffpatch` 最重要的配置接口，它允许你精细地控制差异比较的行为。

```typescript
export interface Options {
  // ... 字段 ...
}
```

- `objectHash?: (item: object, index?: number) => string | undefined;`
  **核心功能**。当比较数组时，默认情况下 `jsondiffpatch` 只能通过位置来猜测元素的对应关系。但如果数组中的对象顺序变了，这会导致大量的“删除”和“添加”操作。通过提供 `objectHash` 函数，你可以告诉 `jsondiffpatch` 如何唯一地识别一个对象（例如，通过其 `id` 或 `uuid` 属性）。这样，即使对象在数组中的位置移动了，`jsondiffpatch` 也能识别出这是一个“移动”操作，从而生成更精确、更小的 delta。

- `arrays?: { detectMove?: boolean; includeValueOnMove?: boolean; };`

  - `detectMove`: 默认为 `true`。是否启用数组元素的移动检测。如果关闭，元素的重新排序将被视为“删除”旧位置的元素并在新位置“添加”一个新元素。
  - `includeValueOnMove`: 默认为 `false`。当一个元素被检测到移动时，生成的 `MovedDelta` 是否包含该元素的值。设置为 `false` 可以让 delta 更小，但如果你需要仅通过 delta 就知道被移动的元素是什么，可以设为 `true`。

- `textDiff?: { minLength?: number; };`
  配置文本差异的行为。`minLength`（默认 `60`）指定了只有当字符串长度超过这个值时，才会启用 `diff-match-patch` 算法进行文本差异比较。对于短字符串，直接替换通常更高效。

- `propertyFilter?: (name: string, context: DiffContext) => boolean;`
  一个强大的过滤器。在比较两个对象的属性之前，会调用此函数。你可以根据属性名 (`name`) 和当前的差异上下文 (`context`) 来决定是否要包含这个属性的差异。如果函数返回 `false`，该属性将被完全忽略。这对于跳过像 `lastUpdated` 这样的时间戳或临时状态字段非常有用。

- `cloneDiffValues?: boolean | ((value: unknown) => unknown);`
  一个重要的安全选项。默认 `false`。当 `jsondiffpatch` 创建 delta 时，它会引用原始对象中的值。如果你在生成 delta 后修改了原始对象，delta 中的值也会被意外修改。将此项设为 `true` 会对放入 delta 的值进行深拷贝，从而避免这种副作用。你也可以提供一个自定义的克隆函数。

---

### 3. 核心：Delta 类型定义

这部分是 `jsondiffpatch` 的“语言”，定义了所有可能变更的表示方法。

- `export type Delta = ... | undefined;`
  这是一个总的联合类型，涵盖了所有可能的差异格式。如果两个对象完全相同，`diff` 的结果就是 `undefined`。

- `AddedDelta = [unknown];`
  **新增**：一个元素的数组，代表新增的值。

- `ModifiedDelta = [unknown, unknown];`
  **修改**：两个元素的数组，`[oldValue, newValue]`。

- `DeletedDelta = [unknown, 0, 0];`
  **删除**：三个元素的数组，`[oldValue, 0, 0]`。`[0, 0]` 是一个独特的标记，用于区分它与一个恰好包含两个 `0` 的普通数组的修改。

- `ObjectDelta { [property: string]: Delta; }`
  **对象差异**：一个对象，其键是发生变化的属性名，值是对应属性的 `Delta`。

- `ArrayDelta { _t: "a"; ... }`
  **数组差异**：这是最复杂的部分。它是一个特殊对象，由 `_t: "a"` 标记。

  - `[index: number | \`${number}\`]: Delta;`
    - 键是**新数组**的索引，值是 `AddedDelta`（插入）或 `ModifiedDelta`（修改）。
  - `[index: \`\_${number}\`]: DeletedDelta | MovedDelta;`
    - 键是**原数组**的索引，并带有 `_` 前缀。值是 `DeletedDelta`（删除）或 `MovedDelta`（移动）。这种设计允许在一个 delta 中同时表达基于新旧索引的操作。

- `MovedDelta = [unknown, number, 3];`
  **移动**：`[value, newIndex, 3]`。`value` 是被移动的值（可能为 `undefined`），`newIndex` 是它在**新数组**中的索引，`3` 是移动操作的标记。

- `TextDiffDelta = [string, 0, 2];`
  **文本差异**：`[diffString, 0, 2]`。`diffString` 是由 `diff-match-patch` 生成的补丁字符串，`2` 是文本差异的标记。

---

### 4. 内部机制 (`Filter` interface)

```typescript
export interface Filter<TContext extends Context<unknown>> {
  (context: TContext): void
  filterName: string
}
```

这是一个用于 `jsondiffpatch` 内部管道（pipeline）模式的接口。`jsondiffpatch` 的 `diff`、`patch` 等操作是通过一系列的“过滤器”链式处理的。每个过滤器负责处理一种特定的数据类型（如对象、数组、原始类型等）。这个接口定义了每个过滤器的形状：它是一个函数，接收一个上下文对象并对其进行操作，并且它还有一个 `filterName` 属性用于调试。这部分对于普通用户来说通常是透明的，但对于扩展库功能的高级用户则很重要。

---

### 5. 工具函数：类型守卫 (Type Guards)

```typescript
export declare function isAddedDelta(delta: Delta): delta is AddedDelta
export declare function isModifiedDelta(delta: Delta): delta is ModifiedDelta
// ... and so on for other delta types
```

这些是 TypeScript 的**类型守卫**函数。因为 `Delta` 是一个宽泛的联合类型，当你在代码中处理一个 `Delta` 对象时，TypeScript 并不知道它具体是哪一种（是 `AddedDelta` 还是 `ObjectDelta`？）。

这些函数允许你安全地缩小类型范围。例如：

```typescript
import { isAddedDelta, isObjectDelta } from 'jsondiffpatch/lib/types'

function processDelta(delta: Delta) {
  if (isAddedDelta(delta)) {
    // 在这个块内部，TypeScript 知道 delta 的类型是 AddedDelta ([unknown])
    console.log('Added value:', delta[0])
  } else if (isObjectDelta(delta)) {
    // 在这里，TypeScript 知道 delta 是一个 ObjectDelta
    console.log('Changed properties:', Object.keys(delta))
  }
  // ...
}
```

没有这些类型守卫，你就需要自己编写复杂的逻辑（如 `Array.isArray(delta) && delta.length === 1`）来判断 delta 的类型，而这些函数提供了官方、可靠且易读的方式来完成这项工作。

---

## formatters

好的，我们来补充 `formatters` 的使用示例。

`formatters` 用于将 `jsondiffpatch` 生成的 `delta` 对象转换成各种人类可读或机器可读的格式。

我将会在您的 demo.ts 文件末尾追加一个新代码块来演示 `console`, `html`, 和 `jsonpatch` 这三种常用格式化器的用法。

```typescript
// ...existing code...
  const delta = diffpatcherWithObjectHash.diff(listBefore, listAfter)

  console.log(JSON.stringify(delta, null, 2))
}

// 导入 formatters
import { formatters } from 'jsondiffpatch'

// 演示 formatters 的使用
console.log('\n--- 演示 Formatters ---')
{
  const v1 = {
    name: 'Alice',
    hobbies: ['reading', 'swimming'],
    deleted: true
  }

  const v2 = {
    name: 'Alicia',
    hobbies: ['reading', 'cycling', 'swimming'],
    status: 'active'
  }

  const delta = diffpatcher.diff(v1, v2)

  console.log('\n原始 Delta 对象:')
  console.log(JSON.stringify(delta, null, 2))

  // 1. Console Formatter (输出人类可读的文本)
  // formatters.console.format(delta, v1) 会在控制台输出彩色的 diff 结果。
  // 在不支持彩色的环境中，它会使用 +/- 等符号表示。
  console.log('\n1. Console Formatter Output:')
  // 注意：直接在 Node.js 终端运行时，你会看到带颜色的输出
  formatters.console.log(delta)

  // 2. HTML Formatter (生成 HTML 可视化 diff)
  // 这会生成一个完整的 HTML 字符串，你可以将其保存为 .html 文件并在浏览器中打开查看。
  console.log('\n2. HTML Formatter Output (部分):')
  const htmlOutput = formatters.html.format(delta, v1)
  console.log(htmlOutput.substring(0, 300) + '...') // 只打印前300个字符作为示例

  // 3. JSON Patch Formatter (生成 RFC 6902 标准格式)
  // 这对于与支持标准 JSON Patch 的其他系统（如某些 API 网关、数据库）集成非常有用。
  console.log('\n3. JSON Patch Formatter Output:')
  const jsonpatchOutput = formatters.jsonpatch.format(delta)
  console.log(JSON.stringify(jsonpatchOutput, null, 2))
}
```

---

## with-text-diffs

好的，我们来详细讲解 with-text-diffs.ts 这个文件。

这个文件在 `jsondiffpatch` 包中扮演着一个非常重要的角色：**它提供了一个“开箱即用”的、预配置了文本差异功能的入口点**。

### 核心目的

标准版的 `jsondiffpatch` 为了保持核心库的轻量，默认不包含文本差异比较的功能。文本差异比较依赖于一个外部库 `@dmsnell/diff-match-patch`。

with-text-diffs.ts 的目的就是为了方便用户：它将 `diff-match-patch` 库预先集成好，并导出一系列函数，让用户无需手动配置就能直接使用带有文本差异比较能力的 `jsondiffpatch`。

---

### 代码逐段解析

#### 1. 导入 (Imports)

```typescript
import { diff_match_patch } from '@dmsnell/diff-match-patch'

import type Context from './contexts/context.js'
// ... 其他内部类型导入
import dateReviver from './date-reviver.js'
import DiffPatcher from './diffpatcher.js'
import type { Delta, Options } from './types.js'
```

- `import { diff_match_patch } from "@dmsnell/diff-match-patch";`: 这是最关键的导入。它引入了实际执行文本差异比较的库。
- 其他导入：引入了 `jsondiffpatch` 内部的核心类（`DiffPatcher`）、类型（`Options`, `Delta` 等）和辅助函数（`dateReviver`），以便构建和导出功能。

#### 2. 导出类型和核心类 (Exports)

```typescript
export { dateReviver, DiffPatcher }

export type * from './types.js'
export type { Context, DiffContext, PatchContext, ReverseContext }
```

这部分代码的作用是**创建一个干净的公共 API**。它将所有底层的类型定义和核心的 `DiffPatcher` 类从这个入口点重新导出。这样做的好处是，用户只需要从 `jsondiffpatch/with-text-diffs` 这一个地方导入，就能获得所有需要的类型和类，而无需关心库内部的文件结构。

#### 3. `create` 工厂函数

```typescript
export function create(
  options?: Omit<Options, 'textDiff'> & {
    textDiff?: Omit<Options['textDiff'], 'diffMatchPatch'>
  }
) {
  return new DiffPatcher({
    ...options,
    textDiff: { ...options?.textDiff, diffMatchPatch: diff_match_patch }
  })
}
```

这是这个文件的核心功能之一。它是一个增强版的 `create` 函数。

- **复杂的类型定义**:

  - `Omit<Options, "textDiff">`: 意味着用户可以传入除了 `textDiff` 之外的所有标准 `Options`。
  - `& { textDiff?: Omit<Options["textDiff"], "diffMatchPatch">; }`: 意味着用户**可以**提供一个 `textDiff` 配置对象，但是这个对象里面**不能**包含 `diffMatchPatch` 属性。
  - **为什么这么设计？** 因为这个 `create` 函数的职责就是帮你注入 `diffMatchPatch`。这个类型定义从根本上防止了用户错误地覆盖它。用户仍然可以提供 `textDiff` 的其他配置，比如 `minLength`。

- **实现逻辑**:
  - `return new DiffPatcher({ ... })`: 它创建并返回一个新的 `DiffPatcher` 实例。
  - `...options`: 将用户传入的所有配置（如 `objectHash`）展开。
  - `textDiff: { ...options?.textDiff, diffMatchPatch: diff_match_patch }`: 这是关键。它创建了一个 `textDiff` 配置：
    1.  它首先展开用户可能提供的 `textDiff` 配置（`...options?.textDiff`），比如 `minLength`。
    2.  然后，它强制性地将 `diffMatchPatch` 属性设置为导入的 `diff_match_patch` 库。这确保了文本差异功能一定会被正确配置。

#### 4. 单例模式的便捷函数 (`diff`, `patch`, `unpatch`, etc.)

```typescript
let defaultInstance: DiffPatcher

export function diff(left: unknown, right: unknown) {
  if (!defaultInstance) {
    defaultInstance = new DiffPatcher({
      textDiff: { diffMatchPatch: diff_match_patch }
    })
  }
  return defaultInstance.diff(left, right)
}

// patch, unpatch, reverse, clone 函数与 diff 结构类似
```

这部分代码提供了一组更便捷的顶层函数。

- **设计模式**: 这里使用了**惰性初始化的单例模式 (Lazy Initialized Singleton)**。

  - `let defaultInstance: DiffPatcher;`: 在模块作用域内声明一个变量 `defaultInstance`，但不立即初始化。
  - `if (!defaultInstance) { ... }`: 在 `diff`, `patch` 等函数第一次被调用时，会检查 `defaultInstance` 是否存在。如果不存在，就创建一个新的 `DiffPatcher` 实例并赋值给它。这个实例被创建时，同样自动配置了文本差异功能。
  - `return defaultInstance.diff(...)`: 后续所有对这些函数的调用都会复用这个已经创建好的 `defaultInstance`。

- **为什么提供这些函数？**
  - **简单性**: 对于那些不需要复杂配置、只想快速进行一次 diff 或 patch 操作的用户来说，直接调用 `diff(a, b)` 比先 `create()` 一个实例再调用 `instance.diff(a, b)` 要简单得多。
  - **性能**: 通过复用同一个 `defaultInstance`，避免了每次调用都重新创建 `DiffPatcher` 实例的开销。

### 总结

with-text-diffs.ts 文件是 `jsondiffpatch` 提供给用户的一个高级、便捷的封装层。它解决了“默认功能”与“可选重度功能”之间的平衡问题。

- **对于需要自定义行为的库作者或高级用户**，他们可以使用 `create()` 函数，传入自己的配置（如 `objectHash`），同时免费获得自动配置好的文本差异功能。
- **对于只需要快速进行一次差异比较的普通用户**，他们可以直接使用 `diff()`, `patch()` 等顶层函数，享受最简单的 API，并且这些 API 默认就具备了强大的文本差异比较能力。

# TypeScript's unknown type and type variance

https://marijnhaverbeke.nl/blog/unknown-type-variance.html

这篇文章由开发者 Marijn Haverbeke 撰写，深入探讨了 TypeScript 类型系统中一个既强大又棘手的领域：**`unknown` 类型与类型变变性（Variance）的交互**，特别是在处理泛型数据结构时。

为了让你深入理解，我将文章拆解为几个核心部分进行详细讲解：`unknown` 的本质、变变性的陷阱、以及作者提出的解决方案。

### 1. 三种特殊的“顶级/底层”类型

文章首先通过对比，明确了 TypeScript 中三个最关键的特殊类型：

- **`any` (任意类型)**
  - **地位**：既是所有类型的父类型（Supertype），也是所有类型的子类型（Subtype）。
  - **作用**：完全关闭类型检查。
  - **问题**：它会破坏类型系统的健全性。作者形容它像是“在类型系统上开了一个大到能开卡车通过的洞”。
- **`unknown` (未知类型)**
  - **地位**：是所有类型的父类型（类似 `any`），但**仅**是 `any` 和它自己的子类型。
  - **作用**：表示“这里有一个值，但我不知道它是啥”。它是安全的 `any`。
  - **关键约束**：你不能直接使用 `unknown` 类型的值（比如调用方法或访问属性），必须先进行类型检查（Narrowing）或断言（Casting）。
- **`never` (从不存在)**
  - **地位**：是所有类型的子类型，但**仅**是 `any` 和它自己的父类型。
  - **作用**：表示永远不会发生的值（例如抛出错误的函数的返回值，或不可能存在的联合类型分支）。

### 2. 核心问题：泛型与异构集合

作者举了一个现实场景的例子：**`Widget<T>`**（组件）。

```typescript
type Widget<T> = {
  parameter: T
  type: WidgetType<T> // 定义组件外观和行为
}
```

假设你需要一个数组来存储各种不同类型的组件（例如文本组件 `Widget<string>`，布尔组件 `Widget<boolean>`）。

- **旧方法（TypeScript 3.0 之前）**：使用 `Widget<any>[]`。
  - 但这会失去所有类型保护，不推荐。
- **理想方法**：使用 `Widget<unknown>[]`。
  - 既然 `unknown` 是所有类型的父类型，那么 `Widget<unknown>` 应该是所有 `Widget<T>` 的父类型，对吧？

**陷阱出现了：** 当你尝试这样做时：

```typescript
const textWidget: Widget<string> = { ... };
const list: Widget<unknown>[] = [];

// 报错: Type 'Widget<string>' is not assignable to type 'Widget<unknown>'.
list.push(textWidget);
```

这不仅反直觉，而且令人困惑。为什么 `string` 是 `unknown` 的子类型，但 `Widget<string>` 却不是 `Widget<unknown>` 的子类型？

### 3. "变变性"（Variance）详解

这就是文章的核心难点：**类型变变性（Type Variance）**。

#### 什么是协变（Covariant）与逆变（Contravariant）？

- **协变（Covariant）**: 如果 `A` 是 `B` 的子类型，那么 `T<A>` 也是 `T<B>` 的子类型。
  - _直觉理解_：容器与其内容的方向一致。
  - _例子_：只读属性通常是协变的。如果我需要一个能*读取* `unknown` 的盒子，给我一个能读取 `string` 的盒子也没问题（这在读取时其实也不完全对，但在只输出的情况下成立）。
- **逆变（Contravariant）**: 如果 `A` 是 `B` 的子类型，那么 `T<A>` 反而是 `T<B>` 的父类型（方向反了）。
  - _出现场景_：**函数参数**。
  - _逻辑_：如果一个函数需要处理 `Animal`（父类），你不能给它一个只知道如何处理 `Dog`（子类）的函数。因为调用者可能会传入 `Cat`，而那是你的 `Dog` 处理函数无法处理的。反过来，如果函数需要处理 `Dog`，你给它一个能处理所有 `Animal` 的函数是安全的。

#### `Widget` 的问题所在

让我们看作者给出的结构：

```typescript
type Widget<T> = {
  parameter: T // 这里 T 作为属性值，通常是协变的
  type: {
    render: (parameter: T) => Pixels // ⚡️ 关键点在这里！
  }
}
```

在 `render` 函数中，`T` 出现在了**参数位置**。

1.  函数参数是**逆变**的。
2.  意味着：要让 `(p: string) => void` 赋值给 `(p: unknown) => void`，前者必须能处理所有 `unknown` 类型的数据。
3.  显然，一个只接受 `string` 的函数无法处理 `unknown`（比如数字）的数据。
4.  因此，`Widget<string>` **不是** `Widget<unknown>` 的子类型。

这就是为什么 `unknown` 经常导致复杂的类型错误：一旦泛型参数出现在函数输入位置，父子关系就翻转了。

### 4. 解决方案：投影类型（Projected Type）

既然问题的根源在于那个逆变的函数参数（`renderer`），作者提出的解决方案是创建一个**移除了逆变部分**的新类型。

**方法：使用 `Omit` 工具类型**

```typescript
// 创建一个“投影”类型，去掉了导致逆变的 "type" 字段
type AnyWidget = Omit<Widget<unknown>, 'type'>
```

**深入分析这个解决方案：**

1.  `Widget<unknown>` 本身包含了导致冲突的结构。
2.  通过 `Omit`，我们切除了导致不兼容的部分（即包含函数参数的那部分）。
3.  剩下的部分（如 `parameter: T`）通常是协变的。
4.  现在，`AnyWidget[]` 可以接收 `Widget<string>`、`Widget<boolean>` 等，因为剩余的结构是兼容的。

### 5. 总结与权衡

作者承认这种方法有权衡：

- **优点**：你可以构建一个异构的集合（数组），并且依然保持类型安全（比 `any` 安全）。
- **缺点**：当你从 `AnyWidget[]` 中取出元素想要真正使用它（调用 `render`）时，你需要进行**类型断言（Type Casting）**，将其转回完整的 `Widget<unknown>` 或具体类型。

**代码示例演示最终方案：**

```typescript
type Pixels = { data: string }

type Widget<T> = {
  parameter: T
  type: {
    render: (parameter: T) => Pixels
  }
}

// 1. 尝试直接使用 Widget<unknown> (会失败)
function demoFail() {
  const textWidget: Widget<string> = {
    parameter: 'hello',
    type: { render: s => ({ data: s }) }
  }

  // ❌ Error: Type 'Widget<string>' is not assignable to type 'Widget<unknown>'.
  // const list: Widget<unknown>[] = [textWidget];
}

// 2. 使用作者建议的 AnyWidget
type AnyWidget = Omit<Widget<unknown>, 'type'>

function demoSuccess() {
  const textWidget: Widget<string> = {
    parameter: 'hello',
    type: { render: s => ({ data: s }) }
  }

  // ✅ Works! Base properties are compatible
  const list: AnyWidget[] = [textWidget]

  // 使用时
  const item = list[0]
  // item.type.render() // ❌ Error: Property 'type' does not exist.

  // 需要断言回完整类型来使用 (在受控的范围内)
  const fullWidget = item as Widget<unknown>
  // 注意：这里实际调用依然危险，通常需要配合 Type Guard
}
```

**核心启示**：在使用 TypeScript 高级类型时，由于函数参数的逆变性，直接使用 `unknown` 作为泛型参数往往行不通。通过 `Omit` 创建一个宽容的“投影类型”是处理异构集合的一种优雅的逃生舱口。

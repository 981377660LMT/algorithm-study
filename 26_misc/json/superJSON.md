好的，我们来逐一讲解 `superjson` 导出的这些核心方法和类。

`superjson` 库的核心思想是扩展标准的 JSON，使其能够支持更多 JavaScript 的数据类型，如 `Date`, `Map`, `Set`, `BigInt`, `RegExp` 等，而标准的 `JSON.stringify` 和 `JSON.parse` 会丢失这些类型信息。

### 核心方法

#### 1. `stringify` 和 `parse`

这是最常用的一对方法，功能类似于 `JSON.stringify` 和 `JSON.parse`，但支持更多数据类型。

- `superjson.stringify(object)`: 将一个包含复杂数据类型（如 `Date` 对象）的 JavaScript 对象转换为一个特殊的 JSON 字符串。这个字符串不仅包含数据，还包含了类型信息的元数据。
- `superjson.parse(string)`: 将 `superjson.stringify` 生成的字符串解析回原始的 JavaScript 对象，并正确恢复其特殊数据类型。

**示例:**
在你提供的代码中，`superjson.stringify({ date: new Date(0) })` 会生成一个类似 `{"json":{"date":{"__dot_json":"Date","value":"1970-01-01T00:00:00.000Z"}}}` 的字符串，其中 `__dot_json` 字段就是 `superjson` 添加的元数据，用于在 `parse` 时识别出这是一个 `Date` 对象。

#### 2. `serialize` 和 `deserialize`

这是一对更底层的组合。`stringify` 和 `parse` 实际上是基于这两个方法实现的。它们将序列化过程分成了两步。

- `superjson.serialize(object)`: 将 JavaScript 对象序列化为一个普通的 JavaScript 对象（类型为 `SuperJSONResult`），而不是直接生成字符串。这个结果对象包含 `json`（可被 `JSON.stringify` 处理的数据）和 `meta`（类型元数据）两个部分。这在需要对序列化后的数据进行额外处理或通过不支持字符串传输的渠道（如 React Native 的 bridge）传递数据时很有用。
- `superjson.deserialize(payload)`: 接收 `serialize` 方法生成的 `SuperJSONResult` 对象，并将其反序列化为原始的 JavaScript 对象。

**示例:**
在你提供的代码中，`superjson.serialize({ date: new Date(0) })` 的输出会是一个对象，形如：

```json
{
  "json": { "date": "1970-01-01T00:00:00.000Z" },
  "meta": { "json.date": ["Date"] }
}
```

### 扩展与定制

`superjson` 也允许你注册自定义的类和类型转换器。

- `registerClass(v, options)`: 注册你自己的类。这样 `superjson` 在序列化和反序列化时就能正确地处理你的类实例。
- `registerCustom(transformer, name)`: 注册一个自定义的转换器，用于处理 `superjson` 默认不支持的任意数据类型。你需要提供一个 `is` 函数来判断对象是否属于该类型，以及 `serialize` 和 `deserialize` 函数。
- `registerSymbol(v, identifier)`: 注册 `Symbol`，使其可以被序列化。
- `allowErrorProps(...props)`: 默认情况下，序列化 `Error` 对象时只会包含 `name` 和 `message` 属性。使用此方法可以允许包含其他属性（如 `stack`）。

### 实例与静态方法

`superjson` 导出的方法（如 `superjson.stringify`）实际上是操作一个默认的、全局共享的 `SuperJSON` 实例。

你也可以创建自己的 `SuperJSON` 实例，这在需要不同序列化配置的场景下很有用。

```typescript
import SuperJSON from 'superjson'

// 创建一个独立的 superjson 实例
const mySuperJson = new SuperJSON()

// 使用实例方法
const jsonString = mySuperJson.stringify({ date: new Date() })
```

`new SuperJSON({ dedupe: boolean })` 构造函数中的 `dedupe` 参数用于处理重复引用的对象，如果设置为 `true`，对于同一个对象的多次引用，只有第一次会被完整序列化，后续的引用会被替换为 `null`，以避免循环引用和减小输出体积。

---

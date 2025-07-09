# Lodash 对象类型判断函数详解

这四个函数都用于判断不同类型的对象，但它们有着明确的区别和使用场景。让我详细解释每个函数的用途和区别。

## 1. isObject

**定义**：检查值是否为 `Object` 类型（不是原始值）

```typescript
import { isObject } from 'lodash'

// 基本用法
console.log(isObject({})) // true
console.log(isObject([])) // true
console.log(isObject(function () {})) // true
console.log(isObject(new Date())) // true
console.log(isObject(/regex/)) // true
console.log(isObject(null)) // false
console.log(isObject(undefined)) // false
console.log(isObject(42)) // false
console.log(isObject('string')) // false
console.log(isObject(true)) // false

// 实现原理
function isObject(value: any): boolean {
  const type = typeof value
  return value != null && (type === 'object' || type === 'function')
}
```

**关键点**：

- 包括所有非原始值类型（对象、数组、函数、日期等）
- 排除 `null`（虽然 `typeof null === 'object'`）
- 函数也被认为是对象

## 2. isObjectLike

**定义**：检查值是否为类对象（不是函数，但是对象）

```typescript
import { isObjectLike } from 'lodash'

// 基本用法
console.log(isObjectLike({})) // true
console.log(isObjectLike([])) // true
console.log(isObjectLike(new Date())) // true
console.log(isObjectLike(/regex/)) // true
console.log(isObjectLike(function () {})) // false - 关键区别！
console.log(isObjectLike(null)) // false
console.log(isObjectLike(undefined)) // false
console.log(isObjectLike(42)) // false
console.log(isObjectLike('string')) // false

// 实现原理
function isObjectLike(value: any): boolean {
  return value != null && typeof value === 'object'
}
```

**关键点**：

- 比 `isObject` 更严格，排除了函数
- 只判断 `typeof value === 'object'` 且不为 `null`

## 3. isArrayLikeObject

**定义**：检查值是否为类数组对象（既是类对象，又是类数组）

```typescript
import { isArrayLikeObject } from 'lodash'

// 基本用法
console.log(isArrayLikeObject([])) // true
console.log(isArrayLikeObject('string')) // false - 字符串是类数组但不是对象
console.log(isArrayLikeObject({ length: 0 })) // true
console.log(isArrayLikeObject({ 0: 'a', length: 1 })) // true
console.log(isArrayLikeObject({})) // false - 没有 length 属性
console.log(isArrayLikeObject(function () {})) // false - 虽然函数有 length，但不是类对象
console.log(isArrayLikeObject(null)) // false
console.log(isArrayLikeObject(42)) // false

// DOM 元素集合示例
const nodeList = document.querySelectorAll('div')
console.log(isArrayLikeObject(nodeList)) // true

// Arguments 对象示例
function testArgs() {
  console.log(isArrayLikeObject(arguments)) // true
}
testArgs(1, 2, 3)

// 实现原理
function isArrayLikeObject(value: any): boolean {
  return isObjectLike(value) && isArrayLike(value)
}

function isArrayLike(value: any): boolean {
  return (
    value != null &&
    typeof value !== 'function' &&
    typeof value.length === 'number' &&
    value.length >= 0 &&
    value.length <= Number.MAX_SAFE_INTEGER
  )
}
```

**关键点**：

- 必须同时满足 `isObjectLike` 和 `isArrayLike`
- 排除字符串（虽然字符串是类数组）
- 排除函数（虽然函数有 `length` 属性）

## 4. isPlainObject

**定义**：检查值是否为普通对象（由 `Object` 构造函数创建或原型为 `null`）

```typescript
import { isPlainObject } from 'lodash'

// 基本用法
console.log(isPlainObject({})) // true
console.log(isPlainObject({ a: 1 })) // true
console.log(isPlainObject(Object.create(null))) // true
console.log(isPlainObject(new Object())) // true

// 非普通对象
console.log(isPlainObject([])) // false
console.log(isPlainObject(new Date())) // false
console.log(isPlainObject(/regex/)) // false
console.log(isPlainObject(function () {})) // false
console.log(isPlainObject(null)) // false
console.log(isPlainObject(undefined)) // false
console.log(isPlainObject('string')) // false

// 自定义构造函数创建的对象
class MyClass {}
console.log(isPlainObject(new MyClass())) // false

// 继承的对象
const proto = { inherited: true }
const obj = Object.create(proto)
console.log(isPlainObject(obj)) // false

// 实现原理（简化版）
function isPlainObject(value: any): boolean {
  if (!isObjectLike(value) || Object.prototype.toString.call(value) !== '[object Object]') {
    return false
  }

  if (Object.getPrototypeOf(value) === null) {
    return true
  }

  let proto = value
  while (Object.getPrototypeOf(proto) !== null) {
    proto = Object.getPrototypeOf(proto)
  }

  return Object.getPrototypeOf(value) === proto
}
```

**关键点**：

- 最严格的对象判断
- 只接受普通的 `{}` 对象或 `Object.create(null)` 创建的对象
- 排除所有内置对象（Date、RegExp、Array 等）
- 排除自定义构造函数创建的实例

## 完整对比示例

```typescript
// 综合对比示例
const testCases = [
  {}, // 普通对象
  [], // 数组
  'string', // 字符串
  function () {}, // 函数
  new Date(), // 日期对象
  /regex/, // 正则对象
  null, // null
  undefined, // undefined
  42, // 数字
  { length: 0 }, // 类数组对象
  document.querySelectorAll('div'), // NodeList
  (function () {
    return arguments
  })(), // Arguments
  Object.create(null), // 无原型对象
  new (class MyClass {})() // 自定义类实例
]

console.table(
  testCases.map(value => ({
    value: String(value),
    isObject: isObject(value),
    isObjectLike: isObjectLike(value),
    isArrayLikeObject: isArrayLikeObject(value),
    isPlainObject: isPlainObject(value)
  }))
)
```

## 使用场景建议

1. **isObject**: 判断是否为非原始值，用于区分对象和原始类型
2. **isObjectLike**: 判断是否为真正的对象（排除函数），常用于数据处理
3. **isArrayLikeObject**: 判断是否可以像数组一样遍历的对象，用于类数组处理
4. **isPlainObject**: 判断是否为纯数据对象，常用于配置对象验证和深拷贝

## 总结

这四个函数的严格程度递增：

- `isObject`: 最宽松，包括所有非原始值
- `isObjectLike`: 排除函数
- `isArrayLikeObject`: 在 `isObjectLike` 基础上要求有 `length` 属性
- `isPlainObject`: 最严格，只接受普通对象

选择哪个函数取决于你的具体需求和对"对象"的定义。

找到具有 2 个许可证类型的类似代码

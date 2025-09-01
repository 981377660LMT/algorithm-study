好的，我们来详细讲解 Lodash 中用于计算数组差集的 `difference` 系列函数。

Lodash 中没有一个名为 `diff` 的函数，但它提供了一系列以 `difference` 开头的函数来实现类似的功能，它们的核心作用是找出第一个数组中存在，但在后续其他数组中不存在的元素。

这系列函数主要有三个：`_.difference`、`_.differenceBy` 和 `_.differenceWith`。

### 1. `_.difference(array, ...values)`

这是最基础的差集函数。它接收一个主数组和若干个其他数组，然后返回一个新数组，包含所有在主数组中但不在任何其他数组中的元素。

- **比较方式**：使用 `SameValueZero` 进行相等性比较，这基本上等同于 `===`（严格相等）。这意味着它能很好地处理原始类型（数字、字符串、布尔值等），但对于对象，它只比较引用地址，而不是对象的内容。

- **语法**：

  ```javascript
  _.difference(array, [...values])
  ```

- **参数**：

  - `array`: 要检查的主数组。
  - `[...values]`: 一个或多个要排除的值所在的数组。

- **示例**：

  **a. 基本用法（数字和字符串）**

  ```javascript
  const array1 = [1, 2, 3, 4, 5]
  const array2 = [3, 5, 7]

  // 找出在 array1 中但不在 array2 中的元素
  const result = _.difference(array1, array2)

  console.log(result)
  // => [1, 2, 4]
  ```

  **b. 对象比较（注意陷阱）**

  ```javascript
  const obj1 = { id: 1, name: 'A' }
  const obj2 = { id: 2, name: 'B' }
  const obj3 = { id: 1, name: 'A' } // 内容相同，但引用不同

  const arrayA = [obj1, obj2]
  const arrayB = [obj3]

  // 因为 obj1 和 obj3 的引用地址不同，所以 _.difference 认为它们是不同的元素
  const result = _.difference(arrayA, arrayB)

  console.log(result)
  // => [{ id: 1, name: 'A' }, { id: 2, name: 'B' }]
  // 结果是整个 arrayA，因为没有找到引用相同的对象
  ```

  对于对象数组，`_.difference` 通常不是你想要的，这时就需要 `_.differenceBy`。

### 2. `_.differenceBy(array, ...values, iteratee)`

这个函数是 `_.difference` 的增强版。它允许你提供一个 `iteratee`（迭代器），在比较前对每个元素应用这个迭代器，然后比较迭代器返回的结果。

- **比较方式**：对每个元素调用 `iteratee` 函数，然后比较返回的结果。
- **`iteratee` 可以是**：

  - 一个函数：`item => item.id`
  - 一个属性字符串：`'id'`

- **语法**：

  ```javascript
  _.differenceBy(array, [...values], iteratee)
  ```

- **示例**：

  **a. 比较对象数组的某个属性**

  ```javascript
  const arrayA = [
    { id: 1, name: 'Apple' },
    { id: 2, name: 'Banana' }
  ]
  const arrayB = [
    { id: 2, name: 'Banana V2' },
    { id: 3, name: 'Cherry' }
  ]

  // 根据 'id' 属性来找差集
  const result = _.differenceBy(arrayA, arrayB, 'id')

  console.log(result)
  // => [{ 'id': 1, 'name': 'Apple' }]
  // 因为 id 为 2 的对象在 arrayB 中也存在
  ```

  **b. 使用函数作为迭代器**

  ```javascript
  // 找出在 [2.1, 1.2] 中但不在 [2.3, 3.4] 中的整数部分
  const result = _.differenceBy([2.1, 1.2], [2.3, 3.4], Math.floor)

  console.log(result)
  // => [1.2]
  // 过程:
  // 1. Math.floor(2.1) -> 2
  // 2. Math.floor(1.2) -> 1
  // 3. Math.floor(2.3) -> 2
  // 4. Math.floor(3.4) -> 3
  // 5. 比较 [2, 1] 和 [2, 3] 的差集，得到 [1]
  // 6. 返回原始数组中对应的值 [1.2]
  ```

### 3. `_.differenceWith(array, ...values, comparator)`

这是最灵活、最强大的差集函数。它允许你提供一个自定义的 `comparator` (比较器) 函数来定义两个元素是否“相等”。

- **比较方式**：`comparator` 函数会被调用，接收两个参数 `(arrVal, othVal)`，分别来自主数组和其他数组。如果 `comparator` 返回 `true`，则认为两个元素相等。

- **语法**：

  ```javascript
  _.differenceWith(array, [...values], comparator)
  ```

- **示例**：

  **a. 自定义复杂的对象比较逻辑**
  假设我们认为只要 `id` 相同且 `type` 也相同的对象就是同一个对象。

  ```javascript
  const arrayA = [
    { id: 1, type: 'fruit' },
    { id: 2, type: 'vegetable' }
  ]
  const arrayB = [
    { id: 1, type: 'fruit', price: 10 },
    { id: 3, type: 'fruit' }
  ]

  const customComparator = (objA, objB) => {
    return objA.id === objB.id && objA.type === objB.type
  }

  const result = _.differenceWith(arrayA, arrayB, customComparator)

  console.log(result)
  // => [{ id: 2, type: 'vegetable' }]
  // 因为 arrayA 的第一个元素通过自定义比较器在 arrayB 中找到了匹配项
  ```

### 总结对比

| 函数                   | 用途                                                     | 比较方式                         |
| :--------------------- | :------------------------------------------------------- | :------------------------------- |
| **`_.difference`**     | 计算原始类型数组或基于对象引用的差集。                   | `===` (严格相等)                 |
| **`_.differenceBy`**   | 基于对象数组的**某个特定属性**或**转换后**的值计算差集。 | 比较 `iteratee` 的返回值。       |
| **`_.differenceWith`** | 基于**完全自定义的比较逻辑**计算差集，适用于复杂场景。   | 使用自定义的 `comparator` 函数。 |

### 如何选择？

- 如果你处理的是数字、字符串等简单类型，使用 `_.difference`。
- 如果你处理的是对象数组，并且想根据一个简单的键（如 `id`）来判断是否相同，使用 `_.differenceBy`。
- 如果你需要更复杂的比较逻辑（例如，比较多个属性，或者进行模糊匹配），使用 `_.differenceWith`。

好的，我们来详细讲解 MobX 的核心 API。

MobX 是一个简单、可扩展的状态管理库，它通过透明的函数式响应式编程（TFRP）使得状态管理变得简单和可扩展。

其核心思想是：**任何源自应用状态的东西都应该自动地获得。**

这包括 UI、数据序列化、服务器通信等。

MobX 的 API 主要围绕以下三个核心概念构建：

1.  **State (状态)**：驱动应用的数据。
2.  **Actions (动作)**：修改状态的代码。
3.  **Derivations (派生)**：任何可以从状态中自动派生出来的值。派生有两种形式：
    - **Computed values (计算值)**：从当前可观察状态中派生出的新值。
    - **Reactions (反应)**：当状态改变时需要自动执行的副作用（例如更新 UI）。

---

### 1. 创建可观察状态 (Creating Observable State)

这是 MobX 的基础。你需要告诉 MobX 你想跟踪哪些数据。

#### `makeObservable` & `makeAutoObservable`

这是现代 MobX 中将一个对象变为可观察的首选方式，通常在类的构造函数中使用。

- **`makeObservable(target, annotations)`**: 提供精确的控制，你需要为对象的每个属性指定注解（例如 `observable`, `action`, `computed`）。

  ```javascript
  import { makeObservable, observable, computed, action } from 'mobx'

  class Todo {
    id = Math.random()
    title = ''
    finished = false

    constructor(title) {
      // 第二个参数是注解对象
      makeObservable(this, {
        title: observable,
        finished: observable,
        unfinishedTasks: computed,
        toggle: action
      })
      this.title = title
    }

    get unfinishedTasks() {
      // 这是一个计算值
      return this.finished ? 0 : 1
    }

    toggle() {
      this.finished = !this.finished
    }
  }
  ```

- **`makeAutoObservable(target)`**: 一个更方便的 API，它会自动推断所有属性的类型。

  - `get`ters 会被推断为 `computed`。
  - `set`ters 会被推断为 `action`。
  - 函数/方法会被推断为 `action`。
  - 其他所有属性（包括对象和数组）会被推断为 `observable`。

  ```javascript
  import { makeAutoObservable } from 'mobx'

  class Timer {
    secondsPassed = 0

    constructor() {
      // 自动推断所有属性
      makeAutoObservable(this)
    }

    increase() {
      this.secondsPassed += 1
    }

    reset() {
      this.secondsPassed = 0
    }
  }
  ```

#### `observable(source)`

可以用来创建可观察的对象、数组和 Map。

- **`observable.box(value)`**: 如果你想让一个原始类型（如 string, number, boolean）的值本身成为可观察的，可以使用 `observable.box`。它会创建一个包含 `.get()` 和 `.set()` 方法的对象。

  ```javascript
  import { observable, autorun } from 'mobx'

  const temperature = observable.box(25)

  autorun(() => {
    console.log(`Current temperature: ${temperature.get()}°C`)
  })

  // 输出: "Current temperature: 25°C"

  temperature.set(30)
  // 输出: "Current temperature: 30°C"
  ```

---

### 2. 修改状态 (Modifying State)

MobX 鼓励将所有修改状态的代码都包裹在 `action` 中。

#### `action`

为什么使用 `action`？

1.  **批处理更新**：在一个 `action` 中对状态的所有修改会等到 action 执行完毕后，才通知所有观察者（Reactions）进行更新。这可以防止不必要的中间状态更新，提高性能。
2.  **代码组织**：它清晰地标记了代码中“会修改状态”的部分。

`action` 可以作为函数包装器或 `makeObservable` 中的注解使用（如上例所示）。

```javascript
import { observable, action, runInAction } from 'mobx'

const person = observable({
  firstName: 'John',
  lastName: 'Doe'
})

// 1. 使用 action 函数包装
const changeName = action((firstName, lastName) => {
  person.firstName = firstName
  person.lastName = lastName
})

changeName('Jane', 'Smith')

// 2. 在异步函数中使用 runInAction
async function fetchUser() {
  const response = await fetch('/api/user')
  const user = await response.json()

  // 在 await 之后，需要用 runInAction 来包裹状态修改
  runInAction(() => {
    person.firstName = user.firstName
    person.lastName = user.lastName
  })
}
```

#### `runInAction(fn)`

这是一个简单的工具函数，它接收一个函数并以 `action` 的方式执行它。这在异步代码中特别有用。

---

### 3. 派生和反应 (Derivations and Reactions)

#### `computed`

计算值可以根据其他可观察状态派生出来。MobX 会确保计算值在它的任何依赖项发生变化时自动更新，并且结果会被缓存。只有当依赖项改变时，它才会重新计算。

`computed` 可以作为 getter 或 `makeObservable` 中的注解使用。

```javascript
import { observable, computed, makeObservable } from 'mobx'

class OrderLine {
  price = 0
  amount = 1

  constructor(price) {
    makeObservable(this, {
      price: observable,
      amount: observable,
      total: computed // 标记为 computed
    })
    this.price = price
  }

  // 这个 getter 就是一个计算值
  get total() {
    console.log('Calculating total...')
    return this.price * this.amount
  }
}

const line = new OrderLine(10)
console.log(line.total) // 输出 "Calculating total..." 和 10
console.log(line.total) // 直接输出 10 (因为结果被缓存，没有重新计算)

line.amount = 5
console.log(line.total) // 依赖项改变，重新计算。输出 "Calculating total..." 和 50
```

#### `autorun`

`autorun` 会自动运行你提供给它的函数，并在其依赖的任何可观察状态发生变化时重新运行。它非常适合用于日志记录、网络请求等需要响应状态变化的副作用。

```javascript
import { observable, autorun } from 'mobx'

const user = observable({
  name: 'Guest',
  isLoggedIn: false
})

// autorun 会立即执行一次，然后每当 user.name 或 user.isLoggedIn 变化时再次执行
autorun(() => {
  if (user.isLoggedIn) {
    console.log(`Welcome, ${user.name}!`)
  } else {
    console.log('Please log in.')
  }
})

// 输出: "Please log in."

user.isLoggedIn = true
user.name = 'Alice'
// 输出: "Welcome, Alice!"
```

#### `reaction`

`reaction` 类似于 `autorun`，但提供了更精细的控制。它接收两个函数：

1.  **数据函数 (data function)**：跟踪可观察数据并返回一个值。
2.  **效果函数 (effect function)**：只有在数据函数返回的值发生变化时才会执行。

这可以避免不必要的副作用执行，性能更好。

```javascript
import { observable, reaction } from 'mobx'

const cart = observable({
  itemCount: 0,
  status: 'active'
})

// 第一个函数跟踪 itemCount
// 第二个函数只有在 itemCount 的返回值变化时才执行
const disposer = reaction(
  () => cart.itemCount,
  (count, previousCount) => {
    console.log(`Item count changed from ${previousCount} to ${count}`)
  }
)

cart.status = 'pending' // 不会触发 reaction，因为跟踪的数据没变
cart.itemCount = 1 // 输出: "Item count changed from 0 to 1"
cart.itemCount = 2 // 输出: "Item count changed from 1 to 2"

// 调用 reaction 返回的函数可以清除这个反应
disposer()
```

#### `when`

`when` 观察并运行给定的 `predicate` (第一个函数)，当其返回 `true` 时，执行 `effect` (第二个函数)。`when` 只会执行一次，然后自动被清理。

```javascript
import { observable, when } from 'mobx'

const person = observable({
  isReady: false
})

// 当 person.isReady 变为 true 时，执行回调
when(
  () => person.isReady,
  () => {
    console.log('Person is ready!')
  }
)

setTimeout(() => {
  person.isReady = true // 2秒后输出 "Person is ready!"
}, 2000)
```

---

### 总结

| API                  | 用途                                                            | 何时使用                                                                     |
| -------------------- | --------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `makeAutoObservable` | 自动将类或对象的所有属性转换为可观察的、计算的或动作。          | 在类中快速设置响应式状态，是现代 MobX 的首选。                               |
| `makeObservable`     | 精确控制类或对象中每个属性的类型（`observable`, `action` 等）。 | 当你需要覆盖默认行为或代码更具可读性时。                                     |
| `observable`         | 创建可观察的对象、数组、Map。                                   | 创建独立的响应式数据结构。                                                   |
| `action`             | 将状态修改包装起来，进行批处理。                                | 任何修改状态的地方都应该使用，以获得最佳性能和可预测性。                     |
| `computed`           | 从现有状态派生新值，并缓存结果。                                | 当你需要基于现有状态计算值，并且不希望在每次访问时都重新计算。               |
| `autorun`            | 自动运行一个函数，并在其依赖项变化时重新运行。                  | 用于简单的副作用，如日志记录或调试，当你不关心“什么”改变了，只关心“有”改变。 |
| `reaction`           | 当特定的数据发生变化时，才执行副作用。                          | 当你只想在特定数据变化时触发副作用，以优化性能。                             |
| `when`               | 当某个条件变为真时，执行一次副作用。                            | 用于一次性的、条件性的逻辑，例如等待数据加载完成。                           |

掌握这些核心 API，你就可以有效地使用 MobX 来管理你的应用状态了。

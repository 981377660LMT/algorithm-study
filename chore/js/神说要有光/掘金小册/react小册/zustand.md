近几年出了很多可以替代 redux 的优秀状态管理库，zustand 是其中最优秀的一个。
它的特点有很多：体积小、简单、支持中间件扩展。

`它的核心就是一个 create 函数，传入 state 来创建 store。`
`create 返回的函数可以传入 selector，取出部分 state 在组件里用。`

它的中间件和 redux 一样，就是一个高阶函数，可以对 get、set 做一些扩展。
zustand 内置了 immer、persist 等中间件，我们也自己写了一个 log 的中间件。

zustand 本身的实现也很简单，就是 getState、setState、subscribe 这些功能，然后再加上 **useSyncExternalStore** 来触发组件 rerender。

---

## 1. 安装

```bash
npm install zustand
```

---

## 2. 创建 Store

使用 `create` 创建一个 store。store 就是一个 hook。

```javascript
// store.js
import { create } from 'zustand'

const useStore = create(set => ({
  count: 0,
  increase: () => set(state => ({ count: state.count + 1 })),
  decrease: () => set(state => ({ count: state.count - 1 }))
}))

export default useStore
```

---

## 3. 在组件中使用

直接在组件中调用 store hook 获取和操作状态。

```javascript
import useStore from './store'

function Counter() {
  const { count, increase, decrease } = useStore()
  return (
    <div>
      <button onClick={decrease}>-</button>
      <span>{count}</span>
      <button onClick={increase}>+</button>
    </div>
  )
}
```

---

## 4. 选择性订阅（优化性能）

只订阅部分状态，避免组件不必要的重渲染。

```javascript
const count = useStore(state => state.count)
const increase = useStore(state => state.increase)
```

---

## 5. 中间件（如持久化）

Zustand 支持中间件，比如状态持久化：

```javascript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

const useStore = create(
  persist(
    set => ({
      count: 0,
      increase: () => set(state => ({ count: state.count + 1 }))
    }),
    { name: 'counter-storage' }
  )
)
```

create 方法的参数，它是一个接受 set、get、store 的三个参数的函数。
中间件是装饰器。

---

## 6. 异步操作

可以在 action 里写异步逻辑：

```javascript
const useStore = create(set => ({
  data: null,
  fetchData: async () => {
    const res = await fetch('/api/data')
    const data = await res.json()
    set({ data })
  }
}))
```

## 7. 其他使用

```js
// 回调函数可以拿到当前 state，或者调用 store.getState 也可以拿到 state。
useEffect(() => {
  useXxxStore.subscribe(state => {
    console.log(useXxxStore.getState())
  })
}, [])
```

---

## 总结

- 创建 store：`create`
- 组件中直接用 hook
- 支持选择性订阅
- 支持中间件（如持久化、devtools）
- 支持异步 action

Zustand 适合需要简单、灵活状态管理的 React 项目。

---

https://juejin.cn/book/7294082310658326565/section/7329768564628881434

```js
const createStore = createState => {
  let state
  const listeners = new Set()

  // zustand 在 set 状态的时候默认是合并，你也可以传一个 true 改成替换
  const setState = (partial, replace) => {
    const nextState = typeof partial === 'function' ? partial(state) : partial

    if (!Object.is(nextState, state)) {
      const previousState = state

      if (!replace) {
        state =
          typeof nextState !== 'object' || nextState === null
            ? nextState
            : Object.assign({}, state, nextState)
      } else {
        state = nextState
      }
      listeners.forEach(listener => listener(state, previousState))
    }
  }

  const getState = () => state

  const subscribe = listener => {
    listeners.add(listener)
    return () => listeners.delete(listener)
  }

  const destroy = () => {
    listeners.clear()
  }

  const api = { setState, getState, subscribe, destroy }

  state = createState(setState, getState, api)

  return api
}

function useStore(api, selector) {
  const [, forceRender] = useState(0)
  useEffect(() => {
    api.subscribe((state, prevState) => {
      const newObj = selector(state)
      const oldobj = selector(prevState)

      if (newObj !== oldobj) {
        forceRender(Math.random())
      }
    })
  }, [])
  return selector(api.getState())
}

export const create = createState => {
  const api = createStore(createState)

  const useBoundStore = selector => useStore(api, selector)

  Object.assign(useBoundStore, api)

  return useBoundStore
}
```

代码还可以进一步简化。
react **有一个 hook 就是用来定义外部 store 的，store 变化以后会触发 rerender**

```ts
function useStore(api, selector) {
  function getState() {
    return selector(api.getState())
  }

  return useSyncExternalStore(api.subscribe, getState)
}
```

---

教程
https://zustand.docs.pmnd.rs/getting-started/comparison

## updating state

- Flat updates
  使用新 state 调用提供的 set 函数，它将与 store 中的现有 state 浅合并
- Deeply nested object 深度嵌套的对象
  使用immer

  ```js
    immerInc: () =>
      set(produce((state: State) => { ++state.deep.nested.obj.count })),
  ```

## Flux inspired practice 受 Flux 启发的实践

1. **Single store**
   Your applications global state should be located in a single Zustand store.
   您的应用程序全局状态应位于单个 Zustand 存储中。
   If you have a large application, Zustand supports `splitting the store into slices`.
   如果您有一个大型应用程序，Zustand 支持将 store 拆分为 slice。

```js
import { create } from 'zustand'
import { createBearSlice } from './bearSlice'
import { createFishSlice } from './fishSlice'

export const useBoundStore = create((...a) => ({
  ...createBearSlice(...a),
  ...createFishSlice(...a)
}))

export const createFishSlice = set => ({
  fishes: 0,
  addFish: () => set(state => ({ fishes: state.fishes + 1 }))
})

export const createBearSlice = set => ({
  bears: 0,
  addBear: () => set(state => ({ bears: state.bears + 1 })),
  eatFish: () => set(state => ({ fishes: state.fishes - 1 }))
})

import { useBoundStore } from './stores/useBoundStore'

function App() {
  const bears = useBoundStore(state => state.bears)
  const fishes = useBoundStore(state => state.fishes)
  const addBear = useBoundStore(state => state.addBear)
  return (
    <div>
      <h2>Number of bears: {bears}</h2>
      <h2>Number of fishes: {fishes}</h2>
      <button onClick={() => addBear()}>Add a bear</button>
    </div>
  )
}

export default App
```

## TypeScript Guide

您**必须编写 create<T>**（）（...） 而不是编写 create（...）（ 注意额外的括号 （） 以及 type 参数），其中 T 是用于注释的状态类型。

```ts
import { create } from 'zustand'

interface BearState {
  bears: number
  increase: (by: number) => void
}

const useBearStore = create<BearState>()(set => ({
  bears: 0,
  increase: by => set(state => ({ bears: state.bears + by }))
}))
```

为什么我们不能简单地从初始状态推断类型呢？
`TLDR：Because state generic T is invariant.`

## Api

```ts
createStore<T>()(stateCreatorFn: StateCreator<T, [], []>): StoreApi<T>
createWithEqualityFn<T>()(stateCreatorFn: StateCreator<T, [], []>, equalityFn?: (a: T, b: T) => boolean): UseBoundStore<StoreApi<T>>
create<T>()(stateCreatorFn: StateCreator<T, [], []>): UseBoundStore<StoreApi<T>>
```

## Hooks

```ts
useStore<StoreApi<T>, U = T>(store: StoreApi<T>, selectorFn?: (state: T) => U) => UseBoundStore<StoreApi<T>>

useShallow<T, U = T>(selectorFn: (state: T) => U): (state: T) => U
useStoreWithEqualityFn<T, U = T>(store: StoreApi<T>, selectorFn: (state: T) => U, equalityFn?: (a: T, b: T) => boolean): U
```

## Middlewares

```ts
combine<T, U>(initialState: T, additionalStateCreatorFn: StateCreator<T, [], [], U>): StateCreator<Omit<T, keyof U> & U, [], []>
devtools<T>(stateCreatorFn: StateCreator<T, [], []>, devtoolsOptions?: DevtoolsOptions): StateCreator<T, [['zustand/devtools', never]], []>
immer<T>(stateCreatorFn: StateCreator<T, [], []>): StateCreator<T, [['zustand/immer', never]], []>
persist<T, U>(stateCreatorFn: StateCreator<T, [], []>, persistOptions?: PersistOptions<T, U>): StateCreator<T, [['zustand/persist', U]], []>
redux<T, A>(reducerFn: (state: T, action: A) => T, initialState: T): StateCreator<T & { dispatch: (action: A) => A }, [['zustand/redux', A]], []>

```

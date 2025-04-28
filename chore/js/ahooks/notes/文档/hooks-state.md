## useSetState

通过回调进行更新，可以获取上一次的状态，并且也会自动合并返回的对象。

```ts
export type SetState<S extends Record<string, any>> = <K extends keyof S>(
  state: Pick<S, K> | null | ((prevState: Readonly<S>) => Pick<S, K> | S | null)
) => void
declare const useSetState: <S extends Record<string, any>>(initialState: S | (() => S)) => [S, SetState<S>]
```

## useBoolean

```ts
export interface Actions {
  setTrue: () => void
  setFalse: () => void
  set: (value: boolean) => void
  toggle: () => void
}
export default function useBoolean(defaultValue?: boolean): [boolean, Actions]
```

## useToggle

```ts
export interface Actions<T> {
  setLeft: () => void
  setRight: () => void // 如果传入了 reverseValue, 则设置为 reverseValue。 否则设置为 defaultValue 的反值
  set: (value: T) => void
  toggle: () => void
}
declare function useToggle<T = boolean>(): [boolean, Actions<T>]
declare function useToggle<T>(defaultValue: T): [T, Actions<T>]
declare function useToggle<T, U>(defaultValue: T, reverseValue: U): [T | U, Actions<T | U>]
```

## useUrlState

`通过 url query 来管理 state 的 Hook。`
`将状态同步到 url query 中`。通过设置值为 undefined, 可以从 url query 上彻底删除某个属性，从而使用默认值。
**useUrlState 可以同时管理多个状态**。

npm i @ahooksjs/use-url-state -S
该 Hooks 基于 react-router 的 useLocation & useHistory & useNavigate 进行 query 管理，所以使用该 Hooks 之前，你需要保证

1. 你项目正在使用 react-router 5.x 或 6.x 版本来管理路由
2. 独立安装了 @ahooksjs/use-url-state

```ts
export interface Options {
  navigateMode?: 'push' | 'replace'
  parseOptions?: ParseOptions
  stringifyOptions?: StringifyOptions
}
declare type UrlState = Record<string, any>
declare const useUrlState: <S extends UrlState = UrlState>(
  initialState?: S | (() => S) | undefined,
  options?: Options // 可以通过传入 parseOptions 和 stringifyOptions 自定义转换规则。(来自 query-string)
) => readonly [Partial<{ [key in keyof S]: any }>, (s: React.SetStateAction<Partial<{ [key in keyof S]: any }>>) => void]
```

```ts
import React from 'react'
import useUrlState from '@ahooksjs/use-url-state'

export default () => {
  const [state, setState] = useUrlState(
    { ids: ['1', '2', '3'] },
    {
      parseOptions: {
        arrayFormat: 'comma'
      },
      stringifyOptions: {
        arrayFormat: 'comma'
      }
    }
  )

  return (
    <div>
      <button
        onClick={() => {
          const arr = Array(3)
            .fill(1)
            .map(() => Math.floor(Math.random() * 10))
          setState({ ids: arr })
        }}
      >
        变更数组state
      </button>
      <div>ids: {JSON.stringify(state.ids)}</div>
    </div>
  )
}
```

## useCookieState

一个可以将状态存储在 Cookie 中的 Hook 。
`刷新页面后，可以看到输入框中的内容被从 Cookie 中恢复了。`

可配置属性：默认值、有效时间、路径、域名、协议、跨域等，详见 Options 文档
如果想从 document.cookie 中删除这条数据，可以使用 setState() 或 setState(undefined)

```ts
type State = string | undefined;

type SetState = (
  newValue?: State | ((prevState?: State) => State),
  options?: Cookies.CookieAttributes,
) => void;

const [state, setState]: [State, SetState] = useCookieState(
  cookieKey: string,
  options?: Options,
);
```

## useLocalStorageState

- 自定义序列化和反序列化函数
  useLocalStorageState 在往 localStorage 写入数据前，会先调用一次 serializer，在读取数据之后，会先调用一次 deserializer
- 将 state 与 localStorage 保持同步
  存储值变化时，**所有 key 相同的 useLocalStorageState 会同步状态，包括同一浏览器不同 tab 之间**（尝试打开两个此页面，点击其中一个页面的按钮，另一个页面的 count 会自动更新）

```ts
export declare const SYNC_STORAGE_EVENT_NAME = 'AHOOKS_SYNC_STORAGE_EVENT_NAME'
export type SetState<S> = S | ((prevState?: S) => S)
export interface Options<T> {
  defaultValue?: T | (() => T)
  listenStorageChange?: boolean // 是否同步所有 tab 的状态
  serializer?: (value: T) => string
  deserializer?: (value: string) => T
  onError?: (error: unknown) => void
}
export declare function createUseStorageState(
  getStorage: () => Storage | undefined
): <T>(key: string, options?: Options<T>) => readonly [T | undefined, (value?: SetState<T>) => void]

declare const useLocalStorageState: <T>(
  key: string,
  options?: import('../createUseStorageState').Options<T>
) => readonly [T | undefined, (value?: import('../createUseStorageState').SetState<T>) => void]
```

## useSessionStorageState

用法与 useLocalStorageState 一致。

## useDebounce

```ts
export interface DebounceOptions {
  wait?: number // 默认1000
  maxWait?: number
  leading?: boolean // 默认false
  trailing?: boolean // 默认true
}

declare function useDebounce<T>(value: T, options?: DebounceOptions): T

import React, { useState } from 'react'

export default () => {
  const [value, setValue] = useState<string>()
  const debouncedValue = useDebounce(value, { wait: 500 })

  return (
    <div>
      <input value={value} onChange={e => setValue(e.target.value)} placeholder="Typed value" style={{ width: 280 }} />
      <p style={{ marginTop: 16 }}>DebouncedValue: {debouncedValue}</p>
    </div>
  )
}
```

## useThrottle

```ts
export interface ThrottleOptions {
  wait?: number // 默认1000
  leading?: boolean // 默认true
  trailing?: boolean // 默认true
}

declare function useThrottle<T>(value: T, options?: ThrottleOptions): T
```

## useMap

```ts
declare function useMap<K, T>(
  initialValue?: Iterable<readonly [K, T]>
): readonly [
  Map<K, T>,
  {
    readonly set: (key: K, entry: T) => void
    readonly setAll: (newMap: Iterable<readonly [K, T]>) => void
    readonly remove: (key: K) => void
    readonly reset: () => void
    readonly get: (key: K) => T | undefined
  }
]
```

## useSet

```ts
declare function useSet<K>(initialValue?: Iterable<K>): readonly [
  Set<K>,
  {
    readonly add: (key: K) => void
    readonly remove: (key: K) => void
    readonly reset: () => void
  }
]
```

## usePrevious

保存上一次状态的 Hook。
只有 shouldUpdate function 返回 true 时，才会记录值的变化。
默认 `(a, b) => !Object.is(a, b)`.

```ts
export type ShouldUpdateFunc<T> = (prev: T | undefined, next: T) => boolean
declare function usePrevious<T>(state: T, shouldUpdate?: ShouldUpdateFunc<T>): T | undefined
```

## useRafState

**只在 requestAnimationFrame callback 时更新 state**，一般用于性能优化。
同 useState 用法.

```ts
import { useRafState } from 'ahooks'
import React, { useEffect } from 'react'

export default () => {
  const [state, setState] = useRafState({
    width: 0,
    height: 0
  })

  useEffect(() => {
    const onResize = () => {
      setState({
        width: document.documentElement.clientWidth,
        height: document.documentElement.clientHeight
      })
    }
    onResize()

    window.addEventListener('resize', onResize)

    return () => {
      window.removeEventListener('resize', onResize)
    }
  }, [])

  return (
    <div>
      <p>Try to resize the window </p>
      current: {JSON.stringify(state)}
    </div>
  )
}
```

## useSafeState

用法与 React.useState 完全一样，但是在**组件卸载后异步回调内的 setState 不再执行，避免因组件卸载后更新状态而导致的内存泄漏。**

```ts
import { useSafeState } from 'ahooks'
import React, { useEffect, useState } from 'react'

const Child = () => {
  const [value, setValue] = useSafeState<string>()

  useEffect(() => {
    setTimeout(() => {
      setValue('data loaded from server')
    }, 5000)
  }, [])

  const text = value || 'Loading...'

  return <div>{text}</div>
}

export default () => {
  const [visible, setVisible] = useState(true)
  return (
    <div>
      <button onClick={() => setVisible(false)}>Unmount</button>
      {visible && <Child />}
    </div>
  )
}
```

## useGetState

给 React.useState 增加了一个 getter 方法，**以获取当前最新值**。
**避免闭包问题的比较好的解法。**

```ts
import React, { useEffect } from 'react'
import { useGetState } from 'ahooks'

export default () => {
  const [count, setCount, getCount] = useGetState<number>(0)

  useEffect(() => {
    const interval = setInterval(() => {
      console.log('interval count', getCount())
    }, 3000)

    return () => {
      clearInterval(interval)
    }
  }, [])

  return <button onClick={() => setCount(count => count + 1)}>count: {count}</button>
}
```

```ts
type GetStateAction<S> = () => S
declare function useGetState<S>(initialState: S | (() => S)): [S, Dispatch<SetStateAction<S>>, GetStateAction<S>]
declare function useGetState<S = undefined>(): [S | undefined, Dispatch<SetStateAction<S | undefined>>, GetStateAction<S | undefined>]
```

## useResetState

提供重置 state 方法的 Hooks，用法与 React.useState 基本一致。

```ts
type ResetState = () => void
declare const useResetState: <S>(initialState: S | (() => S)) => [S, Dispatch<SetStateAction<S>>, ResetState]
```

---

## react-router 的 useLocation & useHistory & useNavigate

`useLocation`、`useHistory`、`useNavigate` 都是 React Router（react-router-dom）提供的**路由相关的 Hooks**，用于在函数组件中获取和操作路由信息。

---

### 1. `useLocation`

- **作用**：获取当前路由的 location 信息（如 pathname、search、hash 等）。
- **常用场景**：读取当前 URL、监听路由变化。
- **示例**：
  ```js
  import { useLocation } from 'react-router-dom'
  const location = useLocation()
  console.log(location.pathname) // 当前路径
  ```

---

### 2. `useHistory`（仅 5.x 版本）

- **作用**：操作浏览器历史记录（前进、后退、跳转）。
- **常用场景**：编程式跳转页面。
- **示例**：
  ```js
  import { useHistory } from 'react-router-dom'
  const history = useHistory()
  history.push('/home') // 跳转到 /home
  ```

---

### 3. `useNavigate`（6.x 版本替代 useHistory）

- **作用**：在 React Router v6 中，用于跳转路由。
- **常用场景**：编程式跳转页面。
- **示例**：
  ```js
  import { useNavigate } from 'react-router-dom'
  const navigate = useNavigate()
  navigate('/home') // 跳转到 /home
  ```

---

#### 总结记忆

- `useLocation`：获取当前路由信息
- `useHistory`：5.x 路由跳转
- `useNavigate`：6.x 路由跳转（推荐用）

**一句话理解：它们让你在函数组件里获取和操作路由信息。**

## leading 和 trailing

### 场景一：只想在用户**开始输入**时立即响应（`leading: true, trailing: false`）

**例子：**
你有一个搜索输入框，希望用户一输入就立刻显示“正在搜索...”，但不需要在输入结束后再触发一次。

```js
const debouncedValue = useDebounce(value, { wait: 500, leading: true, trailing: false })
```

**效果：**

- 用户一敲键盘就立即触发一次（比如显示 loading），后续输入不再触发，直到下次输入。

---

### 场景二：只想在用户**停止输入**后再响应（`leading: false, trailing: true`，默认）

**例子：**
你有一个搜索输入框，希望用户输入完停下来 500ms 后再去请求接口，避免频繁请求。

```js
const debouncedValue = useDebounce(value, { wait: 500, leading: false, trailing: true })
```

**效果：**

- 用户输入时不会触发，只有停下来 500ms 后才触发一次（常用于防抖搜索）。

---

### 场景三：**开始和结束都要响应**（`leading: true, trailing: true`）

**例子：**
你既想在用户刚输入时立即响应一次，也想在用户停下来后再响应一次。

```js
const debouncedValue = useDebounce(value, { wait: 500, leading: true, trailing: true })
```

**效果：**

- 用户一输入就触发一次，输入结束后再触发一次。

---

**总结口诀：**

- `leading` 适合“刚开始就要反馈”
- `trailing` 适合“用户停下来再处理”

根据你的业务需求选择即可。

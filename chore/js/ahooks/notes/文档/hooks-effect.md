## useUpdateEffect

useUpdateEffect 用法等同于 useEffect，但是会忽略首次执行，只在依赖更新时执行。

```ts
declare const _default: typeof useEffect | typeof import('react').useLayoutEffect
```

## useUpdateLayoutEffect

## useAsyncEffect

useEffect 支持异步函数。
组件加载时进行`异步的检查`、`中断执行`。

```tsx
import React, { useState } from 'react'

declare function useAsyncEffect(effect: () => AsyncGenerator<void, void, void> | Promise<void>, deps?: DependencyList): void

function mockCheck(): Promise<boolean> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(true)
    }, 3000)
  })
}

export default () => {
  const [pass, setPass] = useState<boolean>()

  useAsyncEffect(async () => {
    setPass(await mockCheck())
  }, [])

  return (
    <div>
      {pass === undefined && 'Checking...'}
      {pass === true && 'Check passed.'}
    </div>
  )
}
```

```tsx
import React, { useState } from 'react'
import { useAsyncEffect } from 'ahooks'

function mockCheck(val: string): Promise<boolean> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(val.length > 0)
    }, 1000)
  })
}

export default () => {
  const [value, setValue] = useState('')
  const [pass, setPass] = useState<boolean>()

  useAsyncEffect(
    async function* () {
      setPass(undefined)
      const result = await mockCheck(value)
      yield // 检查当前副作用（effect）是否还有效，如果已经被清理（比如依赖变化或组件卸载），就停止执行后面的代码。
      setPass(result)
    },
    [value]
  )

  return (
    <div>
      <input
        value={value}
        onChange={e => {
          setValue(e.target.value)
        }}
      />
      <p>
        {pass === null && 'Checking...'}
        {pass === false && 'Check failed.'}
        {pass === true && 'Check passed.'}
      </p>
    </div>
  )
}
```

## useDebounceEffect

为 useEffect 增加防抖的能力。

```tsx
import { useDebounceEffect } from 'ahooks'
import React, { useState } from 'react'

export default () => {
  const [value, setValue] = useState('hello')
  const [records, setRecords] = useState<string[]>([])
  useDebounceEffect(
    () => {
      setRecords(val => [...val, value])
    },
    [value],
    {
      wait: 1000
    }
  )

  return (
    <div>
      <input value={value} onChange={e => setValue(e.target.value)} placeholder="Typed value" style={{ width: 280 }} />
      <p style={{ marginTop: 16 }}>
        <ul>
          {records.map((record, index) => (
            <li key={index}>{record}</li>
          ))}
        </ul>
      </p>
    </div>
  )
}
```

## useDebounceFn

用来处理防抖函数的 Hook。
run 触发执行，cancel 取消执行，flush 立即执行当前的防抖函数。

```ts
type noop = (...args: any[]) => any
declare function useDebounceFn<T extends noop>(
  fn: T,
  options?: DebounceOptions
): {
  run: import('lodash').DebouncedFunc<(...args: Parameters<T>) => ReturnType<T>>
  cancel: () => void
  flush: () => ReturnType<T> | undefined
}
```

```tsx
import { useDebounceFn } from 'ahooks'
import React, { useState } from 'react'

export default () => {
  const [value, setValue] = useState(0)
  const { run } = useDebounceFn(
    () => {
      setValue(value + 1)
    },
    {
      wait: 500
    }
  )

  return (
    <div>
      <p style={{ marginTop: 16 }}> Clicked count: {value} </p>
      <button type="button" onClick={run}>
        Click fast!
      </button>
    </div>
  )
}
```

## useThrottleFn

同上。

## useThrottleEffect

同上.

## useDeepCompareEffect

用法与 useEffect 一致，但 deps 通过 `react-fast-compare` 进行深比较。

```ts
declare const _default: typeof useEffect | typeof import('react').useLayoutEffect
```

```tsx
import { useDeepCompareEffect } from 'ahooks'
import React, { useEffect, useState, useRef } from 'react'

export default () => {
  const [_, setCount] = useState(0)
  const effectCountRef = useRef(0)
  const deepCompareCountRef = useRef(0)

  useEffect(() => {
    effectCountRef.current += 1
  }, [{}]) // 每次渲染都会执行

  useDeepCompareEffect(() => {
    deepCompareCountRef.current += 1
    return () => {
      // do something
    }
  }, [{}]) // 只有在深比较不相等时才会执行

  return (
    <div>
      <p>effectCount: {effectCountRef.current}</p>
      <p>deepCompareCount: {deepCompareCountRef.current}</p>
      <p>
        <button type="button" onClick={() => setCount(c => c + 1)}>
          reRender
        </button>
      </p>
    </div>
  )
}
```

## useDeepCompareLayoutEffect

```ts
declare const _default: typeof import('react').useEffect | typeof useLayoutEffect
```

## useInterval

一个可以处理 setInterval 的 Hook。
支持动态修改 delay 以实现定时器间隔变化与暂停。

```ts
declare const useInterval: (
  fn: () => void,
  delay?: number, // 间隔时间，当设置值为 undefined 时会停止计时器
  options?: {
    immediate?: boolean // 是否立即执行一次
  }
) => () => void
```

```tsx
import React, { useState } from 'react'
import { useInterval } from 'ahooks'

export default () => {
  const [count, setCount] = useState(0)
  const [interval, setInterval] = useState<number | undefined>(1000)

  const clear = useInterval(() => {
    setCount(count + 1)
  }, interval)

  return (
    <div>
      <p> count: {count} </p>
      <p style={{ marginTop: 16 }}> interval: {interval} </p>
      <button onClick={() => setInterval(t => (!!t ? t + 1000 : 1000))} style={{ marginRight: 8 }}>
        interval + 1000
      </button>
      <button
        style={{ marginRight: 8 }}
        onClick={() => {
          setInterval(1000)
        }}
      >
        reset interval
      </button>
      <button onClick={clear}>clear</button>
    </div>
  )
}
```

## useRafInterval

用 requestAnimationFrame 模拟实现 setInterval，API 和 useInterval 保持一致。
好处是**可以在页面不渲染的时候停止执行定时器，比如页面隐藏或最小化等**。

请注意，如下两种情况下很可能是不适用的，优先考虑 useInterval ：

- 时间间隔小于 16ms
- 希望页面不渲染的情况下依然执行定时器

Node 环境下 requestAnimationFrame 会自动降级到 setInterval

## useTimeout

```ts
declare const useTimeout: (fn: () => void, delay?: number) => () => void
```

## useRafTimeout

同上。

## useLockFn

`用于给一个异步函数增加竞态锁，防止并发执行。`

```ts
// 注意这里，Promise 的返回结果多了一个 undefined.
declare function useLockFn<P extends any[] = any[], V = any>(fn: (...args: P) => Promise<V>): (...args: P) => Promise<V | undefined>
```

```tsx
import { useLockFn } from 'ahooks'
import { message } from 'antd'
import React, { useState } from 'react'

function mockApiRequest() {
  return new Promise<void>(resolve => {
    setTimeout(() => {
      resolve()
    }, 2000)
  })
}

export default () => {
  const [count, setCount] = useState(0)

  const submit = useLockFn(async () => {
    message.info('Start to submit')
    await mockApiRequest()
    setCount(val => val + 1)
    message.success('Submit finished')
  })

  return (
    <>
      <p>Submit count: {count}</p>
      <button onClick={submit}>Submit</button>
    </>
  )
}
```

## useUpdate

useUpdate 会返回一个函数，调用该函数会强制组件重新渲染。

```ts
declare const useUpdate: () => () => void
```

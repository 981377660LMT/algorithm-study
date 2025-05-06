## useTrackedEffect

追踪是`哪个依赖变化`触发了 useEffect 的执行。

```ts
type Effect<T extends DependencyList> = (
  changes?: number[], // 变化的依赖项索引
  previousDeps?: T, // 上一次的依赖项
  currentDeps?: T // 当前的依赖项
) => void | (() => void)
declare const useTrackedEffect: <T extends DependencyList>(effect: Effect<T>, deps?: [...T]) => void
```

```tsx
import React, { useState } from 'react'
import { useTrackedEffect } from 'ahooks'

export default () => {
  const [count, setCount] = useState(0)
  const [count2, setCount2] = useState(0)

  useTrackedEffect(
    changes => {
      console.log('Index of changed dependencies: ', changes)
    },
    [count, count2]
  )

  return (
    <div>
      <p>Please open the browser console to view the output!</p>
      <div>
        <p>Count: {count}</p>
        <button onClick={() => setCount(c => c + 1)}>count + 1</button>
      </div>
      <div style={{ marginTop: 16 }}>
        <p>Count2: {count2}</p>
        <button onClick={() => setCount2(c => c + 1)}>count + 1</button>
      </div>
    </div>
  )
}
```

## useWhyDidYouUpdate

帮助开发者排查是哪个属性改变导致了组件的 rerender。

```ts
export type IProps = Record<string, any>
export default function useWhyDidYouUpdate(componentName: string, props: IProps): void
```

```tsx
import { useWhyDidYouUpdate } from 'ahooks'
import React, { useState } from 'react'

const Demo: React.FC<{ count: number }> = props => {
  const [randomNum, setRandomNum] = useState(Math.random())

  useWhyDidYouUpdate('useWhyDidYouUpdateComponent', { ...props, randomNum })

  return (
    <div>
      <div>
        <span>number: {props.count}</span>
      </div>
      <div>
        randomNum: {randomNum}
        <button onClick={() => setRandomNum(Math.random)} style={{ marginLeft: 8 }}>
          🎲
        </button>
      </div>
    </div>
  )
}

export default () => {
  const [count, setCount] = useState(0)

  return (
    <div>
      <Demo count={count} />
      <div>
        <button onClick={() => setCount(prevCount => prevCount - 1)}>count -</button>
        <button onClick={() => setCount(prevCount => prevCount + 1)} style={{ marginLeft: 8 }}>
          count +
        </button>
      </div>
      <p style={{ marginTop: 8 }}>Please open the browser console to view the output!</p>
    </div>
  )
}
```

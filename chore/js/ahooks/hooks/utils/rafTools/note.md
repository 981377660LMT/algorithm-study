除了你已经实现的 `useRafThrottle`（利用 rAF 进行节流），`requestAnimationFrame` (rAF) 在前端开发中还有很多重要的应用场景，主要集中在**高性能动画**、**高频状态更新**和**定时器优化**上。

以下是几个核心的 rAF 相关 Hooks 及其原理讲解：

### 1. `useRafState` (高性能状态更新)

**场景**：当你需要根据 `scroll`、`resize` 或 `mousemove` 事件更新 React state 时，直接 `setState` 会导致极其频繁的重渲染，造成页面卡顿。
**原理**：`useRafState` 保证在每一帧（约 16ms）内只执行一次 `setState`，丢弃掉同一帧内多余的更新请求。

```typescript
import { useRef, useState, useCallback, useEffect } from 'react'

export function useRafState<S>(initialState: S | (() => S)) {
  const [state, setState] = useState(initialState)
  const frame = useRef(0)

  const setRafState = useCallback((value: S | ((prevState: S) => S)) => {
    cancelAnimationFrame(frame.current)

    frame.current = requestAnimationFrame(() => {
      setState(value)
    })
  }, [])

  useEffect(() => {
    return () => cancelAnimationFrame(frame.current)
  }, [])

  return [state, setRafState] as const
}
```

### 2. `useRafLoop` (动画循环)

**场景**：当你需要制作游戏、Canvas 渲染或复杂的 JS 动画时，需要一个持续运行的循环。
**原理**：封装了递归调用 `requestAnimationFrame` 的逻辑，并提供 `start`、`stop` 控制能力。它比 `setInterval` 更精准，且在 Tab 不可见时会自动暂停以省电。

```typescript
import { useCallback, useEffect, useRef } from 'react'

export function useRafLoop(callback: (time: number) => void) {
  const rafRef = useRef<number>()
  const callbackRef = useRef(callback)

  // 记住最新的 callback
  useEffect(() => {
    callbackRef.current = callback
  })

  const loop = useCallback((time: number) => {
    callbackRef.current(time)
    rafRef.current = requestAnimationFrame(loop)
  }, [])

  const start = useCallback(() => {
    if (!rafRef.current) {
      rafRef.current = requestAnimationFrame(loop)
    }
  }, [loop])

  const stop = useCallback(() => {
    if (rafRef.current) {
      cancelAnimationFrame(rafRef.current)
      rafRef.current = undefined
    }
  }, [])

  // 组件卸载时自动停止
  useEffect(() => stop, [stop])

  return { start, stop }
}
```

### 3. `useRafInterval` (更可靠的定时器)

**场景**：`setInterval` 在主线程阻塞时会推迟执行，或者在后台 Tab 中被浏览器强行降频（甚至停止）。
**原理**：使用 `requestAnimationFrame` 模拟 `setInterval`。虽然 rAF 也会在后台暂停，但它能保证回调在页面渲染前执行，且不会像 `setInterval` 那样在某些情况下出现“回调堆积”后的一瞬间爆发执行。

```typescript
import { useEffect, useRef } from 'react'

export function useRafInterval(callback: () => void, delay: number | undefined) {
  const callbackRef = useRef(callback)
  const timerRef = useRef<number>()
  const startTimeRef = useRef<number>()

  useEffect(() => {
    callbackRef.current = callback
  })

  useEffect(() => {
    if (typeof delay !== 'number' || delay < 0) return

    const loop = (time: number) => {
      if (!startTimeRef.current) {
        startTimeRef.current = time
      }

      const elapsed = time - startTimeRef.current

      if (elapsed >= delay) {
        callbackRef.current()
        startTimeRef.current = time // 重置起始时间
      }

      timerRef.current = requestAnimationFrame(loop)
    }

    timerRef.current = requestAnimationFrame(loop)

    return () => {
      if (timerRef.current) {
        cancelAnimationFrame(timerRef.current)
      }
    }
  }, [delay])
}
```

### 4. `useRafTimeout` (防止内存泄漏的延时)

**场景**：类似于 `setTimeout`，但用于处理与 UI 渲染紧密相关的延时操作。
**原理**：确保回调在浏览器准备好绘制下一帧时执行。如果组件在倒计时结束前卸载，rAF 会更容易清理（虽然 `setTimeout` 也可以清理，但 rAF 在处理复杂动画序列时心智模型更统一）。

```typescript
import { useEffect, useRef } from 'react'

export function useRafTimeout(callback: () => void, delay: number | undefined) {
  const callbackRef = useRef(callback)
  const timerRef = useRef<number>()
  const startTimeRef = useRef<number>()

  useEffect(() => {
    callbackRef.current = callback
  })

  useEffect(() => {
    if (typeof delay !== 'number' || delay < 0) return

    const loop = (time: number) => {
      if (!startTimeRef.current) {
        startTimeRef.current = time
      }

      const elapsed = time - startTimeRef.current

      if (elapsed >= delay) {
        callbackRef.current()
      } else {
        timerRef.current = requestAnimationFrame(loop)
      }
    }

    timerRef.current = requestAnimationFrame(loop)

    return () => {
      if (timerRef.current) {
        cancelAnimationFrame(timerRef.current)
      }
    }
  }, [delay])
}
```

### 总结与对比

| Hook 名称            | 核心用途     | 解决的问题                                 | 典型场景                          |
| :------------------- | :----------- | :----------------------------------------- | :-------------------------------- |
| **`useRafThrottle`** | **节流执行** | 防止高频事件触发导致函数执行过多           | 滚动监听、窗口 Resize 计算        |
| **`useRafState`**    | **状态节流** | 防止高频 `setState` 导致 React 渲染卡顿    | 拖拽跟随、实时图表更新            |
| **`useRafLoop`**     | **动画循环** | 提供可控的、高性能的渲染循环               | Canvas 游戏、粒子效果、自定义动画 |
| **`useRafInterval`** | **精准轮询** | 替代 `setInterval`，避免后台堆积和执行漂移 | 倒计时、轮询接口（需页面可见时）  |

这些工具的核心思想都是：**将 JS 的执行时机与浏览器的刷新频率（60Hz）对齐**，从而获得最流畅的视觉体验并节省性能。

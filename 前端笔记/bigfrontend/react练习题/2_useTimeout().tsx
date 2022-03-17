// Create a hook to easily use setTimeout(callback, delay).

// reset the timer if delay changes
// DO NOT reset the timer if only callback changes
import React from 'react'

export function useTimeout(callback: () => void, delay: number) {
  // your code here
  const timerRef = React.useRef<number>(0)
  const callbackRef = React.useRef<() => void>(callback)

  React.useEffect(() => {
    callbackRef.current = callback
  }, [callback])

  React.useEffect(() => {
    timerRef.current = window.setTimeout(() => callbackRef.current(), delay)
    return () => {
      window.clearTimeout(timerRef.current)
    }
  }, [delay])
}

// if you want to try your code on the right panel
// remember to export App() component like below

export function App() {
  return <div>your app</div>
}

// NOT reset timeout when only callback changes
// 需要使用ref来保存回调函数
// React.useCallback能减少组件重复更新时带来的性能消耗，无法做到保存原来函数

// loadingDelay 主要用于防止 loading 状态的闪烁.

import { useState, useRef, useCallback, useEffect } from 'react'

interface IOptions {
  /**
   * 设置为 true，则需要手动调用 run 来触发执行。
   * @default false
   */
  manual?: boolean
  /**
   * 延迟显示 loading 状态的时间，单位为毫秒。
   * 如果异步请求在 `loadingDelay` 内完成，则 loading 状态始终为 false。
   * @default 0
   */
  loadingDelay?: number
}

/**
 * 一个管理异步请求状态的 Hook，支持延迟显示 loading 状态以防止闪烁。
 * @param asyncFn 一个返回 Promise 的异步函数。
 * @param options 配置项。
 * @returns 返回一个包含 loading 状态和 run 执行函数的对象。
 */
export function useDelayedRequest<T extends (...args: any[]) => Promise<any>>(
  asyncFn: T,
  options: IOptions = {}
) {
  const { manual = false, loadingDelay = 0 } = options

  const [loading, setLoading] = useState(false)
  const timerRef = useRef<NodeJS.Timeout | null>(null)

  // 使用 useRef 存储最新的 asyncFn，避免 useCallback 依赖频繁变化
  const asyncFnRef = useRef(asyncFn)
  asyncFnRef.current = asyncFn

  // 组件卸载时，清理所有定时器
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current)
      }
    }
  }, [])

  const run = useCallback(
    async (...args: Parameters<T>): Promise<ReturnType<T> | undefined> => {
      // 1. 设置一个定时器，在 loadingDelay 之后才将 loading 设为 true
      if (loadingDelay > 0) {
        timerRef.current = setTimeout(() => {
          setLoading(true)
        }, loadingDelay)
      } else {
        setLoading(true)
      }

      try {
        // 2. 执行真正的异步函数
        return await asyncFnRef.current(...args)
      } finally {
        // 3. 异步函数执行完毕后（无论成功或失败）
        // a. 清除可能还未执行的定时器
        if (timerRef.current) {
          clearTimeout(timerRef.current)
        }
        // b. 将 loading 状态恢复为 false
        setLoading(false)
      }
    },
    [loadingDelay]
  )

  // 如果不是手动模式，则在组件首次渲染时自动执行
  useEffect(() => {
    if (!manual) {
      run()
    }
  }, [manual, run])

  return { loading, run }
}

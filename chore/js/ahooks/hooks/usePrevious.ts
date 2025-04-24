import { useRef } from 'react'

export type ShouldUpdateFunc<T> = (prev: T | undefined, next: T) => boolean

const defaultShouldUpdate = <T>(a?: T, b?: T) => !Object.is(a, b)

/**
 * 保存上一次状态的 Hook.
 *
 * https://ahooks.js.org/zh-CN/hooks/use-previous
 */
export function usePrevious<T>(state: T, shouldUpdate: ShouldUpdateFunc<T> = defaultShouldUpdate): T | undefined {
  const prevRef = useRef<T>()
  const curRef = useRef<T>()

  if (shouldUpdate(curRef.current, state)) {
    prevRef.current = curRef.current
    curRef.current = state
  }

  return prevRef.current
}

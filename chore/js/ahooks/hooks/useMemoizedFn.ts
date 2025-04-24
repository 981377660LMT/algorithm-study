import { useMemo, useRef } from 'react'

type Noop = (this: any, ...args: any[]) => any

type PickFunction<T extends Noop> = (this: ThisParameterType<T>, ...args: Parameters<T>) => ReturnType<T>

/**
 * 使用 useMemoizedFn，可以省略第二个参数 deps，同时保证函数地址永远不会变化.
 *
 * https://ahooks.js.org/zh-CN/hooks/use-memoized-fn
 */
export function useMemoizedFn<T extends Noop>(fn: T): T {
  const fnRef = useRef<T>(fn)

  // why not write `fnRef.current = fn`?
  // https://github.com/alibaba/hooks/issues/728
  fnRef.current = useMemo(() => fn, [fn])

  const memoizedFn = useRef<PickFunction<T>>()
  if (!memoizedFn.current) {
    memoizedFn.current = function (this, ...args) {
      return fnRef.current.apply(this, args)
    }
  }

  return memoizedFn.current as T
}

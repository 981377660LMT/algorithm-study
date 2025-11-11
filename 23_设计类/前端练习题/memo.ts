/**
 * 函数的记忆化版本，类似于 `React.memo`.
 */
export function memo<T extends (...args: any[]) => any>(
  fn: T,
  areArgsEqual?: (prevArgs: Parameters<T>, nextArgs: Parameters<T>) => boolean
): T {
  let prevArgs: Parameters<T> | null = null
  let preResult: ReturnType<T> | null = null
  let called = false
  const compare = areArgsEqual || shallowCompareArgs
  return function (this: any, ...nextArgs: Parameters<T>): ReturnType<T> {
    if (called && compare(prevArgs!, nextArgs)) return preResult!
    preResult = fn.apply(this, nextArgs)
    prevArgs = nextArgs
    called = true
    return preResult!
  } as T
}

function shallowCompareArgs(prevArgs: readonly any[], nextArgs: readonly any[]): boolean {
  if (prevArgs.length !== nextArgs.length) return false
  for (let i = 0; i < prevArgs.length; i++) {
    if (prevArgs[i] !== nextArgs[i]) return false
  }
  return true
}

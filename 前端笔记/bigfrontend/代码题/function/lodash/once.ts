/**
 * 确保函数只执行一次的通用工具函数
 * @param fn 需要包装的函数
 * @returns 包装后的函数
 */
export function once<T extends (...args: any[]) => any>(fn: T): T {
  let called = false
  let result: ReturnType<T>
  return function (this: any, ...args: Parameters<T>): ReturnType<T> {
    if (!called) {
      called = true
      result = fn.apply(this, args)
    }
    return result
  } as T
}

// 使用示例：

// 1. 同步函数
const syncFn = once((name: string) => {
  console.log('执行同步函数')
  return `Hello, ${name}`
})

console.log(syncFn('Alice')) // 输出: 执行同步函数, Hello, Alice
console.log(syncFn('Bob')) // 输出: Hello, Alice (不再执行逻辑)

// 2. 异步函数
const asyncFn = once(async (id: number) => {
  console.log('执行异步请求')
  return Promise.resolve({ id, data: 'success' })
})

;(async () => {
  console.log(await asyncFn(1)) // 输出: 执行异步请求, { id: 1, data: 'success' }
  console.log(await asyncFn(2)) // 输出: { id: 1, data: 'success' } (返回缓存的 Promise)
})()

type Callback = (error: Error, data: any) => void

type AsyncFunc = (callback: Callback, data: any) => void

/**
 * @param {AsyncFunc[]} funcs
 * @return {(callback: Callback) => void}
 * 请实现一个async helper - sequence()。sequence()像pipe() 那样将异步函数串联在一起。
 * 能否使用Promise完成题目？能否不使用Promise完成该题目？
 */
function sequence(funcs: AsyncFunc[]): (callback: Callback) => void {
  // your code here
}

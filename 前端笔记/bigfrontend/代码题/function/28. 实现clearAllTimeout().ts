/**
 * cancel all timer from window.setTimeout
 * 你能否实现一个clearAllTimeout() 来取消掉所有未执行的timer？
 * 比如当页面跳转的时候我们或许想要清除掉所有的timer。
 * 你需要保证window.setTimeout 和 window.clearTimeout 还是原来的interface，虽然你可以替换其中的逻辑。
 * @description
 * 思路是注册时用一个装饰器
 * 注册时将timer保存到容器
 * clear时清空容器
 */
function clearAllTimeout() {
  const originSetTimeout = window.setTimeout
  const originClearTimeout = window.clearTimeout
  const timerContainer = new Set<number>()

  window.clearAllTimeout = () => {
    for (const timerId of timerContainer) {
      originClearTimeout(timerId)
    }
  }

  // @ts-ignore
  window.setTimeout = (callback: (...args: any[]) => void, ms?: number, ...args: any[]): number => {
    const callbackWrapper = () => {
      callback(...args)
      timerContainer.delete(timer)
    }
    const timer = originSetTimeout(callbackWrapper, ms)
    timerContainer.add(timer)
    return timer
  }

  // @ts-ignore
  window.clearTimeout = (timer: number) => {
    originClearTimeout(timer)
    timerContainer.delete(timer)
  }
}

// setTimeout(func1, 10000)
// setTimeout(func2, 10000)
// setTimeout(func3, 10000)

// // 3个方法都是设定在10秒以后
// clearAllTimeout()

// // 所有方法的timer都被取消掉了

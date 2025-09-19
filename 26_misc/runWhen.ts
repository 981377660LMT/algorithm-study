// runOnce 的扩展版.
// “只执行一次”的职责完全交给“条件函数” (predicate).
// 在函数运行成功时，修改标识位.

/**
 * 创建一个函数，该函数仅在 predicate.when() 返回 true 时才运行 fn。
 * fn 成功运行后（对于异步函数，则指其 Promise 成功 resolve 后），会调用 predicate.done()。
 *
 * @param fn 要运行的函数，可以是同步或异步的。
 * @param predicate 一个包含 when 和 done 方法的对象。
 *                  when: () => boolean - 决定 fn 是否应该运行。
 *                  done: () => void - 在 fn 成功运行后调用。
 * @returns 一个新的包装函数。
 */
function runWhen<Params, Return>(
  fn: (...args: Params[]) => Return,
  predicate: { when: () => boolean; done: () => void }
): (...args: Params[]) => Return | undefined {
  return function (...args: Params[]): Return | undefined {
    if (!predicate.when()) {
      return undefined
    }

    try {
      const result: any = fn(...args)
      if (result && typeof result.then === 'function') {
        return result
          .then((finalResult: any) => {
            predicate.done()
            return finalResult
          })
          .catch((error: any) => {
            throw error
          })
      } else {
        predicate.done()
        return result
      }
    } catch (error) {
      throw error
    }
  }
}

// --- 使用示例 ---
// 1. 同步函数示例
let hasRun = false
const runOncePredicate = {
  when: () => !hasRun,
  done: () => {
    hasRun = true
    console.log('同步任务状态已更新。')
  }
}

const logOnce = runWhen(() => {
  console.log('这个同步函数只会执行一次。')
  return 'Sync Success'
}, runOncePredicate)

function runSyncExample() {
  console.log('--- 同步示例 ---')
  // 无需 await
  const r1 = logOnce()
  console.log('第一次调用结果:', r1)
  const r2 = logOnce()
  console.log('第二次调用结果:', r2)
  console.log('--------------------')
}

// 2. 异步函数示例
let isTaskDone = false
const asyncPredicate = {
  when: () => !isTaskDone,
  done: () => {
    isTaskDone = true
    console.log('异步任务状态已更新为完成。')
  }
}

const doAsyncTask = runWhen(async (duration: number) => {
  console.log('异步任务开始...')
  await new Promise(resolve => setTimeout(resolve, duration))
  console.log('异步任务完成！')
  return 'Async Success'
}, asyncPredicate)

async function runAsyncExample() {
  console.log('--- 异步示例 ---')
  // 需要 await
  const r1 = await doAsyncTask(500)
  console.log('第一次调用结果:', r1)
  const r2 = await doAsyncTask(500)
  console.log('第二次调用结果:', r2)
}

async function main() {
  runSyncExample()
  await runAsyncExample()
}

main()

export {}

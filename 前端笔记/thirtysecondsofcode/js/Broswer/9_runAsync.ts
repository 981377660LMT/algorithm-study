// 通过使用 webworker 在单独的线程中运行函数，允许长时间运行的函数不会阻塞 UI。
// 1.使用 Blob 对象 URL 创建一个新的 Worker(scriptURL: string | URL) ，其内容应该是所提供函数的字符串化版本。
// 2.返回一个新的 Promise () ，监听 onmessage 和 onerror 事件并解析从工作者返回的数据，或者抛出一个错误。

function runAsync(fn: (...args: any[]) => any): Promise<void> {
  // Worker extends EventTarget, AbstractWorker
  const worker = new Worker(URL.createObjectURL(new Blob([`postMessage((${fn})());`])), {
    // @ts-ignore
    type: 'application/javascript; charset=utf-8',
  })

  return new Promise((resolve, reject) => {
    worker.onmessage = ({ data }) => {
      resolve(data)
      worker.terminate()
    }

    worker.onerror = err => {
      reject(err)
      worker.terminate()
    }
  })
}

const longRunningFunction = () => {
  let result = 0
  for (let i = 0; i < 1000; i++)
    for (let j = 0; j < 700; j++) for (let k = 0; k < 300; k++) result = result + i + j + k

  return result
}

runAsync(longRunningFunction).then(console.log) // 209685000000
runAsync(() => 10 ** 3).then(console.log) // 1000

const outsideVariable = 50

/**
  NOTE: Since the function is running in a different context, closures are not supported.
  The function supplied to `runAsync` gets stringified, so everything becomes literal.
  All variables and functions must be defined inside.
**/
runAsync(() => typeof outsideVariable).then(console.log) // 'undefined'
// 因为 fn 被转成 string 所以闭包不起作用

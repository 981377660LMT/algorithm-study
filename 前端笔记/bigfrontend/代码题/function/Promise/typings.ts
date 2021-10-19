type Callback = (error: Error | null, res: any) => void

/**
 * @example 
 * 回调形式读取文件
 * ```js
 * fs.readFile('input.txt', function (err, data) {
    if (err) return console.error(err);
    console.log(data.toString());
  });
  ```
 */
type AsyncFunc = (arg: any, callback: Callback) => void

type PromiseFunc = (arg: any) => Promise<any>

//////////////////////////////////////////////////////////
interface Executor<T> {
  (resolve: Resolve<T>, reject: Reject): void
}

interface Resolve<T> {
  (value: T): void
}

interface Reject {
  (reason?: any): void
}

type Status = 'pending' | 'onfulfilled' | 'onrejected'
export type { Callback, AsyncFunc, PromiseFunc }
export type { Executor, Resolve, Reject, Status }

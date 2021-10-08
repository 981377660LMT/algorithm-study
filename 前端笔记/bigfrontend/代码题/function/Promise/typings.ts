type Callback = (error: Error | null, res: any) => void

/**
 * @example
 * ```js
 * fs.readFile('input.txt', function (err, data) {
    if (err) return console.error(err);
    console.log(data.toString());
  });
  ```
 */
type AsyncFunc = (arg: any, callback: Callback) => void

type PromiseFunc = (arg: any) => Promise<any>

export type { Callback, AsyncFunc, PromiseFunc }

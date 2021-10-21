type Thunk = (callback: Callback) => void // thunk 需要接收一个callback
type Callback = (error: Error | null, data: any | Callback) => void // 注意data可以接受一个callback

const readFile = (cb: Callback, path: string) => {
  setTimeout(() => {
    if (Math.random() > 0.1) cb(null, `${path}文件内容为${Math.random().toFixed(2)}`)
    else cb(new Error('未知错误'), null)
  }, 100)
}

export type { Thunk, Callback }
export { readFile }

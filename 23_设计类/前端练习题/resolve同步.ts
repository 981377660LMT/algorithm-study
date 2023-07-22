let resolve: (value: string) => void
new Promise<string>(r => {
  resolve = (s: string) => {
    console.log('调用resolve')
    r(s)
  }
}).then(() => console.log('调用then'))

resolve!('')
console.log('同步代码')

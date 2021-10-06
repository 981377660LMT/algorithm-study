// Write a higher-order function combinedFetcher, using callback to get all the fetched data.
const fetch = (arg: any) => `res with ${arg}`

const fetcher = (arg: any, cb: (...args: any[]) => any) => {
  const res = fetch(arg)
  cb(res)
}

// 只准使用fetcher
const combinedFetcher = (...args: any[]): ((...args: any) => any) => {
  let cbStore: (...args: any[]) => any
  const resStore: any[] = []

  const innerCB = (res: any) => {
    resStore.push(res)
    if (resStore.length === args.length) cbStore(resStore)
  }

  // 先把cb保存起来
  // 对每个参数获取结果 用参数个数判断是否符合cb的条件
  return (cb: (...args: any) => any) => {
    cbStore = cb
    for (const arg of args) {
      fetcher(arg, innerCB)
    }
  }
}

// [ 'response with: fruits', 'response with: drinks' ]
// function fetch(arg: any) {
//   return `response with: ${arg}`
// }
// function fetcher(arg: any, cb: { (r: any): void; (arg0: string): void }) {
//   let res = fetch(arg)
//   cb(res)
// }
// function combinedFetcher(...args: string[]) {
//   const len = args.length
//   let cnt = 0
//   const res: any[] = []
//   function innerCB(r: any, cb: (arg0: any[]) => void) {
//     cnt++
//     res.push(r)
//     if (cnt === len) cb(res)
//   }

//   return (cb: any) => args.forEach(arg => fetcher(arg, (r: any) => innerCB(r, cb)))
// }

const fetchFruitsAndDrinks = combinedFetcher('fruits', 'drinks')
fetchFruitsAndDrinks(console.log)

export {}

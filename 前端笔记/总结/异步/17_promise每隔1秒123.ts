// Promise 串行
const nums = [1, 2, 3]
nums.reduce((pre, cur) => pre.then(() => work(cur)), Promise.resolve())

function work(num: number) {
  return new Promise<void>(resolve =>
    setTimeout(() => {
      resolve()
    }, 1000)
  )
}

export {}
console.log(2 ** 31 - 1)

const list = [1, 2, 3]
const square = (num: number) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve(num * num)
    }, 1000)
  })
}

function test() {
  list.forEach(async x => {
    const res = await square(x)
    console.log(res)
  })
}
test()

export {}
// forEach是不能阻塞的，默认是请求并行发起，所以是同时输出1、4、9。

async function test2() {
  for (let x of list) {
    const res = await square(x)
    console.log(res)
  }
}
test2()

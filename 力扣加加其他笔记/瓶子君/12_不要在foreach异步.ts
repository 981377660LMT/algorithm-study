const list = [1, 2, 3]
const square = (num: number) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve(num * num)
    }, 1000)
  })
}

// 1s 后同时打印 1、4、9
// function test1() {
//   list.forEach(async x => {
//     const res = await square(x)
//     console.log(res)
//   })
// }
// 隔1s一次输出
async function test2() {
  for (const x of list) {
    const res = await square(x)
    console.log(res)
  }
}
// test1()
test2()

export {}

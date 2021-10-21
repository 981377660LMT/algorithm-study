// 在async中，如果 await后面的内容是一个异常或者错误的话，会怎样呢？
// 如果在async函数中抛出了错误，则终止错误结果，不会继续向下执行。
async function async1() {
  await async2()
  console.log('async1')
  return 'async1 success'
}
async function async2() {
  return new Promise((resolve, reject) => {
    console.log('async2')
    reject('error')
  })
}
async1().then(res => console.log(res))
// 'async2'
// (node:29844) UnhandledPromiseRejectionWarning: error

// 如果想要使得错误的地方不影响async函数后续的执行的话，可以使用try catch
async function async1() {
  try {
    await Promise.reject('error!!!')
  } catch (e) {
    console.log(e) // 2
  }
  console.log('async1') // 3
  return Promise.resolve('async1 success')
}
async1().then(res => console.log(res)) // 4
console.log('script start') // 1

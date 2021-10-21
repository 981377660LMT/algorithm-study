const async1 = async () => {
  console.log('async1') // 2
  setTimeout(() => {
    console.log('timer1') // 7
  }, 2000)
  await new Promise(resolve => {
    console.log('promise1') // 3
  })
  console.log('async1 end') // async函数中await的new Promise要是没有返回值的话则不执行后面的内容(类似题5.5)
  return 'async1 success'
}
console.log('script start') // 1
async1().then(res => console.log(res))
console.log('script end') // 4
Promise.resolve(1)
  .then(2)
  .then(Promise.resolve(3))
  .catch(4)
  .then(res => console.log(res)) // 5
setTimeout(() => {
  console.log('timer2') // 6
}, 1000)

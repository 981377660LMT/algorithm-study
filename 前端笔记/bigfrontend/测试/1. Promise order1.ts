console.log(1) // 1st

// excutor内部函数会立即执行
const promise1 = new Promise<void>(async resolve => {
  console.log(2) // 2nd
  resolve() // 同步立即执行
  console.log(3) // 3rd
})

console.log(4) // 4th

// queued for next tick
promise1
  .then(() => {
    console.log(5) // 6th
  })
  .then(() => {
    console.log(6) // 7th
  })

console.log(7) // 5th

// queued for next tick + 10
setTimeout(() => {
  console.log(8) // 9th
}, 10)

// queued for next tick + 0
setTimeout(() => {
  console.log(9) // 8th
}, 0)

// 1
// 2
// 3
// 4
// 7
// 5
// 6
// 9
// 8

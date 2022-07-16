console.log(1) // 1st

const promise2 = new Promise<void>(async resolve => {
  console.log(2) // 2nd
  await resolve()
  console.log(3) // 5th
})

console.log(4) // 3rd

// queued for next tick
promise2
  .then(() => {
    console.log(5) // 6th
  })
  .then(() => {
    console.log(6) // 7th
  })

console.log(7) // 4th

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
// 4
// 7
// 3
// 5
// 6
// 9
// 8

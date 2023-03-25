console.log(1)

setTimeout(() => {
  console.log(2)
}, 100)

requestAnimationFrame(() => {
  console.log(3)
})

requestAnimationFrame(() => {
  console.log(4)
  setTimeout(() => {
    console.log(5)
  }, 10)
})

const end = Date.now() + 200
while (Date.now() < end) {}

console.log(6)
// 1 6 3 4 2 5
export {}

console.log(1) // 1

document.body.addEventListener('click', () => {
  console.log(2)
})

// micro-task
Promise.resolve().then(() => {
  console.log(3) // 5
})

// macro-task
setTimeout(() => {
  console.log(4)
}, 0) // 6

console.log(5) // 2

document.body.click() // 3

console.log(6) // 4

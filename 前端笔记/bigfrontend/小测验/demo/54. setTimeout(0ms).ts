// This snippet's result may vary on browsers

setTimeout(() => {
  console.log(2)
}, 2)

setTimeout(() => {
  console.log(1)
}, 1)

setTimeout(() => {
  console.log(0)
}, 0)

// 1 0 2

console.log(1)

const mc = new MessageChannel()

mc.port1.onmessage = () => {
  console.log(2)
}

Promise.resolve().then(() => {
  console.log(3)
})

setTimeout(() => {
  console.log(4)
}, 0)

console.log(5)

mc.port2.postMessage('')

console.log(6)
// 1 5 6 3 2 4

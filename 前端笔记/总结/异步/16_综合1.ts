const first = () =>
  new Promise((resolve, reject) => {
    console.log(3) // 1
    let p = new Promise((resolve, reject) => {
      console.log(7) // 2
      setTimeout(() => {
        console.log(5) // 6
        resolve(6)
        console.log(p) // 7  Promise{<resolved>: 1}
      }, 0)
      resolve(1)
    })
    resolve(2)
    p.then(arg => {
      console.log(arg) // 4
    })
  })

first().then(arg => {
  console.log(arg) // 5
})

console.log(4) // 3

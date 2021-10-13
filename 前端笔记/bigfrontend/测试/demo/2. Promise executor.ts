new Promise((resolve, reject) => {
  resolve(1)
  resolve(2)
  reject('error')
}).then(
  value => {
    console.log(value)
  },
  error => {
    console.log('error')
  }
)
// Once a promise is resolved, it does not resolve again.

;(async () => {
  await Promise.all([]).then(
    value => {
      console.log(value)
    },
    error => {
      console.log(error)
    }
  )

  await Promise.all([1, 2, Promise.resolve(3), Promise.resolve(4)]).then(
    value => {
      console.log(value)
    },
    error => {
      console.log(error)
    }
  )

  await Promise.all([1, 2, Promise.resolve(3), Promise.reject('error')]).then(
    value => {
      console.log(value)
    },
    error => {
      console.log(error)
    }
  )
})()

// []
// [ 1, 2, 3, 4 ]
// error

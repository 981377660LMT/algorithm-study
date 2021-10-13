Promise.resolve(1)
  .then(val => {
    console.log(val)
    return val + 1
  })
  .then(val => {
    console.log(val)
  })
  .then(val => {
    console.log(val)
    return Promise.resolve(3).then(val => {
      console.log(val) // 3
      // return undefined
    })
  })
  .then(val => {
    console.log(val) // undefined
    return Promise.reject(4)
  })
  .catch(val => {
    console.log(val) // 4
    // return undefined
  })
  // @ts-ignore
  .finally(val => {
    console.log(val) // undefined: finally has no arguments
    return 10 // no effect on promise object
  })
  .then(val => {
    console.log(val) // undefined 已经被catch住了
  })

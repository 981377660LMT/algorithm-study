Promise.resolve(1)
  .finally(data => {
    console.log(data) // undefined because .finally() does not receive any argument
    return Promise.reject('error')
  })
  .catch(error => {
    console.log(error) // "error" from the rejected promise in previous .finally()
    throw 'error2'
  })
  .finally(data => {
    console.log(data) // undefined because .finally() does not receive any argument
    return Promise.resolve(2).then(console.log) // 2
  })
  .then(console.log) // ignored
  .catch(console.log) // "error2" from previous .catch()

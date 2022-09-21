// 先执行完同步 再执行reject 后面微任务

new Promise((resolve, reject) => {
  reject(new Error('error1'))
  console.log(1)
})
  .catch(() => console.log('catch'))
  .then(
    () => console.log('then'),
    () => console.log('then2')
  )
  .then(console.log)

console.log(2)

// 1
// 2
// catch
// then
// undefined

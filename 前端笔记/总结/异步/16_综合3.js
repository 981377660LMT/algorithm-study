const p1 = new Promise(resolve => {
  setTimeout(() => {
    resolve('resolve3')
    console.log('timer1') // 3 'timer1'
  }, 0)
  resolve('resovle1')
  resolve('resolve2')
})
  .then(res => {
    console.log(res) // 1  'resolve1'
    setTimeout(() => {
      console.log(p1) // 4 Promise{<resolved>: undefined}  最后一个定时器打印出的p1其实是.finally的返回值
    }, 1000)
  })
  .finally(res => {
    console.log('finally', res) // 2 'finally' undefined
  })
// finally()中的res是一个迷惑项(

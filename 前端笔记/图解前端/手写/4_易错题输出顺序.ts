setTimeout(function () {
  console.log('1')
}, 0)

async function async1() {
  console.log('2')
  const data = await async2() // !注意await 后会先执行同步  resolve是异步 会在下面同步log(7)之后执行
  console.log('3')
  return data
}

async function async2() {
  return new Promise<string>(resolve => {
    console.log('4')
    resolve('async2的结果')
  }).then(data => {
    console.log('5')
    return data
  })
}

async1().then(data => {
  console.log('6')
  console.log(data)
})

new Promise<void>(function (resolve) {
  console.log('7')
  // resolve()  // !没resolve改变状态，不执行then
}).then(function () {
  console.log('8')
})

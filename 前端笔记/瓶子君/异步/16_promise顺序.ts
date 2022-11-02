const myPromise = Promise.resolve(Promise.resolve('666'))

async function func2() {
  const res = await myPromise
  console.log(res, 'here')
  console.log(888)
}

function func1() {
  myPromise.then(res => res).then(res => console.log(res))
  console.log(777)
}

func2()
func1()

// 先执行所有同步 (但是注意await会阻塞)  所有微任务入队
// 777=>666,here=>888=>666
export {}

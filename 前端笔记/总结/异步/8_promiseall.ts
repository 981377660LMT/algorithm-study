function promise1() {
  let p = new Promise(resolve => {
    console.log('promise1') // 1
    resolve('1')
  })
  return p
}
function promise2() {
  return new Promise((resolve, reject) => {
    reject('error')
  })
}

promise1()
  .then(res => console.log(res)) // 2
  .catch(err => console.log(err))
  .finally(() => console.log('finally1')) // 4

promise2()
  .then(res => console.log(res))
  .catch(err => console.log(err)) // 3
  .finally(() => console.log('finally2')) // 5

export {}
// 首先定义了两个函数promise1和promise2，先不管接着往下看。
// promise1函数先被调用了，然后执行里面new Promise的同步代码打印出promise1
// 之后遇到了resolve(1)，将p的状态改为了resolved并将结果保存下来。
// 此时promise1内的函数内容已经执行完了，跳出该函数
// 碰到了promise1().then()，由于promise1的状态已经发生了改变且为resolved因此将promise1().then()这条微任务加入本轮的微任务列表(这是第一个微任务)
// 这时候要注意了，代码并不会接着往链式调用的下面走，也就是不会先将.finally加入微任务列表，那是因为.then本身就是一个微任务，它链式后面的内容必须得等当前这个微任务执行完才会执行，因此这里我们先不管.finally()
// 再往下走碰到了promise2()函数，其中返回的new Promise中并没有同步代码需要执行，所以执行reject('error')的时候将promise2函数中的Promise的状态变为了rejected
// 跳出promise2函数，遇到了promise2().catch()，将其加入当前的微任务队列(这是第二个微任务)，且链式调用后面的内容得等该任务执行完后才执行，和.then()一样。
// OK， 本轮的宏任务全部执行完了，来看看微任务列表，存在promise1().then()，执行它，打印出1，然后遇到了.finally()这个微任务将它加入微任务列表(这是第三个微任务)等待执行
// 再执行promise2().catch()打印出error，执行完后将finally2加入微任务加入微任务列表(这是第四个微任务)
// OK， 本轮又全部执行完了，但是微任务列表还有两个新的微任务没有执行完，因此依次执行finally1和finally2。
// 'promise1'
// '1'
// 'error'
// 'finally1'
// 'finally2'

export {}

// setinterval 用来实现循环定时调用 可能会存在一定的问题 能用 settimeout 解决吗
// 扩展：我们能反过来使用 setinterval 模拟实现 settimeout 吗？
// 扩展思考：为什么要用 settimeout 模拟实现 setinterval？setinterval 的缺陷是什么？
function mySetInterval(cb: (...args: any[]) => void, time: number) {
  let timer: NodeJS.Timer

  const interval = () => {
    cb()
    timer = setTimeout(interval, time)
  }
  interval()

  return {
    cancel: () => clearTimeout(timer),
  }
}

if (require.main === module) {
  mySetInterval(() => {
    console.log(222)
  }, 1000)
  mySetTimeout(() => {
    console.log(111)
  }, 1000)
}

function mySetTimeout(cb: (...args: any[]) => void, time: number) {
  const timer = setInterval(() => {
    clearInterval(timer)
    cb()
  }, time)
}

// 每个setTimeout 产生的任务会直接push 到任务队列中；
// 而setInterval 在每次把任务push 到任务队列前，
// 都要进行一下判断(看上次的任务是否仍在队列中，如果有则不添加，没有则添加)。
// 因而我们一般用setTimeout 模拟setInterval，来规避掉上面的缺点

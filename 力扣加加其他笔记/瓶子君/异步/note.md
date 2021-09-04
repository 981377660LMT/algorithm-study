1. new 的时候 promise 会立马执行 excutor 函数
   **需要包裹在函数里**才能不立马执行 即**需要 await 一个返回 Promise 的函数**
   包括语法糖 Promise.resolve/Promise.reject

2. Promise 按顺序串行:then 串联/reduce+then 串联数组

```JS
const tasks: ((...args: any[]) => Promise<unknown>)[]
tasks.reduce<Promise<unknown>>((pre, cur) => pre.then(cur), Promise.resolve())

效果和
for (const task of tasks) {
  await task()   // 加了await才会顺序执行 不加await就不会顺序执行
}
一样


例子：
function asyncAdd(a: number, b: number, callback: (err: Error | null, result: number) => void) {
  setTimeout(function () {
    callback(null, a + b)
  }, 1000)
}

async function sumTwo(a: number, b: number) {
  return new Promise<number>((resolve, reject) => {
    asyncAdd(a, b, (err, result) => {
      if (!err) resolve(result)
      else reject(err)
    })
  })
}

// 多数之和版本1  Promise串行
async function sumAll(...nums: number[]) {
  return new Promise(resolve =>
    nums
      .reduce((pre, cur) => pre.then(total => sumTwo(total, cur)), Promise.resolve(0))
      .then(resolve)
  )
}
```

Promise 并发执行：Promise.all/for of +await/for await of

```JS
// 多数之和版本2  Promise 两两合并任务 并行任务
async function sumAll2(...nums: number[]): Promise<number> {
  console.log(nums)
  // 两两一组分组 不足的拿出来
  if (nums.length === 0) return 0
  if (nums.length === 1) return nums[0]
  if (nums.length === 2) return await sumTwo(nums[0], nums[1])

  const tasks: Promise<number>[] = []
  for (let i = 0; i < nums.length - 1; i += 2) {
    tasks.push(sumTwo(nums[i], nums[i + 1]))
  }
  if (nums.length % 2) tasks.push(Promise.resolve<number>(nums[nums.length - 1]))

  return sumAll2(...(await Promise.all(tasks)))
}
```

不 resolve 的 Promise **加上 await 会阻塞整个函数**

```JS
const task1 = new Promise(resolve => console.log('ok1'))
```

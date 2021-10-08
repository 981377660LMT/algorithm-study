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

// 多数之和版本1  Promise+reduce串行
async function sumAll1(...nums: number[]) {
  return new Promise(resolve =>
    nums
      .reduce((pre, cur) => pre.then(total => sumTwo(total, cur)), Promise.resolve(0))
      .then(resolve)
  )
}

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

async function test() {
  console.log(await sumTwo(1, 2))
  // console.log(await sumAll1(1, 2, 3))
  console.log(await sumAll2(1, 2, 3, 4, 5, 6))
}
test()

export {}

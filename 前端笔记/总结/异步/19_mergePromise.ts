const time = (timer: number) => {
  return new Promise<void>(resolve => {
    setTimeout(() => {
      resolve()
    }, timer)
  })
}
const ajax1 = () =>
  time(2000).then(() => {
    console.log(1)
    return 1
  })
const ajax2 = () =>
  time(1000).then(() => {
    console.log(2)
    return 2
  })
const ajax3 = () =>
  time(1000).then(() => {
    console.log(3)
    return 3
  })

// 这道题有点类似于Promise.all()，不过.all()并发执行。
// 这里需要等上一个执行完毕之后才能执行下一个，串行
function mergePromise<T>(promises: (() => Promise<T>)[]): Promise<T[]> {
  // @ts-ignore
  promises.push(() => Promise.resolve<T>())

  return new Promise(resolve => {
    const res: any[] = []

    promises.reduce(
      (pre, cur) =>
        pre
          // @ts-ignore
          .then(data => {
            data != undefined && res.push(data)
            if (res.length === promises.length - 1) return resolve(res)
            return cur()
          }),
      Promise.resolve()
    )
  })
}

mergePromise([ajax1, ajax2, ajax3]).then(data => {
  console.log('done')
  console.log(data) // data 为 [1, 2, 3]
})

export {}

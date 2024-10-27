function createControllablePromise<T>(): {
  promise: Promise<T>
  resolve: (value: T) => void
  reject: (reason: unknown) => void
} {
  let res: any = {}
  res.promise = new Promise<T>((resolve, reject) => {
    res.resolve = resolve
    res.reject = reject
  })
  return res
}

export {}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const { promise, resolve, reject } = createControllablePromise<number>()
  promise.then(console.log)
  resolve(42)
  reject(new Error('error')) // 无效
}

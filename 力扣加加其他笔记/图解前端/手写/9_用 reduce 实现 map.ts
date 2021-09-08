const a = implementMapUsingReduce([1, 2, 3, 4], (a: number) => a + 1) // [2,3,4,5]
console.log(a)

const b = implementMapUsingReduce(['a', 'b', 'c'], (a: string) => a + '!') // ['a!', 'b!', 'c!']
console.log(b)

function implementMapUsingReduce<T>(arr: T[], func: (...args: T[]) => T) {
  return arr.reduce<T[]>((pre, cur, index) => {
    pre[index] = func(cur)
    return pre
  }, [])
}

export default 1

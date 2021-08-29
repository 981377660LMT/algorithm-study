interface NestedArray<T> extends Array<T | NestedArray<T>> {}
// type NestedArray<T> = Array<T | NestedArray<T>>
const flatten = (nums: NestedArray<number>): number[] => {
  const res: number[] = []
  for (const num of nums) {
    if (typeof num === 'number') res.push(num)
    else res.push(...flatten(num))
  }
  return res
}

console.log(flatten([1, 2, [3, 4, [5]], 5]))

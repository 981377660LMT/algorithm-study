function sum(nums: number[]): number {
  function inner(res = 0, index = 0): number {
    if (index === nums.length) {
      return res
    }

    return inner(res + nums[index], index + 1)
  }

  return inner()
}

// Maximum call stack size exceeded
console.log(sum(Array.from({ length: 10000 }, (v, i) => i)))
export {}

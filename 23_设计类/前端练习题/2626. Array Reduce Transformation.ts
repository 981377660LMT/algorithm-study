type Op = (pre: number, cur: number) => number

function reduce(nums: number[], op: Op, e: number): number {
  nums.forEach(cur => {
    e = op(e, cur)
  })
  return e
}

export {}

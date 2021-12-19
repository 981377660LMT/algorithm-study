// 给定一个乱序的栈，设计算法将其升序排列。

// 维护辅助栈为单调递增栈。
// 保证“倒腾”过程的任何时候，辅助栈的元素都是从小到大排序的！
// 当原始栈为空时，辅助栈便满足[1, 2, 3, 4]。
function sortStack(stack: number[]): number[] {
  const incr: number[] = []

  while (stack.length > 0) {
    const cur = stack.pop()!
    // 进来小的，直接把里面大的放回去
    while (incr.length > 0 && incr[incr.length - 1] > cur) {
      stack.push(incr.pop()!)
    }
    incr.push(cur)
  }

  return incr
}

const stack = [4, 2, 1, 3]
console.log(sortStack(stack))

export {}

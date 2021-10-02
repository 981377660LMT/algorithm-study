// arr是一个可能包含重复元素的整数数组，
// 我们将这个数组分割成几个“块”，并将这些块分别进行排序。
// 之后再连接起来，使得连接的结果和按升序排序后的原数组相同。
// 我们最多能将数组分成多少块？

function maxChunksToSorted(arr: number[]): number {
  const stack: number[] = []

  // 来了小的 要看前面是不是都比它小 如果不是 必须合在一起
  // 考虑整体，比如 [4,2,2,1,1] 这样的测试用例，实际上只应该返回 1，原因是后面碰得到了 1，使得前面不应该分块。
  for (const num of arr) {
    if (stack.length && stack[stack.length - 1] > num) {
      // 不是分割区块，而是`融合`区块,我们需要将融合后的区块的最大值重新放回栈
      const top = stack.pop()!
      while (stack.length && stack[stack.length - 1] > num) {
        stack.pop()
      }
      stack.push(top)
    } else {
      stack.push(num)
    }
  }

  return stack.length
}

console.log(maxChunksToSorted([2, 1, 3, 4, 4]))
// 输出: 4
// 解释:
// 我们可以把它分成两块，例如 [2, 1], [3, 4, 4]。
// 然而，分成 [2, 1], [3], [4], [4] 可以得到最多的块数。

console.log(maxChunksToSorted([5, 4, 3, 2, 1]))
// 输出: 1
// 解释:
// 将数组分成2块或者更多块，都无法得到所需的结果。
// 例如，分成 [5, 4], [3, 2, 1] 的结果是 [4, 5, 1, 2, 3]，这不是有序的数组。

// 实际上本题的单调栈思路和 【力扣加加】从排序到线性扫描(57. 插入区间)
// 以及 394. 字符串解码 都有部分相似，大家可以结合起来理解。
// `融合`与【力扣加加】从排序到线性扫描(57. 插入区间) 相似，
// `重新压栈`和 394. 字符串解码 相似。

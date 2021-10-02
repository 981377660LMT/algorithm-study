// /**
//  * @param {number[]} arr
//  * @return {number}
//  * 进行奇数跳跃时 你将会跳到索引 j，使得 A[i] <= A[j]，A[j] 是可能的最小值
//  * 如果存在多个这样的索引 j，你只能跳到满足要求的最小索引 j 上
//  * 进行偶数跳跃时 你将会跳到索引 j，使得 A[i] >= A[j]，A[j] 是可能的最大值
//  * 如果从某一索引开始跳跃一定次数（可能是 0 次或多次），就可以到达数组的末尾（索引 A.length - 1），那么该索引就会被认为是好的起始索引。
//  * 返回好的起始索引的数量。
//  */
// const oddEvenJumps = function (arr: number[]): number {
//   const n = arr.length
//   const jumps = arr.map((v, i) => [v, i]).sort((a, b) => a[0] - b[0] || a[1] - b[1])

//   // 右侧第一个比自己大的（相同大的，取index小的）index
//   const nextBigger = Array(n).fill(-1)
//   const descStack: number[] = []
//   for (const [_, id] of jumps) {
//     // it means stack[-1]'s next bigger(or equal) is i
//     while (descStack.length && descStack[descStack.length - 1] < id) {
//       nextBigger[descStack.pop()!] = id
//     }
//     descStack.push(id)
//   }

//   // 右侧第一个比自己小的（相同大小，取index小的）index
//   const nextSmaller = Array(n).fill(-1)
//   const incStack: number[] = []
//   for (const [_, id] of jumps.reverse()) {
//     // it means stack[-1]'s next smaller(or equal) is i
//     while (incStack.length && incStack[incStack.length - 1] < id) {
//       nextSmaller[incStack.pop()!] = id
//     }
//     incStack.push(id)
//   }
//   // console.log(nextSmaller, nextBigger)
//   // 是否可从 i 开始通过跳高的方式奇偶跳到达数组最后一项
//   const higher = Array(n).fill(false)
//   const lower = Array(n).fill(false)
//   higher[n - 1] = true
//   lower[n - 1] = true

//   let res = 1
//   for (let i = n - 2; ~i; i--) {
//     higher[i] = lower[nextBigger[i]]
//     lower[i] = higher[nextSmaller[i]]
//     res += higher[i] ? 1 : 0
//   }

//   // 起始奇数
//   return res
// }

// console.log(oddEvenJumps([2, 3, 1, 1, 4]))
// // 输出：3

// // 这种题目一般都是倒着思考比较容易。因为我虽然不知道你从哪开始可以跳到最后，
// // 但是我知道最终被计算进返回值的一定是在数组末尾结束的。

// // 单调栈解决的是右侧第一个比当前大/小的idnex
// // 本题是右侧比当前“最大”/“最小”的index
// // 故需要先排序

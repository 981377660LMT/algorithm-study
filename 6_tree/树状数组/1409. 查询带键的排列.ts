import { BIT } from './树状数组单点更新模板'

// 数组中的元素为 1 到 m 之间的正整数
function processQueries(queries: number[], m: number): number[] {
  const bit = new BIT(10000)
}

console.log(processQueries([3, 1, 2, 1], 5))
// 输出：[2,1,2,1]
// 解释：待查数组 queries 处理如下：
// 对于 i=0: queries[i]=3, P=[1,2,3,4,5], 3 在 P 中的位置是 2，接着我们把 3 移动到 P 的起始位置，得到 P=[3,1,2,4,5] 。
// 对于 i=1: queries[i]=1, P=[3,1,2,4,5], 1 在 P 中的位置是 1，接着我们把 1 移动到 P 的起始位置，得到 P=[1,3,2,4,5] 。
// 对于 i=2: queries[i]=2, P=[1,3,2,4,5], 2 在 P 中的位置是 2，接着我们把 2 移动到 P 的起始位置，得到 P=[2,1,3,4,5] 。
// 对于 i=3: queries[i]=1, P=[2,1,3,4,5], 1 在 P 中的位置是 1，接着我们把 1 移动到 P 的起始位置，得到 P=[1,2,3,4,5] 。
// 因此，返回的结果数组为 [2,1,2,1] 。

export {}
// todo

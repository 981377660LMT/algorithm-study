// https://leetcode.cn/problems/linked-list-random-node/solution/gong-shui-san-xie-xu-shui-chi-chou-yang-1lp9d/
// 给定一个数据流，数据流长度N很大，
// 且N直到处理完所有数据之前都不可知，
// 请问如何在只遍历一遍数据（O(N)）的情况下，
// 能够随机选取出k个不重复的数据。

// 算法的核心在于先以某一种概率选取数，
// 并在后续过程以另一种概率换掉之前已经被选中的数。
// 因此实际上每个数被最终选中的概率都是被选中的概率 * 不被替换的概率。

// 基本思路是
// 1.构建一个大小为 k 的数组，将数据流的前 k 个元素放入数组中
// 2. 从数据流k+1之后的数开始，设为第i个数， [0, i] 之间选一个数 rand，如果rand落在[1,k]则将rand与i交换

// 对于前 k 个数，最终被选择的概率都是
// 1 * 不被 k + 1 替换的概率 * 不被 k + 2 替换的概率 * ... 不被 n 替换的概率=k/n

// 对于 第 i (i > k) 个数，最终被选择的概率是
// 第 i 步被选中的概率 * 不被第 i + 1 步替换的概率 * ... * 不被第 n 步被替换的概率=k/n

const pick = (dataStream: number[], k: number) => {
  const reservoir = Array<number>(k)

  for (let i = 0; i < k; i++) {
    reservoir[i] = dataStream[i]
  }

  for (let i = k; i < dataStream.length; i++) {
    // [0,i]
    const rand = Math.floor(Math.random() * (i + 1))
    if (rand <= k - 1) {
      reservoir[rand] = dataStream[i]
    }
  }

  return reservoir
}

console.log(
  pick(
    Array(100)
      .fill(0)
      .map((_, i) => i + 1),
    3
  )
)

// 力扣想进行的操作有以下三种：

import { BIT2 } from './BIT'

// 给团队的一个成员（也可以是负责人）发一定数量的LeetCoin；
// 给团队的一个成员（也可以是负责人），以及他/她管理的所有人（即他/她的下属、他/她下属的下属，……），发一定数量的LeetCoin；
// 查询某一个成员（也可以是负责人），以及他/她管理的所有人被发到的LeetCoin之和。

// https://leetcode-cn.com/problems/coin-bonus/solution/xiao-ai-lao-shi-li-kou-bei-li-jie-zhen-t-rut3/
// https://mp.weixin.qq.com/s?__biz=MzkyMzI3ODgzNQ==&mid=2247483674&idx=1&sn=263595b26950ac60e5bf789839d70c9e&chksm=c1e6cd86f691449062d780b96d9af6d9590a71872ebfee980d5b045b9963714043261027c16b&token=1500097142&lang=zh_CN#rd
// 1. dfs将管理和他管理的人映射到一个区间(这部分很巧妙)[a,b] b表示自身的id
// 2. 树状数组区间update/query
const MOD = 1e9 + 7
function bonus(n: number, leadership: number[][], operations: number[][]): number[] {
  const adjList = Array.from<number, number[]>({ length: n + 1 }, () => [])
  const start = Array<number>(n + 1).fill(0) // 子树最开始的结点序号
  const end = Array<number>(n + 1).fill(0) // 本身最后映射到几
  // begin[1] = 1, end[1] = 6，表示编号为 1 的人所管理的团队映射到的区间是 [1, 6]，本身映射到 6
  let id = 1

  for (const [u, v] of leadership) adjList[u].push(v)

  dfs(1)

  const res: number[] = []
  const bit = new BIT2(n + 10)
  for (const [optType, optId, optValue] of operations) {
    switch (optType) {
      case 1:
        bit.add(end[optId], end[optId], optValue)
        break
      case 2:
        bit.add(start[optId], end[optId], optValue)
        break
      case 3:
        const queryRes = bit.query(start[optId], end[optId])
        res.push(((queryRes % MOD) + MOD) % MOD)
        break
      default:
        throw new Error('invalid optType')
    }
  }

  return res

  // dfs序
  function dfs(cur: number): void {
    start[cur] = id
    for (const next of adjList[cur]) dfs(next)
    // id在dfs过程中被改变了
    end[cur] = id
    id++
  }
}

console.log(
  bonus(
    6,
    [
      [1, 2],
      [1, 6],
      [2, 3],
      [2, 5],
      [1, 4]
    ],
    [
      [1, 1, 500],
      [2, 2, 50],
      [3, 1],
      [2, 6, 15],
      [3, 1]
    ]
  )
)
// 第一次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 0;
// 第二次查询时，每个成员得到的LeetCoin的数量分别为（按编号顺序）：500, 50, 50, 0, 50, 15.
// 输出：[650, 665]
export {}

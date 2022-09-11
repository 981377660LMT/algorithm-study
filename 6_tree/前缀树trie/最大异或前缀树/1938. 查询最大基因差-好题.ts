/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-shadow */
// # 本题与「1707. 与数组中元素的最大异或值」是非常相似的题
// # 离线算法 + 字典树
// 两个基因值的 基因差 是两者的 异或和
// queries[i] = [nodei, vali] 。对于查询 i ，请你最大化 vali XOR pi 。
// !其中 pi 是节点 nodei 到根之间的任意节点
// 0 <= vali <= 2 * 1e5  说明不超过18位

// 1.离线查询:在dfs过程插入节点，dfs回溯阶段删除节点，
// 保证查询时树里只有到根节点的数字，
// 将每个询问在遍历到对应结点时求解
// 2.XORTrie树：查询num与树里的最大异或值

import { useArrayXORTrie } from './XORTrie-数组节省空间版'

function maxGeneticDifference(parents: number[], queries: number[][]): number[] {
  const res = Array<number>(queries.length).fill(0)

  // 预处处理查询
  const nodeQueries = new Map<number, [qv: number, qi: number][]>()
  for (const [qi, [root, qv]] of queries.entries()) {
    !nodeQueries.has(root) && nodeQueries.set(root, [])
    nodeQueries.get(root)!.push([qv, qi])
  }

  // 建树
  let root = -1
  const adjMap = new Map<number, number[]>()
  for (const [cur, pre] of parents.entries()) {
    if (pre === -1) {
      root = cur
    } else {
      !adjMap.has(pre) && adjMap.set(pre, [])
      adjMap.get(pre)!.push(cur)
    }
  }

  // dfs遍历树，查询trie并更新答案
  const bit = Math.floor(Math.log2(2e5)) + 1
  const xorTrie = useArrayXORTrie(bit)
  dfs(root)
  return res

  function dfs(root: number): void {
    xorTrie.insert(root)
    for (const [qv, qi] of nodeQueries.get(root) || []) res[qi] = xorTrie.search(qv)
    for (const next of adjMap.get(root) || []) dfs(next)
    xorTrie.remove(root)
  }
}

console.log(
  maxGeneticDifference(
    [3, 7, -1, 2, 0, 7, 0, 2],
    [
      [4, 6],
      [1, 15],
      [0, 5]
    ]
  )
)

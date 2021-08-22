import { UnionFind } from '../0_并查集'
import { useUnionFind } from '../推荐使用并查集精简版'

/**
 * @param {string} s
 * @param {number[][]} pairs
 * @return {string}
 * @description 你可以 任意多次交换 在 pairs 中任意一对索引处的字符。(注意:如果pairs重叠则联通)
 * 返回在经过若干次交换后，s 可以变成的按字典序最小的字符串。
 * @summary 很重要的方法就是将每个数的root找到并保存在map中
 * 第一次提交超时
 */
const smallestStringWithSwaps = function (s: string, pairs: number[][]): string {
  const groupMap = new Map<number, string[]>()
  const uf = useUnionFind(s.length)
  for (const [v, w] of pairs) {
    uf.union(v, w)
  }

  // 记录每个字符串所属组的根节点
  for (let i = 0; i < s.length; i++) {
    const root = uf.find(i)
    if (!groupMap.has(root)) groupMap.set(root, [])
    groupMap.get(root)!.push(s[i])
  }

  // 字典序在前的先去 以便之后pop
  for (const [_, v] of groupMap) {
    v.sort((a, b) => b.localeCompare(a))
  }

  const res = Array<string>(s.length)
  for (let i = 0; i < s.length; i++) {
    const root = uf.find(i)
    res.push(groupMap.get(root)!.pop()!)
  }

  return res.join('')
}

console.log(
  smallestStringWithSwaps('dcab', [
    [0, 3],
    [1, 2],
  ])
)
// 输出："bacd"
// 解释：
// 交换 s[0] 和 s[3], s = "bcad"
// 交换 s[1] 和 s[2], s = "bacd"

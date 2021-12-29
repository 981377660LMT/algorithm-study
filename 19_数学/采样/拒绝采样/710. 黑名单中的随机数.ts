import { randint } from '../randint'

type Black = number
type White = number

class Solution {
  private whiteListLength: number
  private black: Map<Black, White>

  /**
   *
   *
   * @param n N 的最大值为 1000000000
   * @param blacklist
   * 给定一个包含 [0，n) 中不重复整数的黑名单 blacklist
   */
  constructor(n: number, blacklist: number[]) {
    this.whiteListLength = n - blacklist.length
    this.black = new Map()
    const set = new Set(blacklist)
    let lastIndex = n - 1
    for (const num of blacklist) {
      // 不在白名单索引里 不管
      if (num >= this.whiteListLength) continue
      // 将黑名单中的数对应到后面的非黑名单数(清理白名单中的杂质)
      while (set.has(lastIndex)) {
        lastIndex--
      }
      this.black.set(num, lastIndex--)
    }
  }

  /**
   * 从 [0, n) 中返回一个不在 blacklist 中的随机整数,黑名单以外的数要等概率出现
   * @description
   * 类似于解决哈希冲突，但是任何确定性的探测方法都会导致非均匀分布
   * @summary
   * 黑名单映射 ：默认白名单为前whiteList个索引 然后利用黑名单映射除去杂质
   * @link
   * https://leetcode-cn.com/problems/random-pick-with-blacklist/solution/zhong-jiang-guan-fang-ti-jie-hei-ming-dan-ying-she/
   */
  pick(): number {
    const rand = randint(0, this.whiteListLength - 1)
    return this.black.get(rand) || rand
  }
}

const s = new Solution(4, [2])
console.log(s.pick())
console.log(s.pick())
console.log(s.pick())
export {}

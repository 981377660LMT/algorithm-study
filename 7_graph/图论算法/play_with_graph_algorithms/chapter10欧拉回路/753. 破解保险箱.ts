/**
 * @param {number} n n 的范围是 [1, 4]
 * @param {number} k k 的范围是 [1, 10] k^n 最大可能为 4096。
 * @return {string}
  返回一个最短的字符串满足：n位每位由0~k-1构成的字符串都是其子串
 * 如何在一个最短的串内枚举所有的n位k进制数排列:这种序列称之为 de Bruijn序列。
 */
const crackSafe = function (n: number, k: number): string {
  const res: string[] = []

  // 开始的选择
  const start = '0'.repeat(n - 1)
  dfs(start, new Set<string>(), res)
  res.push(start)
  return res.join('')

  function dfs(cur: string, visited: Set<string>, path: string[]) {
    for (let i = 0; i < k; i++) {
      const newSubstring = cur + i.toString()
      if (visited.has(newSubstring)) continue
      visited.add(newSubstring)
      dfs(newSubstring.slice(1), visited, path)
      path.push(i.toString())
    }
  }
}

console.log(crackSafe(2, 2))
// 输出: "00110"
// 说明: "01100", "10011", "11001" 也能打开保险箱。

// De Bruijn序列 （德布鲁因序列）
// 能否构造一个长度为2的n次方的二进制环状串，使得二进制环状串中总共2的n次方个长为n的不同截断作为2的n次方个长为n的二进制串来说互不相同。
// 例如,当n为2时，这样的环状串可以是：0011，它的4个长度为2的子串分别为00、01、11、10，这四个子串互不相同。

// 密码共有 n 位，每一个位可以有 k 个数字，
// 那么总共不同的密码总数就有 k^n 个，而每个密码可以公用 n - 1 位，
// 所以破解保险箱的密码最短长度为：(n - 1) + k^n 位。

/**
 * @param {number} n
 * @return {string[]}
 * @description
 * 找到所有长度为 n 的中心对称数。
 * 后序遍历
 */
function findStrobogrammatic(n: number): string[] {
  const mapping = new Map<string, string>([
    ['1', '1'],
    ['8', '8'],
    ['6', '9'],
    ['9', '6'],
  ])

  return dfs(n)

  function dfs(remainLength: number): string[] {
    if (remainLength === 0) return ['']
    else if (remainLength === 1) return ['0', '1', '8']

    const res: string[] = []
    for (const num of dfs(remainLength - 2)) {
      for (const [head, tail] of mapping) {
        res.push(`${head}${num}${tail}`)
      }

      if (remainLength !== n) {
        res.push(`0${num}0`)
      }
    }

    return res
  }
}

if (require.main === module) {
  console.log(findStrobogrammatic(3))
}

export { findStrobogrammatic }

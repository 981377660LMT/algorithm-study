/**
 * @param {string[]} queries
 * @param {string} pattern
 * @return {boolean[]}
 * 我们可以将小写字母插入模式串 pattern 得到待查询项 query，那么待查询项与给定模式串匹配。（
 * @summary
 * 1.如果大写在pattern中不存在 则返回false
 * 2.如果匹配长度小于pattern长度 返回false
 */
function camelMatch(queries: string[], pattern: string): boolean[] {
  return queries.map(query => {
    let hit = 0
    for (const char of query) {
      if (pattern[hit] === char) {
        hit++
      } else if (char === char.toUpperCase()) {
        return false
      }
    }

    return hit === pattern.length
  })
}

console.log(camelMatch(['FooBar', 'FooBarTest', 'FootBall', 'FrameBuffer', 'ForceFeedBack'], 'FB'))

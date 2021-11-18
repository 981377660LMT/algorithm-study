// 按 任意顺序 返回答案。
// 1 <= word.length <= 15
function generateAbbreviations(word: string): string[] {
  const res: string[] = []
  dfs(0, 0, '')
  return res

  /**
   *
   * @param index 现在位置
   * @param count 数字替换了几个字符
   * @param path 暂时的字符串
   */
  function dfs(index: number, count: number, path: string): void {
    if (index === word.length) {
      count > 0 && (path += count.toString())
      res.push(path)
      return
    }

    // 是否消耗掉count
    dfs(index + 1, count + 1, path)
    dfs(index + 1, 0, path + (count > 0 ? count.toString() : '') + word[index])
  }
}

console.log(generateAbbreviations('word'))
// 输出：["4","3d","2r1","2rd","1o2","1o1d","1or1","1ord","w3","w2d","w1r1","w1rd","wo2","wo1d","wor1","word"]

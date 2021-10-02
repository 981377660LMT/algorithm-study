// 有效的单词方块是指此由单词序列组成的文字方块的 第 k 行 和 第 k 列 (0 ≤ k < max(行数, 列数)) 所显示的字符串完全相同。
function validWordSquare(words: string[]): boolean {
  for (const [index, word] of zipLongest(...words).entries()) {
    console.table(zipLongest('', ...words))
    if (words[index] !== word.join('')) return false
  }

  /**
   *
   * @param strs
   * @returns 将每个列组合 则有最大长度的行，原来的单词数的列
   */
  function zipLongest(fillValue = '', ...strs: string[]) {
    const length = Math.max(...strs.map(str => str.length))
    // 行对列，列对行
    const arr = Array.from({ length }, () => Array(strs.length).fill(''))
    for (let i = 0; i < strs.length; i++) {
      for (let j = 0; j < length; j++) {
        arr[j][i] = strs[i][j] || fillValue
      }
    }
    return arr
  }

  return true
}

console.log(validWordSquare(['ball', 'area', 'read', 'lady']))

// 输出：
// true

// 解释：
// 第 1 行和第 1 列都是 "abcd"。
// 第 2 行和第 2 列都是 "bnrt"。
// 第 3 行和第 3 列都是 "crmy"。
// 第 4 行和第 4 列都是 "dtye"。

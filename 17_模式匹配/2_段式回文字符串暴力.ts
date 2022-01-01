/**
 * @param {string} text
 * @return {number}
 * @description 给你一个字符串 text，在确保它满足段式回文的前提下，请你返回 段 的 最大数量 k
 * @description 比较字符串相等需要length复杂度 而比较数字只需要1复杂度
 * @description todo
 */
function longestDecomposition(text: string): number {
  let res = 0
  let l = 1

  while (l < text.length) {
    const sub1 = text.slice(0, l)
    const sub2 = text.slice(text.length - l)
    if (sub1 === sub2) {
      res += 2
      text = text.slice(l, text.length - l)
      l = 1
    } else {
      l++
    }
  }

  return text === '' ? res : res + 1
}

console.log(longestDecomposition('asdkasd'))

export {}

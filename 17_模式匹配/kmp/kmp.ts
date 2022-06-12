/**
 *
 * @param shorter 模式串的next数组
 * @returns `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
 */
function getNext(shorter: string): number[] {
  const next = Array<number>(shorter.length).fill(0)
  let j = 0

  for (let i = 1; i < shorter.length; i++) {
    while (j > 0 && shorter[i] !== shorter[j]) {
      //  前进到最长公共后缀结尾处
      j = next[j - 1]
    }

    if (shorter[i] === shorter[j]) j++
    next[i] = j
  }

  return next
}

export { getNext, getNext as getLPS }

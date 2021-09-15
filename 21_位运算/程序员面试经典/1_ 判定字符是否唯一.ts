/**
 * @param {string} astr 小写字母
 * @return {boolean}
 */
var isUnique = function (astr: string): boolean {
  let set = 0
  for (const str of astr) {
    console.log(str)
    const key = 1 << (str.codePointAt(0)! - 97)
    if (set & key) return false
    set |= key
  }
  return true
}

// console.log(isUnique('abc'))
console.log(isUnique('aa'))
console.log(1 & 3)

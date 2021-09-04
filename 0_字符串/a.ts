/**
 * @param {string} s
 * @return {number}
 */
var lengthOfLastWord = function (s: string): number {
  const tmp = s.split(/\s+/g).filter(Boolean)
  return tmp[tmp.length - 1].length
}

console.log(lengthOfLastWord('Hello World'))

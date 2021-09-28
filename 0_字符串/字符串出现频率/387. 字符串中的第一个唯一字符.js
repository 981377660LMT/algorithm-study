/**
 * @param {string} s
 * @return {number}
 * 哈希表/api
 */
var firstUniqChar = function (s) {
  for (let i = 0; i < s.length; i++) {
    const char = s[i]
    // 判断字符的第一个索引和最后一个索引是否相等
    if (s.indexOf(char) === s.lastIndexOf(char)) return i
  }
  return -1
}

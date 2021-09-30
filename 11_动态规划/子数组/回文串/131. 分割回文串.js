/**
 * @param {string} s  1 <= s.length <= 16
 * @return {string[][]}
 * 请你将 s 分割成一些子串，使每个子串都是 回文串 。返回 s 所有可能的分割方案
 */
var partition = function (s) {
  const res = []
  const isPalindrome = str => str === str.split('').reverse().join('')

  const bt = (remain, path) => {
    // console.log(curLen, path)
    if (remain.length === 0) {
      res.push(path.slice())
      return
    }

    for (let i = 0; i < remain.length; i++) {
      const sub = remain.slice(0, i + 1)
      if (isPalindrome(sub)) {
        path.push(sub)
        bt(remain.slice(i + 1), path)
        path.pop()
      }
    }
  }
  bt(s, [])

  return res
}

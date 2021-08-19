/**
 * @param {string} s
 * @param {string} t
 * @return {boolean}
 * # 代表退格字符。
 * 逆序地遍历字符串，就可以立即确定当前字符是否会被删掉。
 * 我们定义 skip 表示当前待删除的字符的数量
 */
const backspaceCompare = function (s: string, t: string): boolean {
  /*
      从后向前，遍历并比较两个字符串：
      (1)首先遍历S(或是 先遍历T也行)：
          1、遇到'#'，就记录个数(sWell++)，让指针前移
          2、若不是'#'，但skipS大于0(表示还有未抵消的'#')，则抵消当前字符，让指针前移
          3、若上述两点都不满足，则结束当前循环，进行后续步骤
      (2)遍历T，如上进行操作
      (3)比较当前S和T的字符，若不相等，则返回false
      (4)往复循环如上步骤，直至任何一个字符串遍历完毕
  */

  let i = s.length - 1
  let j = t.length - 1
  let skipS = 0
  let skipT = 0
  while (i >= 0 || j >= 0) {
    console.log(i, j)
    while (i >= 0) {
      if (s.charAt(i) === '#') {
        skipS++
        i--
      } else if (skipS > 0) {
        skipS--
        i--
      } else {
        break
      }
    }

    // (2)遍历j，如上进行操作
    while (j >= 0) {
      if (t.charAt(j) === '#') {
        skipT++
        j--
      } else if (skipT > 0) {
        skipT--
        j--
      } else {
        break
      }
    }

    // (3)比较当前S和T的字符，若不相等，则返回false
    if (i >= 0 && j >= 0) {
      if (s.charAt(i) !== t.charAt(j)) return false
    } else {
      // 一个走完一个没走完
      if (i >= 0 || j >= 0) return false
    }

    i--
    j--
  }

  return true
}

console.log(backspaceCompare('xywrrmp', 'xywrrmu#p'))

export default 1

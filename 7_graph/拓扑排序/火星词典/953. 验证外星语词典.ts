// 只有当给定的单词在这种外星语中按字典序排列时，返回 true；否则，返回 false。

// 因为递增排序，所以只要比较相邻即可
function isAlienSorted(words: string[], order: string): boolean {
  const ailenOrder = Object.fromEntries([...order.split('').entries()].map(([i, v]) => [v, i]))

  for (let i = 0; i < words.length - 1; i++) {
    const [pre, cur] = [words[i], words[i + 1]]

    let allEqual = true
    for (let j = 0; j < Math.min(pre.length, cur.length); j++) {
      if (pre[j] !== cur[j]) {
        allEqual = false

        if (ailenOrder[pre[j]] > ailenOrder[cur[j]]) {
          return false
        }

        break // 只比较一次
      }
    }

    if (allEqual) {
      if (pre.length > cur.length) {
        return false
      }
    }
  }

  return true
}

console.log(isAlienSorted(['hello', 'leetcode'], 'hlabcdefgijkmnopqrstuvwxyz'))
// console.log(isAlienSorted(['apple', 'app'], 'abcdefghijklmnopqrstuvwxyz'))
// console.log(isAlienSorted(['word', 'world', 'row'], 'worldabcefghijkmnpqstuvxyz'))
// 当前三个字符 "app" 匹配时，第二个字符串相对短一些，然后根据词典编纂规则
//  "apple" > "app"，因为 'l' > '∅'，其中 '∅' 是空白字符，
//  定义为比任何其他字符都小

export {}
// 1.比较每个单词出现的第一个不同的字母，如果index小于后面的单词，就不再比较 大于则return false
// 2.如果没出现不同字母&&前一个单词的长度比后一个单词长度大，返回false

// Python语法糖——for/else循环语句里的break
// 作用是for循环中if条件一直不满足，则最后就执行else语句
// 但是if由break语句中断时，else就不被执行

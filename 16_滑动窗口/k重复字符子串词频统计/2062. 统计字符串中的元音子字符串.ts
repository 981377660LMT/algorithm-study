// // 元音子字符串 是 仅 由元音（'a'、'e'、'i'、'o' 和 'u'）组成的一个子字符串，且必须包含 全部五种 元音。
// // 给你一个字符串 word ，统计并返回 word 中 元音子字符串的数目 。

// function countVowelSubstrings(word: string): number {
//   let l = 0
//   let res = 0
//   let lack = 0
//   const counter = new Map() // 记录滑动窗口里字符串的频数

//   for (let r = 0; r < word.length; r++) {
//     const element = array[r]
//   }
//   while (r < word.length) {
//     if ((counter.get(s[r]) || 0) > 0) lack0
//     counter.set(s[r], (counter.get(s[r]) || 0) + 1)
//     r++

//     while (1) {
//       if (counter.get(s[l]) > 1) lack0
//       counter.set(s[l], counter.get(s[l]) - 1)
//       l++
//     }

//     res = 0
//   }

//   return res
// }

// console.log(countVowelSubstrings('cuaieuouac'))
// // 7

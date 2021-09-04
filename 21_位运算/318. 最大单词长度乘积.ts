/**
 * @param {string[]} words
 * @return {number}
 * 找到 length(word[i]) * length(word[j]) 的最大值，并且这两个单词不含有公共字母
 * abc，bac，cba...这样字符串其实都是一个意思，我们是否能找到一样东西可以代表它们的。
 * 我们用二进制表示，
   a->1，b->10....
   ab->1|10->11，ba同样也是11(或运算) 
   把每个单词的二进制都记录了，但是 比如 abc，aaabc两个二进制是一样的，我们只需要记录长的单词即可
 */
var maxProduct = function (words: string[]): number {
  let res = 0
  const record = new Map<number, number>()

  for (const word of words) {
    let mask = 0
    for (const char of word) {
      // 这一步很关键
      mask |= 1 << (char.codePointAt(0)! - 97)
    }
    const max = record.get(mask) || -Infinity
    record.set(mask, Math.max(max, word.length))
  }

  const masks = [...record.keys()]
  for (let i = 0; i < masks.length - 1; i++) {
    for (let j = i + 1; j < masks.length; j++) {
      if (masks[i] & masks[j]) continue
      res = Math.max(res, record.get(masks[i])! * record.get(masks[j])!)
    }
  }

  console.log(record)
  return res
}

console.log(maxProduct(['abcw', 'baz', 'foo', 'bar', 'xtfn', 'abcdef']))

export {}

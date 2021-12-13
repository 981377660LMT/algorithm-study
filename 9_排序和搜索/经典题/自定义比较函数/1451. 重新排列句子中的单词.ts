// 请你重新排列 text 中的单词，使所有单词按其长度的升序排列。如果两个单词的长度相同，则`保留其在原句子中的相对顺序`。
// 本身sort就是稳定排序 只要排序长度即可
function arrangeWords(text: string): string {
  const words = text.toLowerCase().split(' ')
  const res = words.sort((w1, w2) => w1.length - w2.length).join(' ')
  return res[0].toUpperCase() + res.slice(1)
}

console.log(arrangeWords('Keep calm and code on'))
// 输出："On and keep calm code"
// 解释：输出的排序情况如下：
// "On" 2 个字母。
// "and" 3 个字母。
// "keep" 4 个字母，因为存在长度相同的其他单词，所以它们之间需要保留在原句子中的相对顺序。
// "calm" 4 个字母。
// "code" 4 个字母。

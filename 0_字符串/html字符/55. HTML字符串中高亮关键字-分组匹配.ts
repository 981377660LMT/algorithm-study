/**
 * @param {string} html
 * @param {string[]} keywords
 * 假设你在实现一个搜索建议。
   当输入关键词的时候，你需要在建议中高亮关键词，你如何做到？
   @description
   前缀树是搜索关键词 这里没必要用
   @summary
   replace正则匹配
 */
function highlightKeywords(html: string, keywords: string[]): string {
  // your code here
  const pattern = new RegExp(keywords.join('|'), 'ig')
  return html
    .split(/\s+/g)
    .map(word => {
      if (keywords.includes(word)) return `<em>${word}</em>`
      return word.replace(pattern, match => `<em>${match}</em>`).replace(/\<\/em\>\<em\>/g, '')
    })
    .join(' ')
}

// '<em>Hello</em> <em>Front</em>End Lovers')
console.log(highlightKeywords('Hello FrontEnd Lovers', ['Hello', 'Front', 'JavaScript']))
// 'Hello <em>FrontEnd</em> Lovers'
console.log(highlightKeywords('Hello FrontEnd Lovers', ['Front', 'FrontEnd', 'JavaScript']))
export {}

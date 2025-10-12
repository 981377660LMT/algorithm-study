const diff = require('./fastDiff.js')

{
  const text1 = 'The quick brown fox jumps over the lazy dog.'
  const text2 = 'The quick red fox leaps over the lazy cat.'

  const differences = diff(text1, text2, null, true)

  console.log('\n可视化结果 (HTML):')
  const htmlResult = prettyPrint(differences)
  console.log(htmlResult)
}

{
  const oldCode = `function hello() {
    console.log("world");
  }`
  const newCode = `function hello() {
    console.log("JavaScript");
    return true;
  }`

  const codeDiffs = diff(oldCode, newCode, null, true)
  console.log('\n代码比较的可视化结果:')
  console.log(prettyPrint(codeDiffs))
}

function prettyPrint(diffs) {
  let html = ''
  for (const [op, data] of diffs) {
    const text = data
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      // .replace(/\n/g, '&para;<br>')
      .replace(/\n/g, '<br>')
    switch (op) {
      case diff.INSERT: // 1
        html += `<ins style="background:#e6ffe6;">${text}</ins>`
        break
      case diff.DELETE: // -1
        html += `<del style="background:#ffe6e6;">${text}</del>`
        break
      case diff.EQUAL: // 0
        html += `<span>${text}</span>`
        break
    }
  }
  return html
}

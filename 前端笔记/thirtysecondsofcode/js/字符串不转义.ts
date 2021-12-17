// 相当于python中的 r
const path = `C:\web\index.html` // 'C:web.html'
console.log(path)

const r = String.raw

const unescapedPath = r`C:\web\index.html`
console.log(unescapedPath)
export {}

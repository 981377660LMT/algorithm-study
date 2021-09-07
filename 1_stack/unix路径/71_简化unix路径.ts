// 任意多个连续的斜杠（即，'//'）都被视为单个斜杠 '/'
const simplifyPath = (path: string) => {
  const stack: string[] = []
  const input = path.split(/\/+/).filter(s => s.trim().length)

  for (const s of input) {
    if (s === '.') continue
    if (s === '..') stack.pop()
    else stack.push(s)
  }

  return '/' + stack.join('/')
}

console.dir(simplifyPath('/home//foo/'), { depth: null })
console.dir(simplifyPath('/a/./b/../../c/'), { depth: null })

export {}

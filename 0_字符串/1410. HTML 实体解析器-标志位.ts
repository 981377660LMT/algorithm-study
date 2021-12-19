function entityParser2(text: string): string {
  const regs = [
    {
      reg: /&quot;/g,
      value: '"',
    },
    {
      reg: /&apos;/g,
      value: "'",
    },
    {
      reg: /&gt;/g,
      value: '>',
    },
    {
      reg: /&lt;/g,
      value: '<',
    },
    {
      reg: /&frasl;/g,
      value: '/',
    },
    {
      reg: /&amp;/g,
      value: '&',
    },
  ]

  for (let r of regs) {
    text = text.replace(r.reg, r.value)
  }

  return text
}

// isReplaceMode 这个想法很好
function entityParser(text: string): string {
  const replaceDict: Record<string, string> = {
    '&quot;': '"',
    '&apos;': "'",
    '&amp;': '&',
    '&gt;': '>',
    '&lt;': '<',
    '&frasl;': '/',
  }

  let isReplaceMode = false
  let keyword: string[] = []
  const res: string[] = []

  for (const char of text) {
    if (char === '&') {
      isReplaceMode = true
      keyword = ['&']
    } else if (isReplaceMode) {
      keyword.push(char)
      if (char === ';') {
        const word = keyword.join('')
        res.push(replaceDict[word] ?? word)
        isReplaceMode = false
      }
    } else res.push(char)
  }

  if (isReplaceMode) res.push(...keyword)
  return res.join('')
}

console.log(entityParser('&amp; is an HTML entity but &ambassador; is not.'))

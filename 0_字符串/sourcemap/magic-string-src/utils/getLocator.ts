/**
 * 根据偏移量找到源码对应位置.
 */
export default function getLocator(source: string) {
  const originalLines = source.split('\n')
  const lineOffsets: number[] = []

  for (let i = 0, pos = 0; i < originalLines.length; i++) {
    lineOffsets.push(pos)
    pos += originalLines[i].length + 1
  }

  return function locate(index: number) {
    let i = 0
    let j = lineOffsets.length
    while (i < j) {
      const m = (i + j) >> 1
      if (index < lineOffsets[m]) {
        j = m
      } else {
        i = m + 1
      }
    }
    const line = i - 1
    const column = index - lineOffsets[line]
    return { line, column }
  }
}

if (require.main === module) {
  const source = `
 The quick brown fox
 jumps over 
 the lazy dog.
 `

  const locate = getLocator(source)
  console.log(locate(0)) // { line: 0, column: 0 }
  console.log(locate(1)) // { line: 0, column: 1 }
  console.log(locate(19)) // { line: 0, column: 19 }
  console.log(locate(20)) // { line: 1, column: 0 }
  console.log(locate(21)) // { line: 1, column: 1 }
  console.log(locate(33)) // { line: 1, column: 13 }
}

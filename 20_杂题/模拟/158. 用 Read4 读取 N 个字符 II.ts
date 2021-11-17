/**
 * Definition for read4()
 * read4 = function(buf4: string[]): number {
 *     ...
 * };
 */

const solution = function (read4: any) {
  const fileBuf: any[] = []
  let readOffset = 0 // fileBuf游标
  let bufSize = 0 // fileBuf读到的字符数

  // 请实现一个方法使其能够读取 n 个字符。
  // 注意：你的 read 方法可能会被调用多次。
  return function (buf: string[], n: number): number {
    for (let i = 0; i < n; i++) {
      const nextChar = getNextCharFromFile()
      if (nextChar === 0) return i
      buf[i] = nextChar
    }

    return n
  }

  function getNextCharFromFile(): any {
    if (readOffset === bufSize) {
      bufSize = read4(fileBuf)
      readOffset = 0

      if (bufSize === 0) return 0
    }

    return fileBuf[readOffset++]
  }
}

export {}

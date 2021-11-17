// 给你一个文件，并且该文件只能通过给定的 read4 方法来读取，请实现一个方法使其能够读取 n 个字符。
// read4:存进buf，返回值为实际存了几个字母

const solution = function (read4: (tmp: any) => number) {
  // 实现一个方法使其能够读取 n 个字符。
  // 返回值为实际读取的字符个数
  return function (buf: string[], n: number): number {
    let offset = 0

    for (let i = 0; i < n; i += 4) {
      // 必须先开辟出4个空间读入到tmp
      const tmp: any[] = []
      const curLen = read4(tmp)
      for (let j = 0; j < curLen; j++) {
        buf[offset] = tmp[j]
        offset++
      }
    }

    return Math.min(offset, n)
  }
}

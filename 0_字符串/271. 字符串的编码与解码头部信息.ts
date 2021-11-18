/**
 * Encodes a list of strings to a single string.
 * 用前4位记录了每个数组的长度，每次get4个字符，然后用这四个字符获取字符长度即可。
 */
function encode(strs: string[]): string {
  const list: string[] = []

  for (const str of strs) {
    const len = str.length.toString().padStart(4, '0')
    list.push(len)
    list.push(str)
  }

  return list.join('')
}

/**
 * Decodes a single string to a list of strings.
 */
function decode(s: string): string[] {
  const res: string[] = []
  let index = 0

  while (index < s.length) {
    const strLength = parseInt(s.slice(index, index + 4), 10)
    res.push(s.slice(index + 4, index + 4 + strLength))
    index += 4 + strLength
  }

  return res
}

console.log(encode(['as', '12']))
console.log(decode('0002as000212'))
// 方法二：分块编码
// 这种方法基于 HTTP v1.1 使用的编码，
// 它不依赖于输入字符集，因此比方法一更具有通用性和有效性。
// 数据流被分成块，每个块前面都有其字节大小。

// 遍历字符串数组。
// 计算每个字符串的长度，并将长度大小转换为 4 个字节的字符串。
// 将长度信息的字符串添加到编码字符串的前面。
// 前面 4 个字节表示了编码字符串的长度。
// 后面跟这字符串本身。
// 返回编码后的字符串。

export {}

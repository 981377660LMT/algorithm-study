// const lengthToString = (str: string) => {
//   const len = str.length
//   const stringBuilder: number[] = []
//   for (let i = 3; ~i; i--) {
//     stringBuilder[3 - i] = (len >> (i * 8)) & 0xff
//   }
//   return stringBuilder.join('')
// }

// const stringToLength = (str: string) => {
//   const len = 0
//   const stringBuilder: number[] = []
//   for (let i = 3; ~i; i--) {
//     stringBuilder[3 - i] = (len >> (i * 8)) & 0xff
//   }
//   return len
// }

// console.log(lengthToString('asassdgfhh'))
/**
 * Encodes a list of strings to a single string.
 * 注意不能用join join会无视空字符串
 */
function encode(strs: string[]): string {
  const split = String.fromCharCode(256)
  let res = ''
  for (const str of strs) {
    res += str
    res += split
  }
  return res
}

/**
 * Decodes a single string to a list of strings.
 */
function decode(s: string): string[] {
  const split = String.fromCharCode(256)
  if (s === split) return ['']
  return s.split(split).slice(0, -1)
}

// 因为字符串可能会包含 256 个合法 ascii 字符中的任何字符，
// 所以您的算法必须要能够处理任何可能会出现的字符。

// 请勿使用 “类成员”、“全局变量” 或 “静态变量” 来存储这些状态，
// 您的编码和解码算法应该是非状态依赖的。(纯函数)

// 请不要依赖任何方法库，例如 eval 又或者是 serialize 之类的方法。
// 本题的宗旨是需要您自己实现 “编码” 和 “解码” 算法。

// ASCII中的0~31为控制字符；32~126为打印字符；127为Delete(删除)命令
// ASCII扩展字符——（为了适应更多字符）128~255
////////////////////////////////////////////////////////////////////
// 方法一：使用非 ASCII 码的分隔符
// 最简单的方法就是分隔符连接字符串
// console.log(String.fromCharCode(255)) // Ā

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
console.log(encode(['', '']))
console.log(['', '', 'a'].join('b'))
console.log('a'.split('a').slice(0, -1))
console.log(['', null, 'a', false, undefined].join('11')) // join会无视空字符串和null/undefined

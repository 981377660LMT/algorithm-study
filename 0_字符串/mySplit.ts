export {}

function mySplit(str: string, separator: string): string[] {
  const result: string[] = []

  // 特殊处理：如果 separator 是空字符串，直接展开
  if (separator === '') {
    return Array.from(str)
  }

  let start = 0
  let index = str.indexOf(separator)

  while (index !== -1) {
    // 截取从上一个切点到当前切点的内容
    result.push(str.substring(start, index))
    // 更新起始位置：跳过这个 separator 本身
    start = index + separator.length
    // 寻找下一个切点
    index = str.indexOf(separator, start)
  }

  // 别忘了最后一段（最后一个 separator 之后的部分，或者没有 separator 的情况）
  result.push(str.substring(start))

  return result
}

// 测试
console.log(mySplit('1101011', '0')) // ["11", "1", "11"]
console.log(mySplit('0011', '0')) // ["", "", "11"]
console.log(mySplit('a,b,c', ',')) // ["a", "b", "c"]

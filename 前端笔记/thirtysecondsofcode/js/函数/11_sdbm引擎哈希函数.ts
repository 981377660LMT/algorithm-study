console.log(sdbm('name')) // -3521204949)

// 本算法是由于在开源项目SDBM（一种简单的数据库引擎）中被应用而得名
// 将输入字符串散列为整数。
function sdbm(str: string): any {
  const arr = str.split('')
  return arr.reduce((hashCode, curValue) => {
    hashCode = curValue.codePointAt(0)! + (hashCode << 6) + (hashCode << 16) - hashCode
    return hashCode
  }, 0)
}

// Python replace() 方法把字符串中的 old（旧字符串） 替换成 new(新字符串)，
// 如果指定第三个参数max，则替换不超过 max 次。
const allIndexOf = function (str: string, searchElement: string) {
  if (searchElement === '') return [] // 否则会死循环
  const res: number[] = []
  let idx = str.indexOf(searchElement)
  while (idx !== -1) {
    res.push(idx)
    idx = str.indexOf(searchElement, idx + 1)
  }
  return res
}

// String.prototype.indexOf(searchString)
// console.log('asa'.indexOf('')) // 0

export { allIndexOf }

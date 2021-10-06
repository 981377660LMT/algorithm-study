const addPoint1 = (str: string, k: number) => {
  const res = str.split('')
  for (let i = str.length - 1 - k; i >= 0; i -= k) {
    res[i] += '.'
  }
  return res.join('')
}

console.log(addPoint1('1000000000', 3))
console.log(Number(10000000).toLocaleString('de'))
console.log('1000000000'.replace(/(\d)(?=(\d{3})+\b)/g, '$1.'))
// 词边界测试 \b 检查位置的一侧是否匹配 \w，而另一侧则不匹配 “\w”。

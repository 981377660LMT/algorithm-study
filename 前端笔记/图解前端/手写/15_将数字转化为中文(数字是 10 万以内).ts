// 比如数字12345，我们转化为一万二千三百四十五。 比如数字10002，我们转化为一万零二
function numToChinese(num: number) {
  const numStr = num.toString()
  const numMapper = ['零', '一', '二', '三', '四', '五', '六', '七', '八', '九']
  const unitMapper = ['', '', '十', '百', '千', '万']
  let res = ''

  for (let i = 0; i < numStr.length; i++) {
    const curNum = numStr[i] === '0' && res[res.length - 1] === '零' ? '' : numMapper[+numStr[i]]
    const curUnit = numStr[i] === '0' ? '' : unitMapper[numStr.length - i]
    res += curNum + curUnit
  }

  return res[res.length - 1] === '零' ? res.slice(0, -1) : res
}

console.log(numToChinese(12345))

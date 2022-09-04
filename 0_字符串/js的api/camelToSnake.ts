// camelcase 转 snakecase 驼峰命名转蛇形命名

// !如果变量名首字母是大写字母，则整个变量名都不做转换
// 将其他出现的大写字母转为下划线加小写字母
// 如果大写字母前已经有下划线，则不添加下划线只转小写
// 当出现连续的大写字母，只在其第一个和最后一个字母前添加下划线

function camelToSnake(sentence: string): string {
  const firstOrd = sentence.charCodeAt(0)
  if (firstOrd >= 65 && firstOrd <= 90) {
    return sentence
  }

  const regexp = /([A-Z]+)/g // 匹配连续的大写字母
  const res = sentence.replace(regexp, (match, _, index) => {
    let cand = ''
    if (match.length === 1) cand = `_${match}`
    else cand = `_${match.slice(0, -1)}_${match.slice(-1)}`
    if (index - 1 >= 0 && sentence[index - 1] === '_') cand = cand.slice(1)
    return cand
  })

  return res.toLowerCase()
}

if (require.main === module) {
  console.log(camelToSnake('AppleTree'))
  console.log(camelToSnake('appleHTMLElemnt__A'))
}

export { camelToSnake }

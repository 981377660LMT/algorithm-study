const fs = require('fs')

// 获取单个用例的所有行，此时 lines[0] 为第一行数据
const lines = fs
  .readFileSync(0)
  .toString()
  .trim()
  .split(/\r\n|\r|\n/)

for (let i = 1; i < lines.length; i++) {
  console.log(getAns(lines[i]))
}

function getAns(s) {
  if (s.length < 2) return 'Wrong'
  const first = s.charCodeAt(0)
  if (first < 65 || first > 122 || (first > 90 && first < 97)) return 'Wrong'

  const reg = new RegExp(/^[0-9]*$/)
  let num = 1
  for (let i = 1; i < s.length; i++) {
    if (!reg.test(s[i])) {
      let x = s.charCodeAt(i)
      if (x < 65 || x > 122 || (x > 90 && x < 97)) return 'Wrong'
    } else {
      num = 2
    }
  }

  return num == 2 ? 'Accept' : 'Wrong'
}

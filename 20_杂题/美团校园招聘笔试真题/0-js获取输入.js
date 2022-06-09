// 美团校园招聘笔试题

const { readFileSync } = require('fs')
const iter = readlines()
const input = () => iter.next().value
function* readlines(path = 0) {
  const lines = readFileSync(path)
    .toString()
    .trim()
    .split(/\r\n|\r|\n/)

  yield* lines
}

if (require.main === module) {
  console.log(input())
  console.log(input())
  console.log(input())
}

export { input }

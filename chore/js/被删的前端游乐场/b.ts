const genRandomStr = (length: number) => {
  const str = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let res = []
  for (let i = 0; i < length; i++) {
    res.push(str[~~(Math.random() * str.length)])
  }
  return res.join('')
}

const randomStr = Array.from({ length: 1e6 }, () => genRandomStr(50))
const searchStr = genRandomStr(5)
console.time('foo')
const res = []
for (let i = 0; i < randomStr.length; i++) {
  if (randomStr[i].toLowerCase().includes(searchStr)) {
    res.push(i)
  }
}
console.timeEnd('foo')

export {}

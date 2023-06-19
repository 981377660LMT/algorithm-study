// forin遍历对象比Object.keys()更快

const bigObj = {}
for (let i = 0; i < 1000000; i++) {
  bigObj[i] = i
}

console.time('object.keys')
const keys = Object.keys(bigObj)
for (let i = 0; i < keys.length; i++) {
  const key = keys[i]
  const value = bigObj[key]
}
console.timeEnd('object.keys')

console.time('for in keys')
for (const key in bigObj) {
  const value = bigObj[key]
}
console.timeEnd('for in keys')

// 但是for in 会遍历到原型链上的属性，所以需要使用hasOwnProperty()来判断是否是自身属性
console.time('for in keys')
for (const key in bigObj) {
  if (bigObj.hasOwnProperty(key)) {
    const value = bigObj[key]
  }
}
console.timeEnd('for in keys')

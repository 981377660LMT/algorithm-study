function isSameArray(a, b) {
  if (a.length !== b.length) return false
  for (let i = 0; i < a.length; ++i) {
    if (a[i] !== b[i]) return false
  }
  return true
}

function isSameArray2(a, b) {
  return a[0] === b[0] && a[1] === b[1] && a[2] === b[2]
}

const a = [1, 2, 3]
const b = [1, 2, 3]
console.time('isSameArray')
for (let i = 0; i < 1e7; ++i) {
  isSameArray(a, b)
}
console.timeEnd('isSameArray') // isSameArray: 41.443ms

console.time('isSameArray2')
for (let i = 0; i < 1e7; ++i) {
  isSameArray2(a, b)
}
console.timeEnd('isSameArray2') // isSameArray2: 19.535ms

function isSameObj(a, b) {
  const keys1 = Object.keys(a)
  const keys2 = Object.keys(b)
  if (keys1.length !== keys2.length) return false
  for (let i = 0; i < keys1.length; ++i) {
    const key = keys1[i]
    if (a[key] !== b[key]) return false
  }
  return true
}

function isSameObj2(o1, o2) {
  return o1.a === o2.a && o1.b === o2.b && o1.c === o2.c
}

const o1 = { a: 1, b: 2, c: 3 }
const o2 = { a: 1, b: 2, c: 3 }
console.time('isSameObj')
for (let i = 0; i < 1e7; ++i) {
  isSameObj(o1, o2)
}
console.timeEnd('isSameObj') // isSameObj: 630.697ms

console.time('isSameObj2')
for (let i = 0; i < 1e7; ++i) {
  isSameObj2(o1, o2)
}
console.timeEnd('isSameObj2') // isSameObj2: 10.656ms

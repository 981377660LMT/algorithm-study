import { List } from 'immutable'

const arr = List()

console.time('push')
for (let i = 0; i < 1e6; i++) {
  arr.push(i)
}
console.timeEnd('push')

const nums = []
console.time('push')
for (let i = 0; i < 1e6; i++) {
  nums.push(i)
}
console.timeEnd('push')

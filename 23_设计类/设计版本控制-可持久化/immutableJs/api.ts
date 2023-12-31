import Immutable, { List, Map } from 'immutable'

const arr = List()

console.time('push')
for (let i = 0; i < 1e6; i++) {
  arr.push(i)
}

arr.splice(0, 2)
arr.pu
console.log(arr)
console.timeEnd('push')

const nums = []
console.time('push')
for (let i = 0; i < 1e6; i++) {
  nums.push(i)
}
console.timeEnd('push')

let mp = Immutable.Map()
mp = mp.withMutations(m => {
  for (let i = 0; i < 5; i++) {
    m.set(i, i)
  }
})

console.dir(mp, { depth: null })

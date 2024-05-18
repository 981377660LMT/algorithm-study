// 解决前端精度问题的JS库-math.js

import * as math from 'mathjs'

const M = math.create(math.all, {
  number: 'BigNumber',
  precision: 20
})

const add = (a: number, b: number) => M.format(M.add(a, b), { precision: 16 })
console.log(add(0.1, 0.2))
console.log(math.sqrt(-4))
// true
console.log(math.isComplex(math.sqrt(-4)))
console.log(math.evaluate('sqrt(-4)'))

console.log(math.chain(3).add(4).multiply(2).done())

const _ = require('lodash')

// 将数组（array）拆分成多个 size 长度的区块
console.log(_.chunk(['a', 'b', 'c', ['d', 'e']], 2))

// bisectLeft
console.log(_.sortedIndex([1, 2, 3], 3))
// bisectRight
console.log(_.sortedLastIndex([1, 2, 3], 3))
// bisectLeft 带key
console.log(
  _.sortedIndexBy(
    [
      [1, 'a'],
      [2, 'a'],
      [3, 'd'],
    ],
    [2, 'as'],
    0
  )
)
// bisectRight 带key
console.log(
  _.sortedLastIndexBy(
    [
      [1, 'a'],
      [2, 'a'],
      [3, 'd'],
    ],
    [2, 'as'],
    0
  )
)

// zip
console.log(_.zip([1, 2], ['a', 'b']))

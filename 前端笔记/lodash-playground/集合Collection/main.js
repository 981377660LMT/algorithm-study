const _ = require('lodash')

// 计数
const counter = _.countBy([1, 2, 9, 9])
console.log(counter)

// groupBy神器
const groups = _.groupBy('11122233333', val => val)
console.log(groups)
// {
//   '1': [ '1', '1', '1' ],
//   '2': [ '2', '2', '2' ],
//   '3': [ '3', '3', '3', '3', '3' ]
// }

// 从collection（集合）中获得 n 个随机元素。
console.log(_.sampleSize([1, 2, 3], 2))
console.log(_.shuffle([2, 1]))

// 获取size  同 python 的 len
console.log(_.size({}))

// defer
// 推迟调用func，直到当前堆栈清理完毕
_.defer(text => console.log('defer' + text), 'haha')

// 记忆化 效果很差
// const obj = { a: 1, b: 2 }
// const values = _.memoize(
//   () => obj,
//   (...args) => args.join('_')
// )
// console.log(values)
// values.cache.clear()

// 深浅拷贝
const matrix = [
  [1, 2],
  [3, 4],
]
console.log(_.cloneDeep(matrix))

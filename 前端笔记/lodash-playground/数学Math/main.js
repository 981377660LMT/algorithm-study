const _ = require('lodash')

// 带key的max
const objects = [{ n: 1 }, { n: 2 }]
console.log(_.maxBy(objects, obj => obj.n))

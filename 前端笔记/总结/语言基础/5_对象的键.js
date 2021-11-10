// example 1
var a = {},
  b = '123',
  c = 123
a[b] = 'b'
a[c] = 'c'
console.log(a[b])

// ---------------------
// example 2
var a = {},
  b = Symbol('123'),
  c = Symbol('123')
a[b] = 'b'
a[c] = 'c'
console.log(a[b])

// ---------------------
// example 3
var a = {},
  b = { key: '123' },
  c = { key: '456' }
a[b] = 'b'
a[c] = 'c'
console.log(a[b])

var a = {},
  b = function () {
    return 1
  },
  c = function () {
    return 2
  }
a[b] = 'b'
a[c] = 'c'
console.log(a[b])

// 第一题，JS里对象的key都是字符串或者symbol，其他类型的都会被转换成字符串，所以b和c都是'123'，指向同一个属性，所以会被覆盖，输出 c
// 第二题，symbol可以作为对象的key且任何一个symbol的值都是唯一的不会被覆盖，所以输出b
// 第三题，对象当作key时，会调用toString方法转换成字符串，所以两个key都是[object Object]，所以输出c
// 第四题，function也会调用toString方法，所以结果是b

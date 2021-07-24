// map的两种构造方式
const m = new Map()
const m2 = new Map([['key', 'value']])

// 增
m.set('z', 'aa')
m.set('gg', 'aa')
m.set('a', 'aa')
m.set('b', 'bb')

// 删
m.delete('b')
// m.clear();

// 改
m.set('a', 'aaa')

// 查
console.log(m.get('a'))

// object转map
console.log(new Map(Object.entries({ foo: 'bar' })))
console.log(new Map([[2, 'value']]))

// js Map和java 的map不一样，内部是按顺序存储，可以用来做LRU
console.log(m.keys().next())

// 删除前n个值
// or (let i = 0; i < 3; i ++) {
//   aa.delete(aa.keys().next());
// }

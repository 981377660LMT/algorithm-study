export {}

let sem = { name: 'Semlinker' }

let map = new Map()
map.set(sem, '全栈修仙之路')
sem = null as any // 覆盖引用

// sem被存储在map中
// 我们可以使用map.keys()来获取它
console.log(map.keys())
// 我们使用对象作为常规 Map 的键，那么当 Map 存在时，该对象也将存在。它会占用内存，并且不会被垃圾回收机制回收。
// WeakMap 的 key 是不可枚举的 (没有方法能给出所有的 key)。如果key 是可枚举的话，其列表将会受垃圾回收机制的影响，从而得到不确定的结果。

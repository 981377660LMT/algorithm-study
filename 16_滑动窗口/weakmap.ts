let john = { name: 'John' }

let weakMap = new WeakMap()
weakMap.set(john, 'secret documents')
// 如果 john 消失，secret documents 将会被自动清除
// @ts-ignore
john = null // 覆盖引用
// john 被从内存中删除了！
/// ////////////////////////////////////////////////////////////
let cmnx = { name: 'cmnx' }
let map = new Map()
map.set(cmnx, '...')

// @ts-ignore
cmnx = null
// cmnx这个key对应的value无法被垃圾回收
console.log(map.keys())
// cmnx仍然存在于内存中

export {}

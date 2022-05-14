let map = new Map()

// object转map
map = new Map(Object.entries({ foo: 'bar' }))
map = new Map([[2, 'value']])

// js 的 map 是 LinkedHashMap
console.log(map.keys().next())

// 删除前3个值
for (let i = 0; i < 3; i++) {
  map.delete(map.keys().next())
}

const a = [['name', 'cmnx']]
const map1 = new Map(a)

console.log(map1)

const map2 = new Map(map1)

map2.delete('name')
console.log(map1, map2)

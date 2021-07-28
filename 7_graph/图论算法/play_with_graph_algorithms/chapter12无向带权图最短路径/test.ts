const obj = {
  a: { b: [1, 2, 3] },
  b() {
    console.log(1)
  },
}

console.dir(Object.entries(obj), { depth: null })

let map = new Map(Object.entries(obj))
map = new Map(map)
const copy = Object.fromEntries(map.entries())

console.log(copy)
obj.a.b.push(4)
console.log(obj, copy)

// map无法深拷贝
export {}

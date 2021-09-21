const w = new WeakMap()
const a = [1, 2]
w.set(a, 1)
console.log(w.get(a))
export {}

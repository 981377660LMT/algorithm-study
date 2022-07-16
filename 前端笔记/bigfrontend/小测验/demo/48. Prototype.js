function Foo() {}

const record = Foo.prototype // {}

Foo.prototype.bar = 1
const a = new Foo()
console.log(a.bar)

Foo.prototype.bar = 2
const b = new Foo()
console.log(a.bar)
console.log(b.bar)

Foo.prototype = { bar: 3 }

const c = new Foo()
// record.bar = 987
console.log(a.bar) // 2 还是指向record
console.log(b.bar) // 2 还是指向record
console.log(c.bar) // 指向{ bar: 3 }
// 1
// 2
// 2
// 2
// 2
// 3

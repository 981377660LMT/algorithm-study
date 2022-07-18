const obj = Object.create({ name: 'foo' })
const proto = Object.getPrototypeOf(obj)
Object.setPrototypeOf(proto, null)
console.log(obj)

function Person() {}
var person = new Person()
console.log(person.constructor === Person) // true

const obj = Object.create({ name: 'foo' })
const proto = Object.getPrototypeOf(obj)
Object.setPrototypeOf(proto, null)

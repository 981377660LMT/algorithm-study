const p1 = Promise.resolve(1) // p1 is 1
const p2 = new Promise(resolve => resolve(p1)) // p2 is a pending promise
const p3 = Promise.resolve(p1) // p3 is 1 (p3 is essentially Promise.resolve(1) )
const p4 = p2.then(() => new Promise(resolve => resolve(p3))) // p4 is a pending promise since .then() returns a new promise
const p5 = p4.then(() => p4) //  p5 is a pending promise since  .then() returns a new promise

console.log(p1 == p2)
console.log(p1 == p3)
console.log(p3 == p4)
console.log(p4 == p5)
// false
// true
// false
// false

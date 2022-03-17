const generator = function* (x, y) {
  let result = yield x + y
  return result
}
const gen = generator(1, 2)
console.log(gen.next())
console.log(gen.next(1))
console.log(gen.next(1))
console.log(gen.return(1))
console.log(gen.return(1))

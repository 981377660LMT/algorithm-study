const cycleGenerator = function* <T>(arr: T[]) {
  let i = 0
  while (true) {
    yield arr[i++ % arr.length]
  }
}
const binaryCycle = cycleGenerator([0, 1])
binaryCycle.next() // { value: 0, done: false }
binaryCycle.next() // { value: 1, done: false }
binaryCycle.next() // { value: 0, done: false }
console.log(binaryCycle.next()) // { value: 1, done: false }

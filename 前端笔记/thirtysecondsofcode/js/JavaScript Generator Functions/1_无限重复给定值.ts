const repeater = repeatGenerator(5)
repeater.next() // { value: 5, done: false }
repeater.next() // { value: 5, done: false }
repeater.next(4) // { value: 4, done: false }
repeater.next() // { value: 4, done: false }

function* repeatGenerator(num: number): Generator<number, number, number> {
  let state = num
  while (true) {
    const newState = yield state
    if (newState != undefined) state = newState
  }
}

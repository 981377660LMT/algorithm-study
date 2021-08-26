function* gen(letter: '1' | '2' | '3') {
  switch (letter) {
    case '1':
      yield* ['a', 'b', 'c']
      break
    case '2':
      yield* ['d', 'e']
      break
    case '3':
      yield* ['f', 'g', 'h']
      break
    default:
      throw new Error('invalid input')
  }
}

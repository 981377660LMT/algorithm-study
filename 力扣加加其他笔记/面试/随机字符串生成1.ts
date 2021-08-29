function* getDigit(count: number) {
  for (let i = 0; i < count; i++) {
    yield ~~(Math.random() * 10)
  }
}

function* getLowercase(count: number) {
  const BASE = 97
  for (let i = 0; i < count; i++) {
    const offset = ~~(Math.random() * 26)
    yield String.fromCodePoint(BASE + offset)
  }
}

function* getUppercase(count: number) {
  const BASE = 65
  for (let i = 0; i < count; i++) {
    const offset = ~~(Math.random() * 26)
    yield String.fromCodePoint(BASE + offset)
  }
}

function* getId(length: number) {
  const Type = 3
  for (let i = 0; i < length; i++) {
    const random = ~~(Math.random() * Type)
    switch (random) {
      case 0:
        yield* getDigit(1)
        break
      case 1:
        yield* getLowercase(1)
        break
      case 2:
        yield* getUppercase(1)
        break
      default:
        throw new Error('invalid randomInt')
    }
  }
  // yield* [yield* getLowercase(length), yield* getUppercase(length), yield* getDigit(length)]
}

const genId = getId(4)
console.log([...genId].join(''))

export {}

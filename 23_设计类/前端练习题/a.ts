function* fibGenerator(): Generator<number, any, number> {
  // 0, 1, 1, 2, 3, 5, 8, 13 。
  let a = 0
  let b = 1
  while (true) {
    yield a
    const c = a + b
    a = b
    b = c
  }
}

/**
 * const gen = fibGenerator();
 * gen.next().value; // 0
 * gen.next().value; // 1
 */
export {}

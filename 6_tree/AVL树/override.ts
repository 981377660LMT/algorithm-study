class Animal {
  move(meters: number): number {
    return 1
  }
}

class Snake extends Animal {
  override move(meters: number, time: number): number {
    return 2
  }
}
// 4.3
// When a method is marked with override,
// TypeScript will always make sure that a method with the same name exists in a the base class.

// ts的override 必须类型要一样 (子类型)

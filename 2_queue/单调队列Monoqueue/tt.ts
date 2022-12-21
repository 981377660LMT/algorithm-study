/* eslint-disable @typescript-eslint/no-explicit-any */
// 1. 想限定容器每个元素的类型,并且取

class Container<E extends { value: V }, V = E['value']> {
  private readonly _v: V
  private readonly _compare: (a: V, b: V) => number

  constructor(obj: E, compare: (a: V, b: V) => number = (a: any, b: any) => a - b) {
    this._v = obj.value
    this._compare = compare
  }

  hello(): void {
    console.log(this._v)
    console.log(this._compare)
  }
}

const stringFoo = new Container({ value: 'foo', index: 1 }, (a, b) => a.localeCompare(b))
const numberFoo = new Container({ value: 2 })

export {}

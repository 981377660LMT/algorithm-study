import { debounce } from '../1_实现一个debounce装饰器'

jest.useFakeTimers()

let a: any
let mockFunc: jest.Mock
beforeEach(() => {
  mockFunc = jest.fn()
  class Test {
    @debounce(1000)
    sayHi() {
      mockFunc()
    }
  }
  a = new Test()
})

describe('debounce:', () => {
  test('debounced function should be called after the delay time', () => {
    a.sayHi()
    expect(mockFunc).toHaveBeenCalledTimes(0)
    jest.advanceTimersByTime(1000)
    expect(mockFunc).toHaveBeenCalledTimes(1)
  })

  test('debounced function should not be called before the delay time', () => {
    a.sayHi()
    expect(mockFunc).toHaveBeenCalledTimes(0)
    let count = 100
    while (count--) {
      a.sayHi()
    }
    expect(mockFunc).toHaveBeenCalledTimes(0)

    count = 100
    while (count--) {
      jest.advanceTimersByTime(999)
      a.sayHi()
    }
    expect(mockFunc).toHaveBeenCalledTimes(0)
  })
})

/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * !Slide Window Aggrigation
 * Api:
 *   1. append value to tail,O(1).
 *   2. pop value from head,O(1).
 *   3. query aggrigated value in window,O(1).
 *
 * @param e unit element
 * @param op merge function
 */
function slidingWindowAggrigation<E>(e: () => E, op: (a: E, b: E) => E) {
  const stack0: E[] = []
  const stack1: E[] = []
  const stack2: E[] = []
  const stack3: E[] = []
  let e0 = e()
  let e1 = e()
  let _size = 0

  function append(value: E): void {
    if (!stack0.length) {
      _push0(value)
      _transfer()
    } else {
      _push1(value)
    }

    _size++
  }

  function popLeft(): void {
    if (!_size) return
    if (!stack0.length) _transfer()
    stack0.pop()
    stack2.pop()
    e0 = stack2.length ? stack2[stack2.length - 1] : e()
    _size--
  }

  function query(): E {
    return op(e0, e1)
  }

  return {
    append,
    popLeft,
    query,
    get size() {
      return _size
    }
  }

  function _push0(value: E): void {
    stack0.push(value)
    e0 = op(value, e0)
    stack2.push(e0)
  }

  function _push1(value: E): void {
    stack1.push(value)
    e1 = op(e1, value)
    stack3.push(e1)
  }

  function _transfer(): void {
    while (stack1.length) {
      _push0(stack1.pop()!)
    }
    while (stack3.length) stack3.pop()
    e1 = e()
  }
}

export { slidingWindowAggrigation }

if (require.main === module) {
  const gcd = (a: number, b: number): number => (b === 0 ? a : gcd(b, a % b))
  const window = slidingWindowAggrigation(() => 0, gcd)
  console.log(window.query())
  window.append(1)
  console.log(window.query())
  window.append(2)
  console.log(window.query())
  window.append(4)
  console.log(window.query())
  window.popLeft()
  console.log(window.size, 'size')
  console.log(window.query())
  window.append(8)
  console.log(window.query())
  window.popLeft()
  console.log(window.query())
}

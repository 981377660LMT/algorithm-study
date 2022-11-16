/* eslint-disable no-inner-declarations */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * !Sliding Window Aggregation
 *
 * Api:
 *   1. append value to tail,O(1).
 *   2. pop value from head,O(1).
 *   3. query aggregated value in window,O(1).
 *
 * @param e unit element
 * @param op merge function
 */
function slidingWindowAggregation<E>(e: () => E, op: (a: E, b: E) => E) {
  const _stack0: E[] = []
  const _stack1: E[] = []
  const _stack2: E[] = []
  const _stack3: E[] = []
  let _e0 = e()
  let _e1 = e()
  let _size = 0

  function append(value: E): void {
    if (!_stack0.length) {
      _push0(value)
      _transfer()
    } else {
      _push1(value)
    }

    _size++
  }

  function popLeft(): void {
    if (!_size) return
    if (!_stack0.length) _transfer()
    _stack0.pop()
    _stack2.pop()
    _e0 = _stack2.length ? _stack2[_stack2.length - 1] : e()
    _size--
  }

  function query(): E {
    return op(_e0, _e1)
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
    _stack0.push(value)
    _e0 = op(value, _e0)
    _stack2.push(_e0)
  }

  function _push1(value: E): void {
    _stack1.push(value)
    _e1 = op(_e1, value)
    _stack3.push(_e1)
  }

  function _transfer(): void {
    while (_stack1.length) {
      _push0(_stack1.pop()!)
    }
    while (_stack3.length) _stack3.pop()
    _e1 = e()
  }
}

export { slidingWindowAggregation }

if (require.main === module) {
  const gcd = (a: number, b: number): number => {
    if (Number.isNaN(a) || Number.isNaN(b)) return NaN
    return b === 0 ? a : gcd(b, a % b)
  }

  const window = slidingWindowAggregation(() => 0, gcd)
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

  // 滑动窗口最大值
  function maxSlidingWindow(nums: number[], k: number): number[] {
    const maxWindow = slidingWindowAggregation(() => -2e15, Math.max)
    const res: number[] = []

    for (let i = 0; i < nums.length; i++) {
      maxWindow.append(nums[i])
      if (i >= k) {
        maxWindow.popLeft()
      }
      if (i >= k - 1) {
        res.push(maxWindow.query())
      }
    }

    return res
  }
}

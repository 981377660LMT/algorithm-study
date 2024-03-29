/* eslint-disable max-len */

type CompareFunc<T> = (a: T, b: T) => number
type MutableArrayLike<T> = {
  length: number
  [n: number]: T
}

/**
 * `ZRender`库中的`TimSort`原地排序.
 * @see {@link https://github.com/ecomfe/zrender/blob/master/src/core/timsort.ts}
 */
function timSort<T = number>(array: MutableArrayLike<T>, compare: CompareFunc<T> = (a: any, b: any) => a - b, lo = 0, hi = array.length): void {
  let remaining = hi - lo

  if (remaining < 2) {
    return
  }

  let runLength = 0

  if (remaining < DEFAULT_MIN_MERGE) {
    runLength = makeAscendingRun<T>(array, lo, hi, compare)
    binaryInsertionSort<T>(array, lo, hi, lo + runLength, compare)
    return
  }

  let ts = _TimSort<T>(array, compare)

  let minRun = minRunLength(remaining)

  do {
    runLength = makeAscendingRun<T>(array, lo, hi, compare)
    if (runLength < minRun) {
      let force = remaining
      if (force > minRun) {
        force = minRun
      }

      binaryInsertionSort<T>(array, lo, lo + force, lo + runLength, compare)
      runLength = force
    }

    ts.pushRun(lo, runLength)
    ts.mergeRuns()

    remaining -= runLength
    lo += runLength
  } while (remaining !== 0)

  ts.forceMergeRuns()
}

const DEFAULT_MIN_MERGE = 32

const DEFAULT_MIN_GALLOPING = 7

function minRunLength(n: number): number {
  let r = 0

  while (n >= DEFAULT_MIN_MERGE) {
    r |= n & 1
    n >>= 1
  }

  return n + r
}

function makeAscendingRun<T>(array: MutableArrayLike<T>, lo: number, hi: number, compare: CompareFunc<T>) {
  let runHi = lo + 1

  if (runHi === hi) {
    return 1
  }

  if (compare(array[runHi++], array[lo]) < 0) {
    while (runHi < hi && compare(array[runHi], array[runHi - 1]) < 0) {
      runHi++
    }

    reverseRun<T>(array, lo, runHi)
  } else {
    while (runHi < hi && compare(array[runHi], array[runHi - 1]) >= 0) {
      runHi++
    }
  }

  return runHi - lo
}

function reverseRun<T>(array: MutableArrayLike<T>, lo: number, hi: number) {
  hi--

  while (lo < hi) {
    let t = array[lo]
    array[lo++] = array[hi]
    array[hi--] = t
  }
}

function binaryInsertionSort<T>(array: MutableArrayLike<T>, lo: number, hi: number, start: number, compare: CompareFunc<T>) {
  if (start === lo) {
    start++
  }

  for (; start < hi; start++) {
    let pivot = array[start]

    let left = lo
    let right = start
    let mid: number

    while (left < right) {
      mid = (left + right) >>> 1

      if (compare(pivot, array[mid]) < 0) {
        right = mid
      } else {
        left = mid + 1
      }
    }

    let n = start - left

    switch (n) {
      case 3:
        array[left + 3] = array[left + 2]

      // eslint-disable-next-line no-fallthrough
      case 2:
        array[left + 2] = array[left + 1]

      // eslint-disable-next-line no-fallthrough
      case 1:
        array[left + 1] = array[left]
        break
      default:
        while (n > 0) {
          array[left + n] = array[left + n - 1]
          n--
        }
    }

    array[left] = pivot
  }
}

function gallopLeft<T>(value: T, array: MutableArrayLike<T>, start: number, length: number, hint: number, compare: CompareFunc<T>) {
  let lastOffset = 0
  let maxOffset = 0
  let offset = 1

  if (compare(value, array[start + hint]) > 0) {
    maxOffset = length - hint

    while (offset < maxOffset && compare(value, array[start + hint + offset]) > 0) {
      lastOffset = offset
      offset = (offset << 1) + 1

      if (offset <= 0) {
        offset = maxOffset
      }
    }

    if (offset > maxOffset) {
      offset = maxOffset
    }

    lastOffset += hint
    offset += hint
  } else {
    maxOffset = hint + 1
    while (offset < maxOffset && compare(value, array[start + hint - offset]) <= 0) {
      lastOffset = offset
      offset = (offset << 1) + 1

      if (offset <= 0) {
        offset = maxOffset
      }
    }
    if (offset > maxOffset) {
      offset = maxOffset
    }

    let tmp = lastOffset
    lastOffset = hint - offset
    offset = hint - tmp
  }

  lastOffset++
  while (lastOffset < offset) {
    let m = lastOffset + ((offset - lastOffset) >>> 1)

    if (compare(value, array[start + m]) > 0) {
      lastOffset = m + 1
    } else {
      offset = m
    }
  }
  return offset
}

function gallopRight<T>(value: T, array: MutableArrayLike<T>, start: number, length: number, hint: number, compare: CompareFunc<T>) {
  let lastOffset = 0
  let maxOffset = 0
  let offset = 1

  if (compare(value, array[start + hint]) < 0) {
    maxOffset = hint + 1

    while (offset < maxOffset && compare(value, array[start + hint - offset]) < 0) {
      lastOffset = offset
      offset = (offset << 1) + 1

      if (offset <= 0) {
        offset = maxOffset
      }
    }

    if (offset > maxOffset) {
      offset = maxOffset
    }

    let tmp = lastOffset
    lastOffset = hint - offset
    offset = hint - tmp
  } else {
    maxOffset = length - hint

    while (offset < maxOffset && compare(value, array[start + hint + offset]) >= 0) {
      lastOffset = offset
      offset = (offset << 1) + 1

      if (offset <= 0) {
        offset = maxOffset
      }
    }

    if (offset > maxOffset) {
      offset = maxOffset
    }

    lastOffset += hint
    offset += hint
  }

  lastOffset++

  while (lastOffset < offset) {
    let m = lastOffset + ((offset - lastOffset) >>> 1)

    if (compare(value, array[start + m]) < 0) {
      offset = m
    } else {
      lastOffset = m + 1
    }
  }

  return offset
}

function _TimSort<T>(array: MutableArrayLike<T>, compare: CompareFunc<T>) {
  let minGallop = DEFAULT_MIN_GALLOPING

  let runStart: number[]
  let runLength: number[]
  let stackSize = 0

  let tmp: T[] = []

  runStart = []
  runLength = []

  function pushRun(_runStart: number, _runLength: number) {
    runStart[stackSize] = _runStart
    runLength[stackSize] = _runLength
    stackSize += 1
  }

  function mergeRuns() {
    while (stackSize > 1) {
      let n = stackSize - 2

      if ((n >= 1 && runLength[n - 1] <= runLength[n] + runLength[n + 1]) || (n >= 2 && runLength[n - 2] <= runLength[n] + runLength[n - 1])) {
        if (runLength[n - 1] < runLength[n + 1]) {
          n--
        }
      } else if (runLength[n] > runLength[n + 1]) {
        break
      }
      mergeAt(n)
    }
  }

  function forceMergeRuns() {
    while (stackSize > 1) {
      let n = stackSize - 2

      if (n > 0 && runLength[n - 1] < runLength[n + 1]) {
        n--
      }

      mergeAt(n)
    }
  }

  function mergeAt(i: number) {
    let start1 = runStart[i]
    let length1 = runLength[i]
    let start2 = runStart[i + 1]
    let length2 = runLength[i + 1]

    runLength[i] = length1 + length2

    if (i === stackSize - 3) {
      runStart[i + 1] = runStart[i + 2]
      runLength[i + 1] = runLength[i + 2]
    }

    stackSize--

    let k = gallopRight<T>(array[start2], array, start1, length1, 0, compare)
    start1 += k
    length1 -= k

    if (length1 === 0) {
      return
    }

    length2 = gallopLeft<T>(array[start1 + length1 - 1], array, start2, length2, length2 - 1, compare)

    if (length2 === 0) {
      return
    }

    if (length1 <= length2) {
      mergeLow(start1, length1, start2, length2)
    } else {
      mergeHigh(start1, length1, start2, length2)
    }
  }

  function mergeLow(start1: number, length1: number, start2: number, length2: number) {
    let i = 0

    for (i = 0; i < length1; i++) {
      tmp[i] = array[start1 + i]
    }

    let cursor1 = 0
    let cursor2 = start2
    let dest = start1

    array[dest++] = array[cursor2++]

    if (--length2 === 0) {
      for (i = 0; i < length1; i++) {
        array[dest + i] = tmp[cursor1 + i]
      }
      return
    }

    if (length1 === 1) {
      for (i = 0; i < length2; i++) {
        array[dest + i] = array[cursor2 + i]
      }
      array[dest + length2] = tmp[cursor1]
      return
    }

    let _minGallop = minGallop
    let count1
    let count2
    let exit

    while (1) {
      count1 = 0
      count2 = 0
      exit = false

      do {
        if (compare(array[cursor2], tmp[cursor1]) < 0) {
          array[dest++] = array[cursor2++]
          count2++
          count1 = 0

          if (--length2 === 0) {
            exit = true
            break
          }
        } else {
          array[dest++] = tmp[cursor1++]
          count1++
          count2 = 0
          if (--length1 === 1) {
            exit = true
            break
          }
        }
      } while ((count1 | count2) < _minGallop)

      if (exit) {
        break
      }

      do {
        count1 = gallopRight<T>(array[cursor2], tmp, cursor1, length1, 0, compare)

        if (count1 !== 0) {
          for (i = 0; i < count1; i++) {
            array[dest + i] = tmp[cursor1 + i]
          }

          dest += count1
          cursor1 += count1
          length1 -= count1
          if (length1 <= 1) {
            exit = true
            break
          }
        }

        array[dest++] = array[cursor2++]

        if (--length2 === 0) {
          exit = true
          break
        }

        count2 = gallopLeft<T>(tmp[cursor1], array, cursor2, length2, 0, compare)

        if (count2 !== 0) {
          for (i = 0; i < count2; i++) {
            array[dest + i] = array[cursor2 + i]
          }

          dest += count2
          cursor2 += count2
          length2 -= count2

          if (length2 === 0) {
            exit = true
            break
          }
        }
        array[dest++] = tmp[cursor1++]

        if (--length1 === 1) {
          exit = true
          break
        }

        _minGallop--
      } while (count1 >= DEFAULT_MIN_GALLOPING || count2 >= DEFAULT_MIN_GALLOPING)

      if (exit) {
        break
      }

      if (_minGallop < 0) {
        _minGallop = 0
      }

      _minGallop += 2
    }

    minGallop = _minGallop

    minGallop < 1 && (minGallop = 1)

    if (length1 === 1) {
      for (i = 0; i < length2; i++) {
        array[dest + i] = array[cursor2 + i]
      }
      array[dest + length2] = tmp[cursor1]
    } else if (length1 === 0) {
      throw new Error()
    } else {
      for (i = 0; i < length1; i++) {
        array[dest + i] = tmp[cursor1 + i]
      }
    }
  }

  function mergeHigh(start1: number, length1: number, start2: number, length2: number) {
    let i = 0

    for (i = 0; i < length2; i++) {
      tmp[i] = array[start2 + i]
    }

    let cursor1 = start1 + length1 - 1
    let cursor2 = length2 - 1
    let dest = start2 + length2 - 1
    let customCursor = 0
    let customDest = 0

    array[dest--] = array[cursor1--]

    if (--length1 === 0) {
      customCursor = dest - (length2 - 1)

      for (i = 0; i < length2; i++) {
        array[customCursor + i] = tmp[i]
      }

      return
    }

    if (length2 === 1) {
      dest -= length1
      cursor1 -= length1
      customDest = dest + 1
      customCursor = cursor1 + 1

      for (i = length1 - 1; i >= 0; i--) {
        array[customDest + i] = array[customCursor + i]
      }

      array[dest] = tmp[cursor2]
      return
    }

    let _minGallop = minGallop

    while (true) {
      let count1 = 0
      let count2 = 0
      let exit = false

      do {
        if (compare(tmp[cursor2], array[cursor1]) < 0) {
          array[dest--] = array[cursor1--]
          count1++
          count2 = 0
          if (--length1 === 0) {
            exit = true
            break
          }
        } else {
          array[dest--] = tmp[cursor2--]
          count2++
          count1 = 0
          if (--length2 === 1) {
            exit = true
            break
          }
        }
      } while ((count1 | count2) < _minGallop)

      if (exit) {
        break
      }

      do {
        count1 = length1 - gallopRight<T>(tmp[cursor2], array, start1, length1, length1 - 1, compare)

        if (count1 !== 0) {
          dest -= count1
          cursor1 -= count1
          length1 -= count1
          customDest = dest + 1
          customCursor = cursor1 + 1

          for (i = count1 - 1; i >= 0; i--) {
            array[customDest + i] = array[customCursor + i]
          }

          if (length1 === 0) {
            exit = true
            break
          }
        }

        array[dest--] = tmp[cursor2--]

        if (--length2 === 1) {
          exit = true
          break
        }

        count2 = length2 - gallopLeft<T>(array[cursor1], tmp, 0, length2, length2 - 1, compare)

        if (count2 !== 0) {
          dest -= count2
          cursor2 -= count2
          length2 -= count2
          customDest = dest + 1
          customCursor = cursor2 + 1

          for (i = 0; i < count2; i++) {
            array[customDest + i] = tmp[customCursor + i]
          }

          if (length2 <= 1) {
            exit = true
            break
          }
        }

        array[dest--] = array[cursor1--]

        if (--length1 === 0) {
          exit = true
          break
        }

        _minGallop--
      } while (count1 >= DEFAULT_MIN_GALLOPING || count2 >= DEFAULT_MIN_GALLOPING)

      if (exit) {
        break
      }

      if (_minGallop < 0) {
        _minGallop = 0
      }

      _minGallop += 2
    }

    minGallop = _minGallop

    if (minGallop < 1) {
      minGallop = 1
    }

    if (length2 === 1) {
      dest -= length1
      cursor1 -= length1
      customDest = dest + 1
      customCursor = cursor1 + 1

      for (i = length1 - 1; i >= 0; i--) {
        array[customDest + i] = array[customCursor + i]
      }

      array[dest] = tmp[cursor2]
    } else if (length2 === 0) {
      throw new Error()
      // throw new Error('mergeHigh preconditions were not respected');
    } else {
      customCursor = dest - (length2 - 1)
      for (i = 0; i < length2; i++) {
        array[customCursor + i] = tmp[i]
      }
    }
  }

  return {
    mergeRuns,
    forceMergeRuns,
    pushRun
  }
}

export { timSort }

if (require.main === module) {
  const arr1 = Array(1e6)
  const arr2 = Array(1e6)
  for (let i = 0; i < arr1.length; i++) {
    arr1[i] = (Math.random() * 1e9) | 0
    arr2[i] = arr1[i]
  }

  console.time('sort2')
  arr1.sort((a, b) => b - a)
  console.timeEnd('sort2')
  console.time('sort1')
  timSort(arr2, (a, b) => b - a)
  console.timeEnd('sort1')

  // eslint-disable-next-line no-inner-declarations
  function sumImbalanceNumbers(N: number[]): number {
    const nums = new Uint16Array(N)
    const copy = nums.slice()
    let res = 0
    for (let i = 0; i < nums.length; i++) {
      for (let j = i; j < nums.length; j++) {
        const sub = nums.subarray(i, j + 1)
        timSort(sub)
        for (let k = 0; k < sub.length - 1; k++) res += +(sub[k + 1] - sub[k] > 1)
        nums.set(copy.subarray(i, j + 1), i)
      }
    }
    return res
  }
}

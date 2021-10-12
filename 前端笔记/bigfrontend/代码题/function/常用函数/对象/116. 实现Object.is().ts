// Object.is() 和===基于一致，除了以下情况：

// Object.is(0, -0) // false
// 0 === -0 // true

// Object.is(NaN, NaN) // true
// NaN === NaN // false

/**
 * @param {any} a
 * @param {any} b
 * @return {boolean}
 */
function is(a: any, b: any): boolean {
  if (typeof a === 'number' && typeof b === 'number') {
    // NaN, NaN 要相等
    if (Number.isNaN(a) && Number.isNaN(b)) {
      return true
    }

    // 0, -0 要不相等
    if (a === 0 && b === 0 && 1 / a !== 1 / b) {
      return false
    }
  }

  return a === b
}

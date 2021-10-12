/**
 * @param {any[]} arr
 * @returns {?} - sorry no type hint for this
 * 对arr执行的所有操作，都必须反映在原来的数组中
 */
function wrap<T>(arr: T[]): T[] {
  return new Proxy(arr, {
    get(target: T[], prop: PropertyKey) {
      if (isNumerical(prop)) prop = normalizeIndex(Number(prop), target.length)
      return Reflect.get(target, prop)
    },
    set(target: T[], prop: PropertyKey, value: any) {
      if (isNumerical(prop)) {
        prop = normalizeIndex(Number(prop), target.length)
        assert(prop)
      }
      return Reflect.set(target, prop, value)
    },
  })
}

function assert(index: number) {
  if (index < 0) throw new Error('incorrect index')
}

function normalizeIndex(index: number, length: number) {
  return index >= 0 ? index : length + index
}

function isNumerical(prop: any) {
  return typeof prop === 'string' && !isNaN(parseFloat(prop)) && isFinite(parseFloat(prop))
}

if (require.main === module) {
  const originalArr = [1, 2, 3]
  const arr = wrap(originalArr)

  // 访问的key 是 string
  console.log(arr[3]) // undefined
  console.log(arr[-1]) // 3

  // 访问Symbol.iterator 和 length 属性
  for (const num of arr) {
    console.log(num)
  }

  // 访问length 属性
  arr.forEach(v => console.log(v))
}

export { wrap }

// Cannot convert a Symbol value to a string
// 注意:Symbol转数字和字符串/比较都会抛出错误

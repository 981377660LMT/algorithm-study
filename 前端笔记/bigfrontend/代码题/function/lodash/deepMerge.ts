import _ from 'lodash'

type Patch<T> = Partial<{ [P in keyof T]: Patch<T[P]> }>

function deepMerge<T>(
  target: T,
  patch: Patch<T>,
  path?: (string | number)[],
  overWritePath?: string[]
): T {
  const result = target || ((Array.isArray(patch) ? [] : {}) as T)
  const entries = Object.entries(patch) as [keyof T, T[keyof T]][]
  for (const [key, value] of entries) {
    const originValue = result[key as keyof T]
    if (Object.prototype.toString.call(originValue) !== Object.prototype.toString.call(value)) {
      // 类型不等，使用覆盖更新
      result[key as keyof T] = value
    } else if (value === Object(value)) {
      // 对象 ｜ 数组
      if (
        path &&
        overWritePath?.length &&
        overWritePath?.indexOf([...path, key].join('.')) !== -1
      ) {
        // overWritePath 使用覆盖更新
        if (!isEqual(originValue, value)) {
          result[key as keyof T] = value
        }
      } else if (Array.isArray(value)) {
        // 长度相等的数组，使用深对比
        if (value.length === (originValue as Array<any>).length) {
          result[key as keyof T] = deepMerge(
            result[key as keyof T],
            value,
            [...(path || []), key as string],
            overWritePath
          )
        } else {
          result[key as keyof T] = value
        }
      } else {
        result[key as keyof T] = deepMerge(
          result[key as keyof T],
          value,
          [...(path || []), key as string],
          overWritePath
        )
      }
    } else {
      if (result[key as keyof T] !== value) {
        result[key as keyof T] = value
      }
    }
  }
  return result
}

console.log(deepMerge({ a: 2, foo: [0, { a: 1 }] }, { a: 1, b: { c: 2 }, foo: [{ a: 2, b: 2 }] })) // !数组长度不同，使用覆盖更新
console.log(_.merge({ a: 2, foo: [0, { a: 1 }] }, { a: 1, b: { c: 2 }, foo: [{ a: 2, b: 2 }] })) // 逐项merge

export {}

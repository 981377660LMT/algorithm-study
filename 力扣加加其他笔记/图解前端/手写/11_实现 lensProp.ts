import { NestedDict } from '../../../6_tree/重构json/10_获取nestdict的唯一识别'

function lensProp(lens: string, obj: NestedDict<string | number>) {
  const keys = lens.split('.')
  if (keys.length < 1) return undefined
  return keys.reduce<any>((pre, cur) => (pre !== void 0 ? pre[cur] : pre), obj)
}

const a = lensProp('a', { a: 1 }) // 1
const b = lensProp('b', { a: 1 }) // undefined
const c = lensProp('a.b', { a: { b: 'c' } }) // c
const d = lensProp('a.b.c.d.e.f', { a: { b: 'c' } }) // undefined

console.log(a)
console.log(b)
console.log(c)
console.log(d)

export {}

// 注意这个取不到索引 则为undefined
console.log((Number(1) as any)[1])
console.log((String(1) as any)[1])

function foo<T extends object>(obj: T): T {
  return obj
}

function foo2<T extends Record<PropertyKey, unknown>>(obj: T): T {
  return obj
}

interface Data {
  code: string
  name: string
}

type Data2 = {
  code: string
  name: string
}
const a: Data = { code: '1', name: '2' }
foo<Data>(a) // ok
foo2<Data>({ code: '1', name: '2' }) // 类型“Data”中缺少类型“string”的索引签名。ts(2344)

export {}

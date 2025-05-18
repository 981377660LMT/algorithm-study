// zustand为什么必须传泛型参数而不是自动推断
// Because state generic T is invariant.
// !As long as the generic to be inferred is invariant (i.e. both covariant and contravariant), TypeScript will be unable to infer it

declare const create: <T>(f: (get: () => T) => T) => T

const x = create(get => ({
  foo: 0,
  bar: () => get()
}))
// `x` is inferred as `unknown` instead of
// interface X {
//   foo: number,
//   bar: () => X
// }

export {}

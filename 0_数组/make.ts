// interface MakeArray<ArrayItem = number> {
//   (...lengths: number[]): ArrayItem[]
//   fill: (value: ArrayItem) => void
// }

// const make: MakeArray<number> = () => []
// make.fill = () => {}

// export { make }

type MakeArray<
  Length extends string,
  Result extends any[] = []
> = `${Result['length']}` extends Length ? Result : MakeArray<Length, [...Result, 1]>

type PopLength<T extends any[]> = T extends [...infer F, number] ? [...F]['length'] : 0

type MinusOne<N extends number> = PopLength<MakeArray<`${N}`>>

type NestedArray<ArrayItem, Depth extends number> = Depth extends 1
  ? ArrayItem[]
  : NestedArray<ArrayItem[], MinusOne<Depth>>

/**
 *
 * @param initValue 数组初始化值
 * @param sizes 数组各个维度的大小
 * @returns
 */
function make<T = number, S extends number[] = []>(
  initValue: T,
  ...sizes: S
): NestedArray<T, S['length']> {
  const dimension = sizes.length
  if (dimension === 1) return Array(sizes[0]).fill(initValue) as any
  return Array.from({ length: sizes.shift()! }, () => make(initValue, ...sizes)) as any
}

if (require.main === module) {
  console.log(make(2, 10))
  console.log(make(2, 2, 3))
  console.log(make(0, 2, 3, 4))
}

export { make }

import { ODT } from '../ODT-fastset'

const INF = 2e15

class RangeAssignRangeFreq {
  private readonly _odt: ODT<number>

  constructor(arr: ArrayLike<number>) {
    this._odt = new ODT(arr.length, INF)
    for (let i = 0; i < arr.length; ++i) {
      this._odt.set(i, i + 1, arr[i])
    }
  }

  assign(start: number, end: number, value: number): void {
    this._odt.set(start, end, value)
  }

  freq(start: number, end: number, value: number): number {
    let res = 0
    this._odt.enumerateRange(
      start,
      end,
      (s, e, v) => {
        if (value === v) {
          res += e - s
        }
      },
      false
    )
    return res
  }
}

export { RangeAssignRangeFreq }

if (require.main === module) {
  const arr = [1, 2, 3, 4, 5]
  const rar = new RangeAssignRangeFreq(arr)
  console.log(rar.freq(0, 5, 3))
  console.log(rar.freq(0, 5, 4))
}

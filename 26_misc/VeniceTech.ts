// VeniceTech
// !Venice技巧:使用全局变量维护偏移量
// https://www.cnblogs.com/KIMYON/p/14794243.html
// https://github.com/ShahjalalShohag/code-library/blob/76afa8d832e691c2bf274df341e69831351055f3/Data%20Structures/Venice%20Technique.cpp#L1
// 该技巧可以看成一种数据结构.
// 并支持以下操作:
// 1.添加一个新元素到集合中.
// 2.把一个元素从集合中删去.
// 3.把集合中的所有元素加上一个值.
// 4.得到所有元素的最小值
// 实际上其和普通set的区别仅在于操作3，对于该操作其实只需要使用一个全局delta，修改和查询时都加上/减去该值即可。

import { SortedListFast } from '../22_专题/离线查询/根号分治/SortedList/SortedListFast'

class VeniceTech {
  private readonly _sl: SortedListFast<number> = new SortedListFast()
  private _delta = 0

  add(value: number): void {
    this._sl.add(value + this._delta)
  }

  discard(value: number): void {
    this._sl.discard(value + this._delta)
  }

  addAll(delta: number): void {
    this._delta -= delta
  }

  min(): number | undefined {
    if (!this.size()) return undefined
    return this._sl.min! - this._delta
  }

  size(): number {
    return this._sl.length
  }
}

export {}

if (require.main === module) {
  const vt = new VeniceTech()
  vt.add(1)
  vt.add(2)
  vt.add(3)
  console.log(vt.min())
  vt.addAll(10)
  console.log(vt.min())
  vt.discard(11)
  console.log(vt.min())
}

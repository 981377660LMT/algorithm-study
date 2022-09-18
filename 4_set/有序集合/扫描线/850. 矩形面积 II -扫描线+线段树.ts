// // !找出平面中所有矩形叠加覆盖后的总面积。 由于答案可能太大，请返回它对 10 ^ 9 + 7 取模的结果。
// // 1.从左往右，线性扫描 注意排序
// // 2.扫描，看左边界还是有右边界
// // 3.计算底边宽*高，高度使用线段树维护

// const MOD = 1e9 + 7

// /**
//  * !线段树维护y轴被覆盖的长度length 以及被覆盖的次数cnt
//  */
// class SegmentTree {
//   private readonly _size: number
//   private readonly _count: Uint32Array // 被覆盖的次数
//   private readonly _length: number[] // 被覆盖的长度
//   private readonly _lazyValue: number[]
//   private readonly _nums: number[]

//   constructor(nOrNums: number | number[]) {
//     this._size = Array.isArray(nOrNums) ? nOrNums.length : nOrNums
//     this._count = new Uint32Array(this._size << 2)
//     this._length = Array(this._size << 2).fill(0)
//     this._lazyValue = Array(this._size << 2).fill(0)
//     this._nums = Array.isArray(nOrNums) ? nOrNums : Array(this._size).fill(0)
//   }

//   update(left: number, right: number, delta: number): void {
//     if (left < 1) left = 1
//     if (right > this._size) right = this._size
//     if (left > right) return
//     this._update(1, left, right, 1, this._size, delta)
//   }

//   queryLength(): number {
//     return this._length[1]
//   }

//   private _update(rt: number, L: number, R: number, l: number, r: number, delta: number): void {
//     if (L <= l && r <= R) {
//       this._count[rt] += delta
//       this._pushUp(rt, l, r)
//       return
//     }

//     const mid = Math.floor((l + r) / 2)
//     if (L <= mid) this._update(rt << 1, L, R, l, mid, delta)
//     if (mid < R) this._update((rt << 1) | 1, L, R, mid + 1, r, delta)
//     this._pushUp(rt, l, r)
//   }

//   private _pushUp(rt: number, l: number, r: number): void {
//     if (this._count[rt] > 0) {
//       this._length[rt] = this._nums[r] - this._nums[l - 1]
//     } else if (l !== r) {
//       this._length[rt] = this._length[rt << 1] + this._length[(rt << 1) | 1]
//     } else {
//       this._length[rt] = 0
//     }
//   }
// }

// function rectangleArea(rectangles: number[][]): number {
//   const events: [x: number, type: 0 | 1, y1: number, y2: number][] = []
//   const ySet = new Set<number>()
//   rectangles.forEach(([x1, y1, x2, y2]) => {
//     events.push([x1, 0, y1, y2]) // !左边界进入
//     events.push([x2, 1, y1, y2]) // !右边界离开
//     ySet.add(y1).add(y2) // !记录所有y坐标
//   })

//   events.sort(([x1], [x2]) => x1 - x2) // !按x排序
//   const yNums = [...ySet].sort((a, b) => a - b) // !离散化y坐标
//   const mp = new Map<number, number>()
//   yNums.forEach((num, i) => mp.set(num, i + 1)) // !y坐标映射到1~n

//   const tree = new SegmentTree(yNums)
//   let res = 0
//   events.forEach(([x, type, y1, y2], i) => {
//     if (i) {
//       const yLen = tree.queryLength() // !y轴上被覆盖的长度
//       res += (x - events[i - 1][0]) * yLen
//       res %= MOD
//     }

//     if (type === 0) {
//       tree.update(mp.get(y1)!, mp.get(y2)! - 1, 1) // !左边界进入
//     } else {
//       tree.update(mp.get(y1)!, mp.get(y2)! - 1, -1) // !右边界离开
//     }
//   })

//   return res
// }

// if (require.main === module) {
//   const rectangles = [[0, 0, 1000000000, 1000000000]]
//   console.log(rectangleArea(rectangles))
// }

// export {}

// !超出2**53-1了
// TODO
// https://www.acwing.com/solution/content/1027/
// 亚特兰蒂斯

// 对每个下标，查询 最右侧/最左侧 lower/floor/ceiling/higher 的元素

// import { useBlock } from './SqrtDecomposition/useBlock'

// class RightMostLeftMostQuery {
//   private readonly _nums: number[]
//   private readonly _belong: Uint16Array
//   private readonly _blockStart: Uint32Array
//   private readonly _blockEnd: Uint32Array
//   private readonly _blockCount: number

//   private readonly _blockMin: number[]
//   private readonly _blockMax: number[]
//   private readonly _blockLazy: number[]

//   constructor(nums: number[]) {
//     this._nums = nums.slice()
//     const { blockCount, belong, blockStart, blockEnd } = useBlock(nums)
//     this._belong = belong
//     this._blockStart = blockStart
//     this._blockEnd = blockEnd
//     this._blockCount = blockCount
//     this._blockMin = Array(blockCount).fill(Infinity)
//     this._blockMax = Array(blockCount).fill(-Infinity)
//     nums.forEach((num, i) => {
//       const bid = belong[i]
//       this._blockMin[bid] = Math.min(this._blockMin[bid], num)
//       this._blockMax[bid] = Math.max(this._blockMax[bid], num)
//     })
//   }

//   set(index: number, value: number): void {}

//   addRange(start: number, end: number, delta: number): void {}

//   /**
//    * 查询`index`右侧最远的下标`j`，使得 `nums[j] < target`.
//    * 如果不存在，返回`-1`.
//    */
//   rightMostLower(index: number): number {}

//   /**
//    * 查询`index`右侧最远的下标`j`，使得 `nums[j] <= target`.
//    * 如果不存在，返回`-1`.
//    */
//   rightMostFloor(index: number): number {}

//   /**
//    * 查询`index`右侧最远的下标`j`，使得 `nums[j] >= target`.
//    * 如果不存在，返回`-1`.
//    */
//   rightMostCeiling(index: number): number {}

//   /**
//    * 查询`index`右侧最远的下标`j`，使得 `nums[j] > target`.
//    * 如果不存在，返回`-1`.
//    */
//   rightMostHigher(index: number): number {}

//   /**
//    * 查询`index`左侧最远的下标`j`，使得 `nums[j] < target`.
//    * 如果不存在，返回`-1`.
//    */
//   leftMostLower(index: number): number {}

//   /**
//    * 查询`index`左侧最远的下标`j`，使得 `nums[j] <= target`.
//    * 如果不存在，返回`-1`.
//    */
//   leftMostFloor(index: number): number {}

//   /**
//    * 查询`index`左侧最远的下标`j`，使得 `nums[j] >= target`.
//    * 如果不存在，返回`-1`.
//    */
//   leftMostCeiling(index: number): number {}

//   /**
//    * 查询`index`左侧最远的下标`j`，使得 `nums[j] > target`.
//    * 如果不存在，返回`-1`.
//    */
//   leftMostHigher(index: number): number {}
// }

// export { RightMostLeftMostQuery }

// if (require.main === module) {
//   // 962. 最大宽度坡
//   // https://leetcode.cn/problems/maximum-width-ramp/

//   function maxWidthRamp(nums: number[]): number {}
// }

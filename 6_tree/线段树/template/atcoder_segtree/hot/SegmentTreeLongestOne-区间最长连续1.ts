/* eslint-disable no-inner-declarations */

// SegmentTreeLongestOne-区间最长连续1

import { createPointSetRangeLongestOne } from '../SegmentTreeUtils'

export { createPointSetRangeLongestOne as SegmentTreeLongestOne }

if (require.main === module) {
  const { fromElement, tree } = createPointSetRangeLongestOne([0, 1, 0, 1, 1, 1])
  console.log(tree.queryAll().pairCount)
  tree.set(2, fromElement(1))
  console.log(tree.queryAll().pairCount)
  tree.set(3, fromElement(0))
  console.log(tree.queryAll().pairCount)

  testTime()
  function testTime(): void {
    const n = 1e5
    const arr = Array(n)
      .fill(0)
      .map<0 | 1>((_, i) => (i & 1) as 0 | 1)

    console.time('SegmentTreeLongestOne')
    const { fromElement, tree } = createPointSetRangeLongestOne(arr)
    for (let i = 0; i < n; ++i) {
      tree.set(i, fromElement((i & 1) as 0 | 1))
      tree.query(i, n)
    }
    console.timeEnd('SegmentTreeLongestOne') // SegmentTreeLongestOne: 357.718ms
  }
}

/* eslint-disable no-inner-declarations */

import { createPointSetRangeLongestRepeating } from '../SegmentTreeUtils'

export { createPointSetRangeLongestRepeating as SegmentTreeLongestRepeating }

if (require.main === module) {
  function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
    const { tree, fromElement } = createPointSetRangeLongestRepeating(s)
    const res: number[] = Array(queryIndices.length)
    for (let i = 0; i < queryIndices.length; i++) {
      const pos = queryIndices[i]
      const char = queryCharacters[i]
      tree.set(pos, fromElement(char))
      res[i] = tree.queryAll().max
    }
    return res
  }
}

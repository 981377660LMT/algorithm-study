import { SegmentTreePointUpdateRangeQuery } from '../template/atcoder_segtree/SegmentTreePointUpdateRangeQuery'

type E = {
  size: number
  preMax: number
  sufMax: number
  max: number
  leftChar: string
  rightChar: string
}

function longestRepeating(s: string, queryCharacters: string, queryIndices: number[]): number[] {
  const n = s.length
  const leaves: E[] = Array(n)
  for (let i = 0; i < n; i++) {
    leaves[i] = { size: 1, preMax: 1, sufMax: 1, max: 1, leftChar: s[i], rightChar: s[i] }
  }

  const seg = new SegmentTreePointUpdateRangeQuery<E>(
    leaves,
    () => ({
      size: 0,
      preMax: 0,
      sufMax: 0,
      max: 0,
      leftChar: '',
      rightChar: ''
    }),
    (a, b) => {
      const res = {
        size: a.size + b.size,
        leftChar: a.leftChar,
        rightChar: b.rightChar,
        preMax: 0,
        sufMax: 0,
        max: 0
      }
      if (a.rightChar === b.leftChar) {
        res.preMax = a.preMax
        if (a.preMax === a.size) res.preMax += b.preMax
        res.sufMax = b.sufMax
        if (b.sufMax === b.size) res.sufMax += a.sufMax
        res.max = Math.max(a.max, b.max, a.sufMax + b.preMax)
      } else {
        res.preMax = a.preMax
        res.sufMax = b.sufMax
        res.max = Math.max(a.max, b.max)
      }
      return res
    }
  )

  const res: number[] = Array(queryIndices.length)
  for (let i = 0; i < queryIndices.length; i++) {
    const pos = queryIndices[i]
    const char = queryCharacters[i]
    seg.set(pos, { size: 1, preMax: 1, sufMax: 1, max: 1, leftChar: char, rightChar: char })
    res[i] = seg.queryAll().max
  }
  return res
}

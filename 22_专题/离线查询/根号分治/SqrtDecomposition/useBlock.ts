function useBlock(arr: ArrayLike<unknown>): {
  /** 下标所属的块. */
  belong: Uint16Array
  /** 每个块的起始下标(包含). */
  blockStart: Uint32Array
  /** 每个块的结束下标(不包含). */
  blockEnd: Uint32Array
  /** 块的数量. */
  blockCount: number
} {
  const n = arr.length
  const blockSize = (Math.sqrt(n) + 1) | 0
  const blockCount = 1 + ((n / blockSize) | 0)
  const blockStart = new Uint32Array(blockCount)
  const blockEnd = new Uint32Array(blockCount)
  const belong = new Uint16Array(n)
  for (let i = 0; i < blockCount; i++) {
    blockStart[i] = i * blockSize
    blockEnd[i] = Math.min((i + 1) * blockSize, n)
  }
  for (let i = 0; i < n; i++) {
    belong[i] = (i / blockSize) | 0
  }

  return {
    belong,
    blockStart,
    blockEnd,
    blockCount
  }
}

export { useBlock }

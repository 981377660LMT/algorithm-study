class ArrayReader {
  // This is the ArrayReader's API interface.
  // You should not implement it, or speculate about its implementation
  // @ts-ignore
  get(index: number): number {}
}

// 给定一个升序整数数组，写一个函数搜索 nums 中数字 target
// 这个数组的大小是未知的。你只可以通过 ArrayReader 接口访问这个数组
function search(reader: ArrayReader, target: number): number {
  if (reader.get(0) == target) return 0
  let l = 0
  let r = 1

  // 初始化边界
  while (reader.get(r) < target) {
    l = r
    r <<= 1
  }

  // 因此当 left <= right 的时候，解空间都不为空，此时我们都需要继续搜索
  while (l <= r) {
    const mid = (l + r) >> 1
    const midElement = reader.get(mid)
    if (midElement === target) return mid
    else if (midElement < target) l = mid + 1
    else if (midElement > target) r = mid - 1
  }

  return -1
}

export {}

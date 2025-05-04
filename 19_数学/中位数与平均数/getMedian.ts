/** 有序数组中位数(向下取整). */
function getMedian(arr: ArrayLike<number>, start = 0, end = arr.length): number {
  if (start < 0) start = 0
  if (end > arr.length) end = arr.length
  if (start >= end) return 0
  if ((end - start) & 1) {
    return arr[(end + start) >>> 1]
  }
  return Math.floor((arr[(end + start) >>> 1] + arr[((end + start) >>> 1) - 1]) / 2)
}

export { getMedian }

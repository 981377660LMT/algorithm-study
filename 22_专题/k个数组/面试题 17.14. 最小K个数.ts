function smallestK(arr: number[], k: number): number[] {
  const partition = (nums: number[], l: number, r: number): number => {
    /**
     * @description 生成[start,end]的随机整数
     */
    const randint = (start: number, end: number) => {
      if (start > end) throw new Error('invalid interval')
      const amplitude = end - start
      return Math.floor((amplitude + 1) * Math.random()) + start
    }

    // 优化，随机取标定点，以解决近乎有序的列表
    const randIndex = randint(l, r)
    ;[[nums[l], nums[randIndex]]] = [[nums[randIndex], nums[l]]]

    let pivotIndex = l
    const pivot = nums[l]
    for (let i = l + 1; i <= r; i++) {
      if (nums[i] < pivot) {
        // 这里要先移pivotIndex是因为不能动最左边的比较元素 比较元素要最后移到自己的位置
        pivotIndex++
        ;[[nums[i], nums[pivotIndex]]] = [[nums[pivotIndex], nums[i]]]
      }
    }

    // pivot放到中间应有的位置
    ;[[nums[l], nums[pivotIndex]]] = [[nums[pivotIndex], nums[l]]]

    return pivotIndex
  }

  // 最后要让pivotIndex为k 前面的k个数就是满足题意的
  const divide = (nums: number[], l: number, r: number, k: number): void => {
    if (l < r) {
      const pivotIndex = partition(nums, l, r)
      if (pivotIndex > k) divide(nums, l, pivotIndex - 1, k)
      else if (pivotIndex < k) divide(nums, pivotIndex + 1, r, k)
    }
  }

  // 找第(前)k小的数 k>=1
  const findKSmallest = (nums: number[], k: number) => {
    divide(nums, 0, nums.length - 1, k)
    return nums.slice(0, k)
  }

  return findKSmallest(arr, k)
}

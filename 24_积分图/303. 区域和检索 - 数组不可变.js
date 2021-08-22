// 给定一个整数数组  nums，求出数组从索引 i 到 j（i ≤ j）范围内元素的总和，包含 i、j 两点。
function NumArray(nums) {
  this.sums = []
  var sum = 0
  for (var i = 0; i < nums.length; i++) {
    sum += nums[i]
    this.sums.push(sum)
  }
}

// 注意是i-1
NumArray.prototype.sumRange = function (i, j) {
  return this.sums[j] - (i > 0 ? this.sums[i - 1] : 0)
}

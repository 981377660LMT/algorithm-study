**判断是否存在**:while l <= r
当元素不重复时，如果 nums[i] <= nums[j]，说明区间 [i,j] 是「连续递增」的。
判断区间是否「连续递增」，只需比较区间边界值：如果 nums[left] <= nums[mid]，则区间 [left,mid] 连续递增；反之，区间 [mid,right] 连续递增。但是上述判断仅适用于数组中不含重复元素的情况，如果数组中包含重复元素，那么在 nums[left]==nums[mid] 时将退化为线性查找
有重复元素(nums[left]===nums[mid])时 left++

**查找最小值**: while l < r
需要与右边元素比较看最小值在哪节 **不断逼近左边**
有重复元素(nums[mid]===nums[right])时向左逼近 right--

链接：https://leetcode-cn.com/problems/find-minimum-in-rotated-sorted-array/solution/yi-wen-jie-jue-4-dao-sou-suo-xuan-zhuan-pai-xu-s-3/

小技巧:

一般是这样,

当 while left < right 是循环外输出

当 while left <= right 是循环里输出

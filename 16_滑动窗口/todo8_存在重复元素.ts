// 请你判断是否存在 两个不同下标 i 和 j，使得 abs(nums[i] - nums[j]) <= t ，
// 同时又满足 abs(i - j) <= k 。

// 暴力:O(n^2)
// 思路：在大小为k的滑动窗口中，每多一个元素就看查找表中是否有nums[r]-t到nums[r]+t的元素
// 希望有一个顺序性的查找表
// 参考java里的treeset时间O(nlongk) - treeset的底层是红黑树实现的
// 或者是二分搜索树 todo
const containsNearbyAlmostDuplicate = (nums: number[], k: number, t: number): boolean => {}

console.log(containsNearbyAlmostDuplicate([1, 5, 9, 1, 5, 9], 2, 3))
console.log(containsNearbyAlmostDuplicate([1, 2, 3, 1], 3, 0))

export {}

# 1095. 山脉数组中查找目标值


# https://leetcode.cn/problems/find-in-mountain-array/
# 给你一个 山脉数组 mountainArr，
# !请你返回能够使得 mountainArr.get(index) 等于 `target 最小 的下标 index 值`。
# 如果不存在这样的下标 index，就请返回 -1。
# 你将 不能直接访问该山脉数组，必须通过 MountainArray 接口来获取数据：
# 对 MountainArray.get 发起超过 100 次调用的提交将被视为错误答案。
# 3 <= mountain_arr.length() <= 10000


# !先查找出峰所在的下标 peek，然后分别在 [0, peek) 和 [peek, mountainArr.length()) 两个区间内查找答案。


class MountainArray:
    def get(self, index: int) -> int:
        ...

    def length(self) -> int:
        ...


class Solution:
    def findInMountainArray(self, target: int, mountain_arr: "MountainArray") -> int:
        n = mountain_arr.length()
        left = 0
        right = n - 1
        while left <= right:
            mid = (left + right) >> 1
            if mountain_arr.get(mid) < mountain_arr.get(mid + 1):
                left = mid + 1
            else:
                right = mid - 1

        peak = left

        # find target in the left of peak
        left, right = 0, peak
        while left <= right:
            mid = (left + right) >> 1
            if mountain_arr.get(mid) < target:
                left = mid + 1
            elif mountain_arr.get(mid) > target:
                right = mid - 1
            else:
                return mid

        # find target in the right of peak
        left, right = peak, n - 1
        while left <= right:
            mid = (left + right) >> 1
            if mountain_arr.get(mid) > target:
                left = mid + 1
            elif mountain_arr.get(mid) < target:
                right = mid - 1
            else:
                return mid

        return -1


# 输入：array = [1,2,3,4,5,3,1], target = 3
# 输出：2
# 解释：3 在数组中出现了两次，下标分别为 2 和 5，我们返回最小的下标 2。

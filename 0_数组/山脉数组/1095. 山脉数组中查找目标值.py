class MountainArray:
    def get(self, index: int) -> int:
        ...

    def length(self) -> int:
        ...


# 给你一个 山脉数组 mountainArr，请你返回能够使得 mountainArr.get(index) 等于 `target 最小 的下标 index 值`。
# 如果不存在这样的下标 index，就请返回 -1。
# 你将 不能直接访问该山脉数组，必须通过 MountainArray 接口来获取数据：
# 对 MountainArray.get 发起超过 100 次调用的提交将被视为错误答案。
# 3 <= mountain_arr.length() <= 10000
class Solution:
    def findInMountainArray(self, target: int, mountain_arr: 'MountainArray') -> int:
        n = mountain_arr.length()
        l = 0
        r = n - 1
        while l <= r:
            mid = (l + r) >> 1
            if mountain_arr.get(mid) < mountain_arr.get(mid + 1):
                l = mid + 1
            else:
                r = mid - 1

        peak = l

        # find target in the left of peak
        l, r = 0, peak
        while l <= r:
            mid = (l + r) >> 1
            if mountain_arr.get(mid) < target:
                l = mid + 1
            elif mountain_arr.get(mid) > target:
                r = mid - 1
            else:
                return mid
        # find target in the right of peak
        l, r = peak, n - 1
        while l <= r:
            mid = (l + r) >> 1
            if mountain_arr.get(mid) > target:
                l = mid + 1
            elif mountain_arr.get(mid) < target:
                r = mid - 1
            else:
                return mid

        return -1


# 输入：array = [1,2,3,4,5,3,1], target = 3
# 输出：2
# 解释：3 在数组中出现了两次，下标分别为 2 和 5，我们返回最小的下标 2。


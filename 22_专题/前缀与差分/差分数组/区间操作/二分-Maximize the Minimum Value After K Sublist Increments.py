# k次操作 每次可以范围为区间加1
# 求k次操作后最大的最小值

# 1526. 形成目标数组的子数组最少增加次数 copy
class Solution:
    def solve(self, nums, size, k):
        def check(target):
            diff = [0] * n
            res = curAdd = 0
            for i in range(n):
                curAdd += diff[i]
                delta = target - (nums[i] + curAdd)
                if delta > 0:
                    res += delta
                    curAdd += delta
                    if i + size < n:
                        diff[i + size] -= delta

            return res <= k

        n = len(nums)

        left, right = 0, int(1e20)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        return right

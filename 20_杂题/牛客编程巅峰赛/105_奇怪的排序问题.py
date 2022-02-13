# 需要把队伍变成升序，也就是从矮到高排序
# 现在给出数 n 和一个 1 到 n 的排列，求最少的选择次数，使数组变为升序。
class Solution:
    def wwork(self, n: int, nums: list[int]) -> int:
        """倒序遍历维护最小值即可"""
        min_ = nums[-1]
        res = 0
        for i in range(n - 2, -1, -1):
            cur = nums[i]
            if cur > min_:
                res += 1
            else:
                min_ = cur
        return res


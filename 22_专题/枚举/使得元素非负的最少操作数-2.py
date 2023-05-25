# https://leetcode.cn/problems/minimum-operations-to-make-numbers-non-positive/
# 2702. Minimum Operations to Make Numbers Non-positive
# 所有元素变成非正数的最少操作数
# 一个数组开始全为正数,每回合选择一个下标i。
# !将nums[i]减去x,`然后`将除开下标i以外的所有元素减去y。
# 求使得数组所有元素都<=0的最少操作数。


# 1 <= nums.length <= 1e5
# 1 <= nums[i] <= 1e9
# 1 <= y < x <= 1e9


# 注意到条件y<x.
# !等价于:执行k次所有元素减去y的操作,然后执行k次某个元素减去(y-x)的操作.
# !二分答案.


from typing import List


class Solution:
    def minOperations(self, nums: List[int], x: int, y: int) -> int:
        def check(mid: int) -> bool:
            """操作mid轮是否能使所有元素非正数."""
            arr = [num - mid * y for num in nums]
            todo = 0
            for num in arr:
                if num > 0:
                    todo += (num + diff - 1) // diff
            return todo <= mid

        diff = x - y
        left, right = 0, int(1e10)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


if __name__ == "__main__":
    print(Solution().minOperations([3, 4, 1, 7, 6], 4, 2))

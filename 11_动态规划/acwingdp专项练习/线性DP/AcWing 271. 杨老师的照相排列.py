# AcWing 271. 杨老师的照相排列
# 1 2 3
# 4 5
# 6
# 在合影时要求每一排从左到右身高递减，每一列从后到前身高也递减。
# 问一共有多少种安排合影位置的方案？
# !k注意身高最高的是 1，最低的是 6
# 1≤k≤5,学生总人数不超过  30  人。
# 状态表示：用(a,b,c,d,e)来表示当前的状态(轮廓)
# 状态转移：最矮的一个人被安排在那个位置 (a,b,c,d,e) -> (a,b,c,d,e-1)
# 时间复杂度(30 / 5 + 1) ^ 5 * 5

# dp优化爆搜 方案划分集合

from functools import lru_cache
from typing import List, Tuple


def main(nums: List[int]) -> int:
    # @lru_cache(None)
    # def dfs(a: int, b: int, c: int, d: int, e: int) -> int:
    #     if a < 0 or b < 0 or c < 0 or d < 0 or e < 0:
    #         return 0
    #     if a == b == c == d == e == 0:
    #         return 1
    #     res = 0
    #     """注意每排人数大于等于前一排"""
    #     if a - 1 >= b:
    #         res += dfs(a - 1, b, c, d, e)
    #     if b - 1 >= c:
    #         res += dfs(a, b - 1, c, d, e)
    #     if c - 1 >= d:
    #         res += dfs(a, b, c - 1, d, e)
    #     if d - 1 >= e:
    #         res += dfs(a, b, c, d - 1, e)
    #     if e - 1 >= 0:
    #         res += dfs(a, b, c, d, e - 1)
    #     return res

    # nums = nums + [0] * (5 - len(nums))
    # return dfs(nums[0], nums[1], nums[2], nums[3], nums[4])
    @lru_cache(None)
    def dfs(nums: Tuple[int, ...]) -> int:
        if any(num < 0 for num in nums):
            return 0
        if all(num == 0 for num in nums):
            return 1
        res = 0
        """注意每排人数大于等于前一排"""
        for i, (a, b) in enumerate(zip(nums, nums[1:])):
            if a - 1 >= b:
                res += dfs(tuple(nums[:i] + (a - 1,) + nums[i + 1 :]))
        return res

    nums = nums + [0] * (5 - len(nums))
    nums.append(0)  # 哨兵
    return dfs(tuple(nums))


while True:
    k = int(input())  # 总排数
    if k == 0:
        break
    # 从后向前每排的具体人数
    nums = list(map(int, input().split()))
    print(main(nums))

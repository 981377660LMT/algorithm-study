from typing import List

# 返回存在长度 恰好 为 m 的 一组 1  的`最后`步骤

# 思路: 维护groups的长度
# length[i]=val表示 index[i]位于长为val的组里
# This solution keeps track of the right/left sides of each of the groups.
# The left side references the right side, the right side references the left side.

#  思路:update length[a - left], length[a + right] to left + right + 1.


class Solution:
    def findLatestStep(self, arr: List[int], m: int) -> int:
        if m > len(arr):
            return -1
        if m == len(arr):
            return m

        # +2表示前后哨兵
        length = [0] * (len(arr) + 2)
        res = -1
        for i, a in enumerate(arr):
            left, right = length[a - 1], length[a + 1]
            if left == m or right == m:
                res = i
            length[a - left] = length[a + right] = left + right + 1
            print(length)
        return res


print(Solution().findLatestStep(arr=[3, 5, 1, 2, 4], m=1))
# 输出：4
# 解释：
# 步骤 1："00100"，由 1 构成的组：["1"]
# 步骤 2："00101"，由 1 构成的组：["1", "1"]
# 步骤 3："10101"，由 1 构成的组：["1", "1", "1"]
# 步骤 4："11101"，由 1 构成的组：["111", "1"]
# 步骤 5："11111"，由 1 构成的组：["11111"]
# 存在长度为 1 的一组 1 的最后步骤是步骤 4 。

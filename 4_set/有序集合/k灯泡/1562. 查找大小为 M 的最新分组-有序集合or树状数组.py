from typing import List
from sortedcontainers import SortedList

# 返回存在长度 恰好 为 m 的 一组 1  的`最后`步骤

# 反向查找第一个 出现m个1的情况


class Solution:
    def findLatestStep(self, arr: List[int], m: int) -> int:
        n = len(arr)

        if m > n:
            return -1
        if m == n:
            return m

        ones = SortedList([-1, n])

        for i in range(n - 1, -1, -1):
            pos = arr[i] - 1
            ones.add(pos)
            index = ones.bisect_left(pos)

            if index + 1 < len(ones):
                rightPos = ones[index + 1]
                if rightPos - pos == m + 1:
                    return i

            if index > 0:
                leftPos = ones[index - 1]
                if pos - leftPos == m + 1:
                    return i

        return -1


print(Solution().findLatestStep(arr=[3, 5, 1, 2, 4], m=1))
print(Solution().findLatestStep(arr=[1, 2], m=1))
# 输出：4
# 解释：
# 步骤 1："00100"，由 1 构成的组：["1"]
# 步骤 2："00101"，由 1 构成的组：["1", "1"]
# 步骤 3："10101"，由 1 构成的组：["1", "1", "1"]
# 步骤 4："11101"，由 1 构成的组：["111", "1"]
# 步骤 5："11111"，由 1 构成的组：["11111"]
# 存在长度为 1 的一组 1 的最后步骤是步骤 4 。

from typing import List, Optional, Tuple

MOD = int(1e9 + 7)

# aliceArrows.length == bobArrows.length == 12
# 1 <= numArrows <= 105
class Solution:
    def maximumBobPoints(self, numArrows: int, aliceArrows: List[int]) -> List[int]:

        # 如果 ak < bk ，则 Bob 得 k 分
        resCand = [0] * 12
        score = 0
        for state in range(1 << 12):
            cur = [0] * 12
            curScore = 0
            for i in range(12):
                if (state >> i) & 1:
                    cur[i] = 1
                    curScore += i
            if curScore < score:
                continue

            minCost = 0
            select = [0] * 12
            for i, num in enumerate(cur):
                if num > 0:
                    minCost += aliceArrows[i] + 1
                    select[i] = aliceArrows[i] + 1
            if minCost > numArrows:
                continue

            if curScore > score:
                if minCost < numArrows:
                    select[0] += numArrows - minCost
                score = curScore
                resCand = select
        return resCand


print(Solution().maximumBobPoints(9, [1, 1, 0, 1, 0, 0, 2, 1, 0, 1, 2, 0]))
print(Solution().maximumBobPoints(89, [3, 2, 28, 1, 7, 1, 16, 7, 3, 13, 3, 5]))
# 预期：
# [21,3,0,2,8,2,17,8,4,14,4,6]

# [0,3,0,2,8,2,17,8,4,14,4,6]
# [21,3,0,2,8,2,17,8,4,14,4,6]


from typing import List, Optional, Tuple

MOD = int(1e9 + 7)

# aliceArrows.length == bobArrows.length == 12
# 1 <= numArrows <= 105

# 二进制枚举与01背包的解法
class Solution:
    def maximumBobPoints(self, numArrows: int, aliceArrows: List[int]) -> List[int]:

        # 如果 ak < bk ，则 Bob 得 k 分
        resCand = [0] * 12
        resScore = 0
        for state in range(1 << 12):
            score, arrow, bobArrows = 0, 0, [0] * 12
            for i in range(12):
                if (state >> i) & 1:
                    score += i
                    arrow += aliceArrows[i] + 1
                    bobArrows[i] = aliceArrows[i] + 1
            if arrow > numArrows:
                continue
            if score > resScore:
                resScore = score
                bobArrows[0] += numArrows - arrow
                resCand = bobArrows

        return resCand


print(Solution().maximumBobPoints(9, [1, 1, 0, 1, 0, 0, 2, 1, 0, 1, 2, 0]))
print(Solution().maximumBobPoints(89, [3, 2, 28, 1, 7, 1, 16, 7, 3, 13, 3, 5]))
# 预期：
# [21,3,0,2,8,2,17,8,4,14,4,6]

# [0,3,0,2,8,2,17,8,4,14,4,6]
# [21,3,0,2,8,2,17,8,4,14,4,6]


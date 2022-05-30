from typing import List
from functools import lru_cache

# 即不同的数字表示不同的颜色。
# 每一轮你可以移除具有相同颜色的连续 k 个盒子（k >= 1），这样一轮之后你将得到 k * k 个积分。
# 当你将所有盒子都去掉之后，求你能获得的最大积分和。

# 1 <= boxes.length <= 100

# 664. 奇怪的打印机.py
# https://leetcode.com/problems/remove-boxes/discuss/1402561/C%2B%2BJavaPython-Top-down-DP-Clear-explanation-with-Picture-Clean-and-Concise

# 时间复杂度 O(n^4)


class Solution:
    def removeBoxes(self, boxes: List[int]) -> int:
        # dp(l, r, k) denote the maximum points we can get in boxes[l..r] if we have extra k boxes which is the same color with boxes[l] in the left side.
        @lru_cache(None)
        def dfs(left: int, right: int, k: int) -> int:
            if left > right:
                return 0
            while left + 1 <= right and boxes[left] == boxes[left + 1]:
                left += 1
                k += 1

            res = (k + 1) ** 2 + dfs(left + 1, right, 0)  # 最原始的方案
            for mid in range(left + 1, right + 1):
                if boxes[left] == boxes[mid]:
                    res = max(res, dfs(left + 1, mid - 1, 0) + dfs(mid, right, k + 1))

            return res

        return dfs(0, len(boxes) - 1, 0)


print(Solution().removeBoxes(boxes=[1, 3, 2, 2, 2, 3, 4, 3, 1]))
# 输出：23
# 解释：
# [1, 3, 2, 2, 2, 3, 4, 3, 1]
# ----> [1, 3, 3, 4, 3, 1] (3*3=9 分)
# ----> [1, 3, 3, 3, 1] (1*1=1 分)
# ----> [1, 1] (3*3=9 分)
# ----> [] (2*2=4 分)


# 直接放弃

from typing import List

MOD = int(1e9 + 7)

# 2 <= h, w <= 10^9
# 1 <= horizontalCuts.length < min(h, 10^5)
class Solution:
    def maxArea(self, h: int, w: int, horizontalCuts: List[int], verticalCuts: List[int]) -> int:
        horizontalCuts.sort()
        verticalCuts.sort()

        hc = [0] + horizontalCuts + [h]
        wc = [0] + verticalCuts + [w]

        hm = max([hc[i] - hc[i - 1] for i in range(1, len(hc))])
        wm = max([wc[i] - wc[i - 1] for i in range(1, len(wc))])

        return wm * hm % MOD


print(Solution().maxArea(h=5, w=4, horizontalCuts=[1, 2, 4], verticalCuts=[1, 3]))

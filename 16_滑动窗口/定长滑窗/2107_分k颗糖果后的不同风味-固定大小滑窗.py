# 你有不同口味的 n 颗糖果，candies[i] 表示某颗糖的口味。
# 妈妈要求你分享连续 k 颗糖给你的妹妹，
# 请问在分享糖果后，最多能保留多少种不同口味的糖果？(统计窗口外的数量)
# 1 <= candies.length <= 105
# 0 <= k <= candies.length

# 如果不要求糖果连续，那么是贪心题（即将最富余的糖果先给出去，
# 最少保留一颗），
# 要求连续则变成了典型的哈希表 + 滑动窗口问题。
from typing import List
from collections import Counter


class Solution:
    def shareCandies(self, candies: List[int], k: int) -> int:
        counter = Counter(candies)
        if k == 0:
            return len(counter)

        res = 0
        for right, cur in enumerate(candies):
            counter[cur] -= 1
            if not counter[cur]:
                del counter[cur]
            if right >= k:
                counter[candies[right - k]] += 1
            if right >= k - 1:
                res = max(res, len(counter))
        return res


print(Solution().shareCandies(candies=[2, 2, 2, 2, 3, 3], k=2))
# 2
print(Solution().shareCandies(candies=[1, 2, 2, 3, 4, 3], k=3))

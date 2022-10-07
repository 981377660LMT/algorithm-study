from typing import List
from collections import defaultdict


MOD = int(1e9 + 7)
INF = int(1e20)

# 你可以选择一段区间的鲜花做成插花，且不能丢弃。
# 在你选择的插花中，如果每一品种的鲜花数量都不超过 cnt 朵，
# 那么我们认为这束插花是 「美观的」。
# !请返回在这一排鲜花中，共有多少种可选择的区间，使得插花是「美观的」。


class Solution:
    def beautifulBouquet(self, flowers: List[int], cnt: int) -> int:
        res, left, n = 0, 0, len(flowers)
        counter = defaultdict(int)
        for right in range(n):
            counter[flowers[right]] += 1
            while left <= right and counter[flowers[right]] > cnt:
                counter[flowers[left]] -= 1
                left += 1
            res += right - left + 1
            res %= MOD
        return res

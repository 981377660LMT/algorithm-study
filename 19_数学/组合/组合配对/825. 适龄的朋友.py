from collections import Counter
from typing import List

# 求总共会发出多少份好友请求?
class Solution:
    def numFriendRequests(self, ages: List[int]) -> int:
        def canSend(a, b):
            return not (b <= 0.5 * a + 7 or b > a or b > 100 and a < 100)

        counter = Counter(ages)
        res = 0
        for a in counter:
            for b in counter:
                if canSend(a, b):
                    if a == b:
                        res += counter[a] * (counter[a] - 1)
                    else:
                        res += counter[a] * counter[b]
        return res

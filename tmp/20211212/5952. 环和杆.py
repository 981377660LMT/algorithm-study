from collections import Counter


class Solution:
    def countPoints(self, rings: str) -> int:
        n = len(rings) // 2
        store = ['R', 'G', 'B']
        counter = [0] * 10
        for i in range(n):
            color, index = rings[2 * i], int(rings[(2 * i) + 1])
            counter[index] |= 1 << (store.index(color))
        return sum(c == 7 for c in counter)


print(Solution().countPoints(rings="B0B6G0R6R0R6G9"))


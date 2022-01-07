class Solution:
    def countPoints(self, rings: str) -> int:
        colors = ['R', 'G', 'B']
        states = [0] * 10
        for c, r in zip(rings[::2], rings[1::2]):
            states[int(r)] |= 1 << (colors.index(c))
        return sum(c == 7 for c in states)


print(Solution().countPoints(rings="B0B6G0R6R0R6G9"))


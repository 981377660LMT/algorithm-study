# Return the maximum value for all 0 ≤ i < j < n:
# |a[i] - a[j]| + |b[i] - b[j]| + |i - j|
class Solution:
    def solve(self, a, b):
        res = 0
        n = len(a)
        for s, t in [(-1, -1), (-1, 1), (1, -1), (1, 1)]:
            cur_min = float("inf")
            cur_max = float("-inf")
            for i in range(n):
                tmp = s * a[i] + t * b[i] + i
                cur_min = min(cur_min, tmp)
                cur_max = max(cur_max, tmp)
            res = max(res, cur_max - cur_min)
        return res

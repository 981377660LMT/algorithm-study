# n^2logn
from bisect import bisect_right

# Return the number of different unique indices i, j, k, l such that A[i] + B[j] + C[k] + D[l] ≤ target.
# 求四元组的数目


class Solution:
    def solve(self, A, B, C, D, target):
        ab = sorted(a + b for a in A for b in B)
        # a + b <= target - c - d
        return sum(bisect_right(ab, target - c - d) for c in C for d in D)


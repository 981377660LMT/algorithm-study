class Solution:
    def isMonotonic(self, A):
        N = len(A)
        inc, dec = True, True
        for i in range(1, N):
            if A[i] < A[i - 1]:
                inc = False
            if A[i] > A[i - 1]:
                dec = False
            if not inc and not dec:
                return False
        return True


# 如果数组是单调递增或单调递减的，那么它是单调的。
# 一次遍历解法

import numpy as np


class Solution:
    def kInversePairs(self, n: int, k: int) -> int:

        a = np.ones(1, dtype=np.int64)
        print(a)
        for i in range(1, n + 1):
            a = np.convolve(a, np.ones(i, dtype=np.int64)) % 1000000007
            a = a[:1001]
        return int(a[k]) if k < len(a) else 0


print(Solution().kInversePairs(3, 1))
print(Solution().kInversePairs(1, 1))


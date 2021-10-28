import collections


all2power = [collections.Counter(str(1 << i)) for i in range(30)]


class Solution:
    def reorderedPowerOf2(self, N: int) -> bool:
        ncount = collections.Counter(str(N))
        return ncount in all2power


print(Solution().reorderedPowerOf2(16))

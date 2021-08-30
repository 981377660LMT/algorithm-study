import collections


class Solution:
    def reorderedPowerOf2(self, N: int) -> bool:
        all2power = [(2 ** i) for i in range(32)]
        all2powercount = [collections.Counter(str(i)) for i in all2power]
        # print(all2powercount, "\n")
        ncount = collections.Counter(str(N))
        # print(ncount)
        return ncount in all2powercount


print(Solution().reorderedPowerOf2(16))

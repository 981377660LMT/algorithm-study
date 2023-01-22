from collections import Counter


POW2_COUNTER = [Counter(str(1 << i)) for i in range(30)]


class Solution:
    def reorderedPowerOf2(self, N: int) -> bool:
        counter = Counter(str(N))
        return counter in POW2_COUNTER


print(Solution().reorderedPowerOf2(16))

from functools import lru_cache
from itertools import accumulate

MOD = int(1e9 + 7)

# jump or not jump

# 1 <= carpetLen <= floor.length <= 1000
class Solution:
    def minimumWhiteTiles(self, floor: str, numCarpets: int, carpetLen: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            # 每个位置跳还是不跳
            if index >= n or remain == 0:
                return 0

            res = 0
            res = (
                max(res, dfs(index + carpetLen, remain - 1))
                + preSum[min(index + carpetLen, n)]
                - preSum[index]
            )
            res = max(res, dfs(index + 1, remain))
            return res

        n = len(floor)
        nums = [1 if floor[i] == '1' else 0 for i in range(n)]
        preSum = [0] + list(accumulate(nums))
        res = dfs(0, numCarpets)
        dfs.cache_clear()
        return floor.count('1') - res


print(Solution().minimumWhiteTiles(floor="101111010", numCarpets=2, carpetLen=4))

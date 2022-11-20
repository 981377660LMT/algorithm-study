from functools import lru_cache

MOD = int(1e9 + 7)
INF = int(1e20)

# 'W' -> 'B'


class Solution:
    def minimumRecolors(self, blocks: str, k: int) -> int:
        # 每一次操作中，你可以选择一个白色块将它 涂成 黑色块。
        # 请你返回至少出现 一次 连续 k 个黑色块的 最少 操作次数。
        # !注意这里是窗口 不要看到连续就想groupby
        # !枚举每个答案窗口的左端点即可
        n = len(blocks)
        res = n
        preSum = [0]
        for char in blocks:
            preSum.append(preSum[-1] + (1 if char == "B" else 0))
        for start in range(n - k + 1):
            one = preSum[start + k] - preSum[start]
            res = min(res, k - one)
        return res

    def minimumRecolors2(self, blocks: str, k: int) -> int:
        """定长滑窗用双指针"""
        n = len(blocks)
        res, black = n, 0
        for right in range(n):
            black += blocks[right] == "B"
            if right >= k:
                black -= blocks[right - k] == "B"
            if right >= k - 1:
                res = min(res, k - black)
        return res


print(Solution().minimumRecolors(blocks="WBBWWBBWBW", k=7))
print(Solution().minimumRecolors(blocks="WBWBBBW", k=2))

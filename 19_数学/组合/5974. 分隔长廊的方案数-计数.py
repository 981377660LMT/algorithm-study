# 请你将走廊用屏风划分为若干段，且每一段内都 恰好有两个座位
class Solution:
    # 统计第 2*k 个座位和第 2*k+1 个座位之间有多少个植物即可。
    # 一个坑：没有座位时算作没有方案

    def numberOfWays(self, corridor: str) -> int:
        seats = [i for i, char in enumerate(corridor) if char == 'S']
        if not seats or len(seats) & 1:
            return 0
        res, mod = 1, int(1e9 + 7)
        for i in range(2, len(seats), 2):
            res *= seats[i] - seats[i - 1]
            res %= mod
        return res


print(Solution().numberOfWays("SSPPSPS"))
print(Solution().numberOfWays("PPSPSP"))
print(Solution().numberOfWays("S"))

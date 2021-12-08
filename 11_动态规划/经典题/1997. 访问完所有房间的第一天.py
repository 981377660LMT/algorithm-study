from typing import List


# 如果算上本次访问，访问 i 号房间的次数为 奇数 ，那么 第二天 需要访问 nextVisit[i] 所指定的房间，
# 如果算上本次访问，访问 i 号房间的次数为 偶数 ，那么 第二天 需要访问 (i + 1) mod n 号房间。
# 请返回你访问完所有房间的第一天的日期编号

MOD = int(1e9 + 7)
# 2 <= n <= 105
# 0 <= nextVisit[i] <= i

#  dp[i] is number of days to reach cell i
#  `We can only reach cell i from the cell i-1`
# summary:
# start => i-1 =>nextVisit[i-1] => i-1 => i
class Solution:
    def firstDayBeenInAllRooms(self, nextVisit: List[int]) -> int:
        n = len(nextVisit)
        dp = [0] * n
        for i in range(1, n):
            dp[i] = ((dp[i - 1] + 1) + (dp[i - 1] - dp[nextVisit[i - 1]]) + 1) % MOD
        return dp[-1]


print(Solution().firstDayBeenInAllRooms([0, 0]))
# 输出：2
# 解释：
# - 第 0 天，你访问房间 0 。访问 0 号房间的总次数为 1 ，次数为奇数。
#   下一天你需要访问房间的编号是 nextVisit[0] = 0
# - 第 1 天，你访问房间 0 。访问 0 号房间的总次数为 2 ，次数为偶数。
#   下一天你需要访问房间的编号是 (0 + 1) mod 2 = 1
# - 第 2 天，你访问房间 1 。这是你第一次完成访问所有房间的那天。

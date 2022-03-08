# 一个商人穿过一个 N×N 的正方形的网格，去参加一个非常重要的商务活动。
# 他要从网格的左上角进，右下角出。
# 每穿越中间 1 个小方格，都要花费 1 个单位时间。
# 商人必须在 (2N−1) 个单位时间穿越出去。
# 而在经过中间的每个小方格时，都需要缴纳一定的费用。
# 这个商人期望在规定时间内用最少费用穿越出去。
# 请问至少需要多少费用？

from typing import List


class Solution:
    def business(self, group: List[List[int]]):
        dp = [[float('inf') for _ in range(len(group[0]) + 1)] for _ in range(len(group) + 1)]
        dp[0][1], dp[1][0] = 0, 0
        for i in range(1, len(dp)):
            for j in range(1, len(dp[0])):
                dp[i][j] = min(dp[i - 1][j], dp[i][j - 1]) + group[i - 1][j - 1]
        return dp[-1][-1]


if __name__ == '__main__':
    solution = Solution()
    row = int(input())
    group = []
    for i in range(row):
        group.append(list(map(int, input().split())))
    res = solution.business(group)
    print(res)


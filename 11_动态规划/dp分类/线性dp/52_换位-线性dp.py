#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# 计算最小交换次数
# @param firstRow int一维数组 第一行的数列数据
# @param secondRow int一维数组 第二行的数列数据
# @return int
#
from typing import List

# 给出两个数列，求使第一个数列严格降序，第二个数列严格升序的最少交换次数
# 若无法交换成严格有序，返回“-1”

# 交换问题：四种情况
class Solution:
    def arrange(self, firstRow: List[int], secondRow: List[int]):
        n = len(firstRow)
        INF = float("inf")
        dp = [[INF, INF] for _ in range(n)]
        dp[0][0], dp[0][1] = 0, 1
        for i in range(1, n):
            # (i - 1)列不换, i列不换
            if firstRow[i - 1] > firstRow[i] and secondRow[i - 1] < secondRow[i]:
                dp[i][0] = min(dp[i][0], dp[i - 1][0])
            # (i - 1)列不换, i列换
            if firstRow[i - 1] > secondRow[i] and secondRow[i - 1] < firstRow[i]:
                dp[i][1] = min(dp[i][1], dp[i - 1][0] + 1)
            # (i - 1)列换, i列不换
            if secondRow[i - 1] > firstRow[i] and firstRow[i - 1] < secondRow[i]:
                dp[i][0] = min(dp[i][0], dp[i - 1][1])
            # (i - 1)列换, i列换
            if secondRow[i - 1] > secondRow[i] and firstRow[i - 1] < firstRow[i]:
                dp[i][1] = min(dp[i][1], dp[i - 1][1] + 1)
        # print(dp)
        return min(dp[n - 1]) if dp[n - 1][0] != INF or dp[n - 1][1] != INF else -1

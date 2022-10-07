# E - NarrowRectangles
# https://atcoder.jp/contests/arc070/tasks/arc070_c

# !dp[i][x] 表示第i个长方形的左端点为x时的最小移动花费
# !dp[i][x]=min(dp[i-1][y]) + abs(x-lefti) , x-(lengthi)<=y<=x+(lengthi) 其中lengthi为第i个长方形的长度
from typing import List
from SlopeTrick import SlopeTrick


class Solution:
    def buildBridge(self, num: int, wood: List[List[int]]) -> int:
        S = SlopeTrick()


if __name__ == "__main__":
    n = int(input())
    woods = [tuple(map(int, input().split())) for _ in range(n)]  # (left,right)

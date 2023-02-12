# E - NarrowRectangles
# https://atcoder.jp/contests/arc070/tasks/arc070_c
# !dp[i][x] 表示第i个长方形的左端点为x时的最小移动花费
# !dp[i][x]=min(dp[i-1][y]) + abs(x-lefti) , x-(lengthi-1)<=y<=x+(lengthi) 其中lengthi为第i个长方形的长度

from typing import List
from SlopeTrick import SlopeTrick


class Solution:
    def buildBridge(self, _num: int, wood: List[List[int]]) -> int:
        """
        Args:
            _num (int): 长度为 _num 的河流 (这个参数其实没用)
            wood (List[List[int]]): wood[i] 记录了第 i 条河道上的浮木初始的覆盖范围。

        Returns:
            int: 勇者跨越这条河流，最少需要花费多少「自然之力」。

        - _num <= 1e9
        - len(wood) <= 1e5
        - wood[i][0] <= wood[i][1] <= 1e9
        """
        st = SlopeTrick()
        length = [b - a for a, b in wood]
        for i, (left, _) in enumerate(wood):
            st.add_l -= length[i]
            if i > 0:
                st.add_r += length[i - 1]
            st.add_abs(left)
        return st.query()[0]


print(Solution().buildBridge(_num=10, wood=[[1, 2], [4, 7], [8, 9]]))
# 将 [1,2] 浮木移动至 [3,4]，花费 2「自然之力」，
# 将 [8,9] 浮木移动至 [7,8]，花费 1「自然之力」，
# 此时勇者可以顺着 [3,4]->[4,7]->[7,8] 跨越河流，
# 因此，勇者最少需要花费 3 点「自然之力」跨越这条河流

if __name__ == "__main__":
    n = int(input())
    woods = [list(map(int, input().split())) for _ in range(n)]  # (left,right)
    print(Solution().buildBridge(n, woods))

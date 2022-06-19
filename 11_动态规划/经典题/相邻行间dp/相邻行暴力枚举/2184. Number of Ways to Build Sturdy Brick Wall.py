from collections import defaultdict
from typing import List

MOD = int(1e9 + 7)

# 首先通过dfs构造出所有合法的瓷砖排列情况，每种情况只需要记录下间隙的位置（代码中记录包含了两端点）
# 所有情况两两作检查，看是否能够间隔排放，用于后续的dp操作
# dp过程：计算以某种排列情况作为当前排所产生的全局排列情况数
# 初始化第一排每种情况数为1
# 每一排选取某种排列的总答案数为上一排与自己能够相邻的排列的答案总和
# 返回最后一排的dp答案总和
# https://leetcode-cn.com/problems/number-of-ways-to-build-sturdy-brick-wall/solution/cdfsgou-zao-dpji-suan-da-an-96ms-by-mono-bfu1/


class Solution:
    def buildWall(self, height: int, width: int, bricks: List[int]) -> int:
        """预处理每行可能的状态后，相邻行间进行dp
        
        滚动数组快了很多
        """

        def dfs(curWidth: int, state: int) -> None:
            if curWidth > width:
                return
            if curWidth == width:
                allStates.add(state)
                return
            for choose in bricks:
                nextWidth = curWidth + choose
                nextState = state
                if nextWidth != width:
                    nextState |= 1 << (nextWidth)
                dfs(nextWidth, nextState)

        bricks = sorted(bricks)
        allStates = set()
        dfs(0, 0)

        dp = defaultdict(int)
        for cur in allStates:
            dp[cur] = 1
        for _ in range(1, height):
            ndp = defaultdict(int)
            for cur in allStates:
                for pre in dp:
                    if not cur & pre:
                        ndp[cur] += dp[pre]
                        ndp[cur] %= MOD
            dp = ndp

        res = 0
        for key in dp:
            res += dp[key]
            res %= MOD
        return res

from functools import lru_cache
from itertools import groupby
from typing import Set


def mex(nums: Set[int]) -> int:
    res = 0
    while res in nums:
        res += 1
    return res


class Solution:
    def canWin(self, currentState: str) -> bool:
        """
        你和朋友轮流将 连续 的两个 "++" 反转成 "--" 。当一方无法进行有效的翻转时便意味着游戏结束，则另一方获胜
        1 <= currentState.length <= 60

        O(n^2)
        """

        @lru_cache(None)
        def grundy(plusLen: int) -> int:
            """状态定义:连续加号的长度"""
            nexts = set()
            for i in range(plusLen - 1):
                remain = plusLen - i - 2
                nexts.add(grundy(remain) ^ grundy(i))
            return mex(nexts)

        groups = [len(list(g)) for c, g in groupby(currentState) if c == "+"]
        sg = 0
        for len_ in groups:
            sg ^= grundy(len_)
        return sg > 0


print(Solution().canWin("++++"))

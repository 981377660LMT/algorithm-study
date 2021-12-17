from typing import List
from itertools import chain

# 1 <= m * n <= 2 * 105 (暗示遍历m*n)
# 如果单词是 水平 放置的，那么该单词左边和右边 相邻 格子不能为 ' ' 或小写英文字母。
# 如果单词是 竖直 放置的，那么该单词上边和下边 相邻 格子不能为 ' ' 或小写英文字母。

# 看到还是没什么思路

# 每一行和每一列的字符串按 '#' 分割为子串，判断每个子串是否匹配 word 或 word[::-1] 即可。
class Solution:
    def placeWordInCrossword(self, board: List[List[str]], word: str) -> bool:
        def canMatch(s: str, t: str) -> bool:
            """"s在对应位置要么是空格要么相等"""

            return all(a in [' ', b] for a, b in zip(s, t))

        for row in chain(board, zip(*board)):
            for space in ''.join(row).split('#'):
                if len(space) == len(word) and (
                    canMatch(space, word) or canMatch(space, word[::-1])
                ):
                    return True
        return False


print(
    Solution().placeWordInCrossword([["#", " ", "#"], [" ", " ", "#"], ["#", "c", " "]], word="abc")
)

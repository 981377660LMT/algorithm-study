from typing import List
from collections import defaultdict

# 所有的单词长度都相同。
# 找出其中所有的 单词方块 (沿正对角线对称)
# 1 <= words.length <= 1000
# 1 <= words[i].length <= 4
# words[i] 长度相同
# words[i] 只由小写英文字母组成
# words[i] 都 各不相同


class Solution:
    def wordSquares(self, words: List[str]) -> List[List[str]]:
        def bt(row: int, path: List[str]) -> None:
            """当前要选择第几行的人选了,已经选择的单词"""
            if len(path) == COL:
                res.append(path[:])
                return

            col = ''.join([w[row] for w in path])
            for next in adjMap[col]:
                path.append(next)
                bt(row + 1, path)
                path.pop()

        # !根据前缀查询单词
        adjMap = defaultdict(set)
        for w in words:
            for i in range(1, len(w)):
                adjMap[w[:i]].add(w)

        res = []
        COL = len(words[0])
        for w in words:
            bt(1, [w])
        return res


print(Solution().wordSquares(["area", "lead", "wall", "lady", "ball"]))
# 输入：
# ["area","lead","wall","lady","ball"]

# 输出：
# [
#   [ "wall",
#     "area",
#     "lead",
#     "lady"
#   ],
#   [ "ball",
#     "area",
#     "lead",
#     "lady"
#   ]
# ]


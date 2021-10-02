from typing import List
from collections import defaultdict

# 所有的单词长度都相同。
# 找出其中所有的 单词方块 (沿正对角线对称)
class Solution:
    def wordSquares(self, words: List[str]) -> List[List[str]]:
        def helper(path: List[str], index: int):
            if len(path[0]) == index:
                return res.append(path[::])
            # 看第index列所有单词组成的前缀
            prefix = ''.join([x[index] for x in path])
            for next_word in dic[prefix]:
                path.append(next_word)
                helper(path, index + 1)
                path.pop()

        # 1.找到前缀树中所有前缀满足要求的词
        res = []
        dic = defaultdict(list)
        for word in words:
            for i in range(1, len(word)):
                dic[word[:i]].append(word)
        print(dic)
        # 2.每次看第index列所有单词组成的前缀
        for word in words:
            helper([word], 1)
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


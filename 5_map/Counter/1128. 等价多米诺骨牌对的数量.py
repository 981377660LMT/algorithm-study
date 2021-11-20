# 1 <= dominoes.length <= 40000
# 1 <= dominoes[i][j] <= 9
from typing import List
from collections import Counter
from math import comb


class Solution:
    def numEquivDominoPairs(self, dominoes: List[List[int]]) -> int:
        # 可替换为位运算进行压缩
        states = [tuple(sorted(entry)) for entry in dominoes]
        counter = Counter(states)
        return sum(comb(count, 2) for count in counter.values())


# 输入：dominoes = [[1,2],[2,1],[3,4],[5,6]]
# 输出：1

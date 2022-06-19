# 给你一个字符串数组 ideas 表示在公司命名过程中使用的名字列表。公司命名流程如下：

# 从 ideas 中选择 2 个 不同 名字，称为 ideaA 和 ideaB 。
# 交换 ideaA 和 ideaB 的首字母。
# 如果得到的两个新名字 都 不在 ideas 中，那么 ideaA ideaB（串联 ideaA 和 ideaB ，中间用一个空格分隔）是一个有效的公司名字。
# 否则，不是一个有效的名字。
# 返回 不同 且有效的公司名字的数目。


from collections import defaultdict
from itertools import product
from typing import List


class Solution:
    def distinctNames(self, ideas: List[str]) -> int:
        """直接遍历所有的 otherother 来统计答案显然会是 O(n^2) 会超时

        !从枚举首字母的角度来计算所有组合
        """
        adjMap = defaultdict(set)
        for w in ideas:
            adjMap[w[0]].add(w[1:])

        res = 0
        for c1, c2 in product(adjMap, repeat=2):
            if c1 == c2:
                continue
            sameSuf = len(adjMap[c1] & adjMap[c2])
            res += (len(adjMap[c1]) - sameSuf) * (len(adjMap[c2]) - sameSuf)
        return res


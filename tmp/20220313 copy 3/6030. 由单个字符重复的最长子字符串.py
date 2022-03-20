from itertools import groupby
from typing import List, Tuple
from collections import defaultdict


from sortedcontainers import SortedList, SortedDict

# 1 <= s.length <= 105

# 1 <= k <= 105


class Solution:
    def longestRepeating(self, s: str, queryCharacters: str, queryIndices: List[int]) -> List[int]:
        """
        第 i 个查询会将 s 中位于下标 queryIndices[i] 的字符更新为 queryCharacters[i] 。
        返回一个长度为 k 的数组 lengths ，其中 lengths[i] 是在执行第 i 个查询 之后 s 中仅由 单个字符重复 组成的 最长子字符串 的 长度 。
        """

        k = len(queryCharacters)
        res = [1] * k
        endsWith = {}
        startWith = {}
        sls = [SortedList(key=lambda x: -x[0]) for _ in range(26)]

        groups = [[char, len(list(group))] for char, group in groupby(s)]

        # print(groups)
        start, end = 0, 0
        for group in groups:
            char, length = group
            sls[ord(char) - 97].add(length)
            end += length - 1
            endsWith[ord(char) - 97][end] = length
            startWith[ord(char) - 97][start] = length
            start = end + 1
            end = start

        for i in range(k):
            qc, qi = queryCharacters[i], queryIndices[i]
            if qc == s[qi]:
                res[i] = sls[ord(qc) - 97][0]
                continue

            # 中断
            if qi not in endsWith and qi not in startWith:
                preIndex = endsWith.bisect_right(qi)
                preItem = endsWith[preIndex]
                preChar = preItem[1]
            # 连接
            pre, next = qi - 1, qi + 1
            # if

        return res


print(Solution().longestRepeating(s="babacc", queryCharacters="bcb", queryIndices=[1, 3, 3]))

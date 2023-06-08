# 子序列自动机
# 如果需要进行大量的子序列匹配，那么就不能用朴素的双指针匹配了

# 1.nexts数组形式 26*n
# !O(26*n) 预处理 O(s) 查询

# 2.二分形式
# !O(nlogn) 预处理 O(slogn) 查询

from bisect import bisect_right
from collections import defaultdict
from typing import DefaultDict, List, Tuple


# !默认charset为26个小写字母


class SubsequenceAutomaton1:
    def __init__(self, s: str) -> None:
        self._s = s
        self._nexts = self._build()
        """
        _nexts[i][j] 表示在 i 右侧的字符 j 的最近位置 (右侧表示下标严格大于i).
        如果不存在，则为 n.
        """

    def match(self, t: str, sStart=0, tStart=0) -> Tuple[int, int]:
        """在 s[sStart:] 的子序列中寻找 t[tStart:]

        :param sStart: s的起始索引
        :param tStart: t的起始索引
        :return: (hit,end) (匹配的前缀长度, 匹配到的前缀对应在s中的结束索引)
        """
        if not self._s or not t:
            return 0, 0
        n, m = len(self._s), len(t)
        si, ti = sStart, tStart
        if self._s[sStart] == t[tStart]:
            ti += 1
        while si < n and ti < m:
            nextPos = self._nexts[si][ord(t[ti]) - 97]
            if nextPos == n:
                return ti - tStart, si
            si, ti = nextPos, ti + 1
        return ti - tStart, si

    def _build(self) -> List[Tuple[int]]:
        n = len(self._s)
        nexts = [None] * n
        last = [n] * 26
        for i in range(n - 1, -1, -1):
            nexts[i] = tuple(last)  # type: ignore
            last[ord(self._s[i]) - 97] = i
        return nexts  # type: ignore


class SubsequenceAutomaton2:
    def __init__(self, s: str) -> None:
        self._s = s
        self._indexes = self._build()

    def match(self, t: str, sStart=0, tStart=0) -> Tuple[int, int]:
        """在 s[sStart:] 的子序列中寻找 t[tStart:]

        :param sStart: s的起始索引
        :param tStart: t的起始索引
        :return: (hit,end) (匹配的前缀长度, 匹配到的前缀对应在s中的结束索引)
        """
        if not self._s or not t:
            return 0, 0
        n, m = len(self._s), len(t)
        si, ti = sStart, tStart
        if self._s[sStart] == t[tStart]:
            ti += 1
        while si < n and ti < m:
            indexes = self._indexes[t[ti]]
            pos = bisect_right(indexes, si)
            if pos == len(indexes):
                return ti - tStart, si
            si, ti = indexes[pos], ti + 1
        return ti - tStart, si

    def _build(self) -> DefaultDict[str, List[int]]:
        indexes = defaultdict(list)
        for i, char in enumerate(self._s):
            indexes[char].append(i)
        return indexes


if __name__ == "__main__":

    sa = SubsequenceAutomaton1("abcdebdde")
    assert sa.match("bde") == (3, 4)
    assert sa.match("bde", 1) == (3, 4)

    sa = SubsequenceAutomaton2("bbabbabbbbabaababab")
    print(sa.match("bbbbbbbbbbbb"))

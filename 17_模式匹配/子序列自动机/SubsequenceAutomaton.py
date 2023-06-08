# 子序列自动机
# 如果需要进行大量的子序列匹配，那么就不能用朴素的双指针匹配了
# 1.nexts数组形式 26*n
# !O(26*n) 预处理 O(s) 查询
# 2.二分形式
# !O(nlogn) 预处理 O(slogn) 查询
# 查询当前位置的下一个特定字符的位置
# 查询是否含有某序列


from bisect import bisect_right
from collections import defaultdict
from typing import DefaultDict, Generic, List, Sequence, Tuple, TypeVar


class SubsequenceAutomaton1:
    __slots__ = ("_s", "_nexts", "_charset", "_offset")

    def __init__(self, s: str, charset=26, offset=97) -> None:
        """O(charset*n) 预处理.

        Args:
            s (str): 待匹配的字符串
            charset (int, optional): 字符集大小. 默认为 26.
            offset (int, optional): 字符集的起始字符. 默认为 97.
        """
        self._s = s
        self._charset = charset
        self._offset = offset
        self._nexts = self._build()
        """
        _nexts[i][j] 表示在 i 右侧的字符 j 的最近位置 (右侧表示下标严格大于i).
        如果不存在，则为 n.
        """

    def move(self, pos: int, char: str) -> int:
        """
        查询当前位置的下一个特定字符的位置(下标严格大于pos).
        如果不存在，则为 n.
        0<=pos<n.
        """
        return self._nexts[pos][ord(char) - self._offset]

    def includes(self, t: str, sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> bool:
        """
        查询s[sStart:sEnd]是否含有某序列t[tStart:tEnd].
        时间复杂度O(len(t)).
        """
        hit, _ = self.match(t, sStart=sStart, sEnd=sEnd, tStart=tStart, tEnd=tEnd)
        if tEnd == -1:
            tEnd = len(t)
        return hit >= tEnd - tStart

    def match(self, t: str, sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> Tuple[int, int]:
        """
        在 s[sStart:sEnd] 中寻找子序列 t[tStart:tEnd].
        时间复杂度 O(len(t)).

        Args:
            t: 待匹配的子序列
            sStart: s的起始索引
            sEnd: s的结束索引
            tStart: t的起始索引
            tEnd: t的结束索引
        Returns:
            (hit,end): (`匹配到的t的长度`, `匹配结束时s的索引`)
        """
        if sEnd == -1:
            sEnd = len(self._s)
        if sStart >= sEnd:
            return 0, sStart
        if tEnd == -1:
            tEnd = len(t)
        if tStart >= tEnd:
            return 0, sStart

        n = len(self._s)
        si, ti = sStart, tStart
        if self._s[sStart] == t[tStart]:
            ti += 1
        while si < sEnd and ti < tEnd:
            nextPos = self.move(si, t[ti])
            if nextPos == n:
                return ti - tStart, si
            si, ti = nextPos, ti + 1
        return ti - tStart, si

    def _build(self) -> List[Tuple[int]]:
        n = len(self._s)
        nexts = [None] * n
        last = [n] * self._charset
        offset = self._offset
        for i in range(n - 1, -1, -1):
            nexts[i] = tuple(last)  # type: ignore
            last[ord(self._s[i]) - offset] = i
        return nexts  # type: ignore


V = TypeVar("V")


class SubsequenceAutomaton2(Generic[V]):
    __slots__ = ("_seq", "_indexes")

    def __init__(self, seq: Sequence[V]) -> None:
        """O(nlogn) 预处理."""
        self._seq = seq
        self._indexes = self._build()

    def move(self, pos: int, newValue: V) -> int:
        """
        查询当前位置的下一个特定字符的位置(下标严格大于pos).
        如果不存在，则为 n.
        0<=pos<n
        """
        indexes = self._indexes[newValue]
        nextPos = bisect_right(indexes, pos)
        return indexes[nextPos] if nextPos < len(indexes) else len(self._seq)

    def includes(self, t: Sequence[V], sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> bool:
        """
        查询s[sStart:sEnd]是否含有某序列t[tStart:tEnd].
        时间复杂度O(len(t)logn).
        """
        hit, _ = self.match(t, sStart=sStart, sEnd=sEnd, tStart=tStart, tEnd=tEnd)
        if tEnd == -1:
            tEnd = len(t)
        return hit >= tEnd - tStart

    def match(self, t: Sequence[V], sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> Tuple[int, int]:
        """
        在 s[sStart:sEnd] 中寻找子序列 t[tStart:tEnd].
        时间复杂度 O(len(t)logn).

        Args:
            t: 待匹配的子序列
            sStart: s的起始索引
            sEnd: s的结束索引
            tStart: t的起始索引
            tEnd: t的结束索引
        Returns:
            (hit,end): (`匹配到的的t的长度`, `匹配结束时s的索引`)
        """
        if sEnd == -1:
            sEnd = len(self._seq)
        if sStart >= sEnd:
            return 0, sStart
        if tEnd == -1:
            tEnd = len(t)
        if tStart >= tEnd:
            return 0, sStart

        n = len(self._seq)
        si, ti = sStart, tStart
        if self._seq[sStart] == t[tStart]:
            ti += 1
        while si < sEnd and ti < tEnd:
            nextPos = self.move(si, t[ti])
            if nextPos == n:
                return ti - tStart, si
            si, ti = nextPos, ti + 1
        return ti - tStart, si

    def _build(self) -> DefaultDict[V, List[int]]:
        indexes = defaultdict(list)
        for i, char in enumerate(self._seq):
            indexes[char].append(i)
        return indexes


if __name__ == "__main__":
    sa = SubsequenceAutomaton1("abcdebdde")
    assert sa.match("bde") == (3, 4)
    assert sa.match("bde", 1) == (3, 4)

    sa = SubsequenceAutomaton1("bbabbabbbbabaababab")

    assert sa.match("bbbbbbbbbbbb") == (12, 18)
    assert sa.includes("bbbbbbbbbbbb")

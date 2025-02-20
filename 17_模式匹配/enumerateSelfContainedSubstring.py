# - 定义:
#   自包含子串(SelfContainedSubstring)是指，该子串中所有的字符，均未在子串以外的部分出现。
#   如abcbca，其中"bcbc"包含b、c两种字符都不在子串以外的部分出现，因此是自包含子串。
#
# - 性质：
#   !1.自包含子串最多只有O(∑)个可能的起点(每个字母第一次出现的位置).
#   !2.自包含子串最多只有O(∑^2)个.
#   !3.自包含子串当且仅当"子串内每种字符的次数之和等于子串长度".


from typing import List, Tuple, Union
from operator import itemgetter


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


def enumerateSelfContainedSubstring(s: str) -> List[Tuple[int, int]]:
    """
    返回所有自包含子串.
    自包含子串是指，该子串中所有的字符，均未在子串以外的部分出现.
    时间复杂度: O(n + ∑^2).

    Args:
        s: 字符串s.
    Returns:
        List[Tuple[int, int]]: 所有自包含子串的起始和结束位置(左闭右开).
    """
    first, last, counter = dict(), dict(), dict()
    for i, c in enumerate(s):
        if c not in first:
            first[c] = i
        last[c] = i
        counter[c] = counter.get(c, 0) + 1
    chars = sorted(first, key=lambda x: first[x])  # 将字符按首次出现的位置排序

    m = len(chars)
    res = []
    for i, c1 in enumerate(chars):  # 枚举起始字符
        left, right, count = first[c1], 0, 0
        for j in range(i, m):  # 枚举结束字符
            c2 = chars[j]
            right = max2(right, last[c2])
            count += counter[c2]
            if count == right - left + 1:
                res.append((left, right + 1))
    return res


def enumerateNonOverlappingSelfContainedSubstring(s: str) -> List[Tuple[int, int]]:
    """
    返回最多的不相交的自包含子串.
    如果有多个解法有相同的子串数目，返回这些子串总长度最小的一个解.
    时间复杂度: O(n + ∑log∑).

    Args:
        s: 字符串s.
    Returns:
        List[Tuple[int, int]]: 最多不相交的自包含子串的起始和结束位置(左闭右开).
    """

    first, last, counter = dict(), dict(), dict()
    for i, c in enumerate(s):
        first.setdefault(c, i)
        last[c] = i
        counter[c] = counter.get(c, 0) + 1
    intervals = [[first[c], last[c], counter[c]] for c in counter]
    intervals.sort(key=itemgetter(1))

    res = []
    merged = []
    for interval in intervals:
        merged.append(interval)
        # 合并有重叠区间
        while len(merged) > 1 and merged[-2][1] > merged[-1][0]:
            l, r, c = merged.pop()
            merged[-1][0] = min2(merged[-1][0], l)
            merged[-1][1] = max2(merged[-1][1], r)
            merged[-1][2] += c
        # 判断无重叠区间是否满足条件，如满足条件则立刻删除
        if merged[-1][1] - merged[-1][0] + 1 == merged[-1][2]:
            a, b, _ = merged.pop()
            res.append((a, b + 1))
    return res


if __name__ == "__main__":

    class Solution:
        # 1520. 最多的不重叠子字符串
        # https://leetcode.cn/problems/maximum-number-of-non-overlapping-substrings/
        # 给定字符串s，找到最多不相交的自包含子串，返回这些子串.
        def maxNumOfSubstrings(self, s: str) -> List[str]:
            return [s[start:end] for start, end in enumerateNonOverlappingSelfContainedSubstring(s)]
            intervals = enumerateSelfContainedSubstring(s)
            intervals.sort(key=lambda x: x[1] - x[0])
            indexes = self.maxNonIntersectingIntervals(intervals)
            res = []
            for i in indexes:
                start, end = intervals[i]
                res.append(s[start:end])
            return res

        # 3104. Find Longest Self-Contained Substring (查找最长的自包含子字符串)
        # https://leetcode.cn/problems/find-longest-self-contained-substring/description/
        # 给定字符串s，找到最长的自包含子字符串，返回其长度，没有则返回-1。
        def maxSubstringLength_3104(self, s: str) -> int:
            cand = enumerateSelfContainedSubstring(s)
            res = -1
            for start, end in cand:
                tmp = end - start
                if tmp != len(s):
                    res = max2(res, tmp)
            return res

        # https://leetcode.cn/problems/select-k-disjoint-special-substrings/description/
        # 3458. 选择 K 个互不重叠的特殊子字符串
        def maxSubstringLength(self, s: str, k: int) -> bool:
            n = len(s)
            intervals = enumerateSelfContainedSubstring(s)
            intervals = [(s, e) for s, e in intervals if e - s != n]
            res = self.maxNonIntersectingIntervals(intervals)
            return len(res) >= k

        def maxNonIntersectingIntervals(
            self,
            intervals: Union[List[Tuple[int, int]], List[List[int]]],
            allowOverlapping=True,
            endInclusive=True,
        ) -> List[int]:
            """
            给定 n 个区间 [left_i,right_i].
            请你在数轴上选择若干区间,使得选中的区间之间互不相交.

            Args:
                intervals: 区间列表,每个区间为[left,right].
                allowOverlapping: 是否允许选择的区间端点重合.默认为False.
                endInclusive: 是否包含区间右端点.默认为True.

            Returns:
                List[int]: 区间索引列表. 如果有多种方案，返回字典序最小的那个.
            """

            n = len(intervals)
            if n == 0:
                return []
            if n == 1:
                return [0]
            order = sorted(range(n), key=lambda x: intervals[x][1])
            res = [order[0]]
            preEnd = intervals[order[0]][1]
            for i in order[1:]:
                start, end = intervals[i]
                if not endInclusive:
                    end -= 1
                if allowOverlapping:
                    if start >= preEnd:
                        res.append(i)
                        preEnd = end
                else:
                    if start > preEnd:
                        res.append(i)
                        preEnd = end
            return res

# 616. 给字符串添加加粗标签
# https://leetcode.cn/problems/add-bold-tag-in-string/description/
# 给定字符串 s 和字符串数组 words。
#
# 对于 s 内部的子字符串，若其存在于 words 数组中， 则通过添加闭合的粗体标签 <b> 和 </b> 进行加粗标记。
#
# 如果两个这样的子字符串重叠，你应该仅使用一对闭合的粗体标签将它们包围起来。
# 如果被粗体标签包围的两个子字符串是连续的，你应该将它们合并。
# 返回添加加粗标签后的字符串 s 。


from itertools import accumulate
from typing import Callable, List, Optional, Sequence, TypeVar

T = TypeVar("T", int, str)


def indexOfAll(
    longer: Sequence[T], shorter: Sequence[T], start=0, nexts: Optional[List[int]] = None
) -> List[int]:
    """kmp O(n+m)求搜索串 `longer` 中所有匹配 `shorter` 的位置."""
    if not shorter:
        return []
    if len(longer) < len(shorter):
        return []
    res = []
    next = getNext(shorter) if nexts is None else nexts
    hitJ = 0
    for i in range(start, len(longer)):
        while hitJ > 0 and longer[i] != shorter[hitJ]:
            hitJ = next[hitJ - 1]
        if longer[i] == shorter[hitJ]:
            hitJ += 1
        if hitJ == len(shorter):
            res.append(i - len(shorter) + 1)
            hitJ = next[hitJ - 1]
    return res


def getNext(needle: Sequence[T]) -> List[int]:
    """kmp O(n)求 `needle`串的 `next`数组
    `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度
    """
    next = [0] * len(needle)
    j = 0
    for i in range(1, len(needle)):
        while j and needle[i] != needle[j]:  # 1. fallback后前进：匹配不成功j往右走
            j = next[j - 1]
        if needle[i] == needle[j]:  # 2. 匹配：匹配成功j往右走一步
            j += 1
        next[i] = j
    return next


K = TypeVar("K")


def groupByKey(n: int, key: Callable[[int], K]):
    end = 0
    while end < n:
        start, leader = end, key(end)
        end += 1
        while end < n and key(end) == leader:
            end += 1
        yield start, end, leader


class Solution:
    def addBoldTag(self, s: str, words: List[str]) -> str:
        n = len(s)
        boldDiff = [0] * (n + 1)
        for w in words:
            matches = indexOfAll(s, w)
            for start in matches:
                boldDiff[start] += 1
                boldDiff[start + len(w)] -= 1
        bold = list(accumulate(boldDiff))[:-1]

        res = []
        for start, end, v in groupByKey(n, lambda i: bold[i] > 0):
            if v:
                res.append("<b>")
                res.append(s[start:end])
                res.append("</b>")
            else:
                res.append(s[start:end])
        return "".join(res)


if __name__ == "__main__":
    sol = Solution()
    s1 = "abcxyz123"
    words1 = ["abc", "123"]
    print(sol.addBoldTag(s1, words1))  # "<b>abc</b>xyz<b>123</b>"

    s2 = "aaabbcc"
    words2 = ["aaa", "aab", "bc"]
    # 区间合并后 [0,4) 包含 "aaab", [4,6) 包含 "bc"
    print(sol.addBoldTag(s2, words2))  # "<b>aaabbc</b>c"

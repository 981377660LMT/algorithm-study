class SlidingWindowContainsAll:
    """
    支持新增元素、删除元素、快速判断容器内是否包含了指定的所有元素.
    """

    __slots__ = "_missingCount", "_missing"

    def __init__(self, supplier):
        self._missing = dict()
        for v in supplier:
            self._missing[v] = self._missing.get(v, 0) + 1
        self._missingCount = len(self._missing)

    def add(self, v) -> bool:
        c = self._missing.get(v, None)
        if c is None:
            return False
        self._missing[v] = c - 1
        if c == 1:
            self._missingCount -= 1
        return True

    def discard(self, v) -> bool:
        c = self._missing.get(v, None)
        if c is None:
            return False
        self._missing[v] = c + 1
        if c == 0:
            self._missingCount += 1
        return True

    def containsAll(self) -> bool:
        return self._missingCount == 0


if __name__ == "__main__":
    # 3298. 统计重新排列后包含另一个字符串的子字符串数目 II
    # https://leetcode.cn/problems/count-substrings-that-can-be-rearranged-to-contain-a-string-ii/description/
    class Solution:
        def validSubstringCount(self, word1: str, word2: str) -> int:
            S = SlidingWindowContainsAll(word2)
            res, left, n = 0, 0, len(word1)
            for right in range(n):
                S.add(word1[right])
                while left <= right and S.containsAll():
                    S.discard(word1[left])
                    left += 1
                res += left
            return res

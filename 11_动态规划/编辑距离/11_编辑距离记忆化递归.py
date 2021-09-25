import functools


class Solution:
    @functools.lru_cache(None)
    def minDistance1(self, word1: str, word2: str) -> int:
        if not word1 or not word2:
            return len(word1) + len(word2)
        if word1[0] == word2[0]:
            return self.minDistance1(word1[1:], word2[1:])
        else:
            inserted = 1 + self.minDistance1(word1, word2[1:])
            deleted = 1 + self.minDistance1(word1[1:], word2)
            replace = 1 + self.minDistance1(word1[1:], word2[1:])
            return min(inserted, deleted, replace)

    # 由于字符串切片是 O(n)所以改成用了索引号。
    def minDistance2(self, word1: str, word2: str) -> int:
        if not word1 or not word2:
            return len(word1) + len(word2)

        @functools.lru_cache(None)
        def helper(i: int, j: int) -> int:
            if i == len(word1) or j == len(word2):
                return len(word1) - i + len(word2) - j
            if word1[i] == word2[j]:
                return helper(i + 1, j + 1)
            else:
                inserted = helper(i, j + 1)
                deleted = helper(i + 1, j)
                replace = helper(i + 1, j + 1)
                return min(inserted, deleted, replace) + 1

        return helper(0, 0)

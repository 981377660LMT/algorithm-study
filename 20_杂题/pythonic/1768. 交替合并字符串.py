from itertools import zip_longest


class Solution:
    def mergeAlternately(self, word1: str, word2: str) -> str:
        return ''.join(s1 + s2 for s1, s2 in zip_longest(word1, word2, fillvalue=''))


print(Solution().mergeAlternately("ab", 'cdef'))

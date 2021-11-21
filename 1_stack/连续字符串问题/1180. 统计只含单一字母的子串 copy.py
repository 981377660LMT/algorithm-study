from itertools import groupby


class Solution:
    def countLetters(self, s: str) -> int:
        return sum(c * (c + 1) // 2 for c in (len(list(group)) for _, group in groupby(s)))


print(Solution().countLetters('aaasdee'))

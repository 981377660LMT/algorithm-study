from string import ascii_lowercase
from typing import List


class Solution:
    def wordCount(self, startWords: List[str], targetWords: List[str]) -> int:
        compress = lambda s: sum(1 << (ord(ch) - ord("a")) for ch in s)
        exist = set(map(compress, startWords))

        res = 0
        for w in startWords:
            state = 0
            for char in w:
                state |= 1 << (ord(char) - ord("a"))

            for next in ascii_lowercase:
                if next in w:
                    continue
                cur = state | (1 << (ord(next) - ord("a")))
                exist.add(cur)

        for w in targetWords:
            state = 0
            for char in w:
                state |= 1 << (ord(char) - ord("a"))
            res += int(state in exist)
        return res


# 2 1 4
print(Solution().wordCount(startWords=["ant", "act", "tack"], targetWords=["tack", "act", "acti"]))
print(Solution().wordCount(startWords=["ab", "a"], targetWords=["abc", "abcd"]))
print(
    Solution().wordCount(
        startWords=["q", "ugqm", "o", "ar", "e"],
        targetWords=[
            "nco",
            "mnwhi",
            "tkuw",
            "ugmiq",
            "fb",
            "oykr",
            "us",
            "sra",
            "dxg",
            "dbp",
            "ql",
            "fq",
        ],
    )
)

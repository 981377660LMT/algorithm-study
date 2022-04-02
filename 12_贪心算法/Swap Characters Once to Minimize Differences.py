# 至多交换一次，让两个字符串不相等的位置最少
# if there are d different characters.
# your answer's possibilities are : d -2, d-1 , d
from string import ascii_lowercase


class Solution:
    def solve(self, s, t):
        badPairs = set()
        res = 0

        for (a, b) in zip(s, t):
            if a == b:
                continue
            badPairs.add((a, b))
            res += 1

        for x, y in badPairs:
            if (y, x) in badPairs:
                return res - 2

        for x, _ in badPairs:
            for ch in ascii_lowercase:
                if (ch, x) in badPairs:
                    return res - 1

        return res


print(Solution().solve(s="abbz", t="zcca"))
# We can swap "a" and "z" to turn s into "zbba". Then there's 2 characters that differ between t.

# 分割字符串使得成为若干个不降/不升子串
# 问最少需要分割多少次(最小值为1)


from itertools import groupby


class Solution:
    def solve(self, s):
        if not s:
            return 0

        res = 1
        isIncreasing = None
        for i in range(len(s) - 1):
            if isIncreasing is None and s[i] != s[i + 1]:
                isIncreasing = s[i] < s[i + 1]
            elif (isIncreasing and s[i] > s[i + 1]) or (isIncreasing is False and s[i] < s[i + 1]):
                res += 1
                isIncreasing = None

        return res


class Solution2:
    def solve(self, s):
        # 连续去重
        s = ''.join(char for char, _ in groupby(s))
        n = len(s)

        res = i = 0
        while i < n:
            res += 1
            if i + 1 == n:
                break
            elif s[i + 1] > s[i]:
                while i + 1 < n and s[i + 1] > s[i]:
                    i += 1
            else:
                while i + 1 < n and s[i + 1] < s[i]:
                    i += 1

            i += 1

        return res


print(Solution().solve(s="abcdcba"))
# We can break s into "abcd" + "cba"
# "acccf" is a non-decreasing string, and "bbba" is a non-increasing string.

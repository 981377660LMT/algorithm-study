# 使得字符串前半与后半单调的最小替换次数
# Return the minimum number of characters that need to be updated such that one of the following three conditions is true for all 0 ≤ i < n/2 and n/2 ≤ j < n:

# s[i] > s[j] or
# s[i] < s[j] or
# s[i] == s[j]

# class Solution:For each char character 'a', 'b', ..., 'z', there are 3 cases:

# make every character equal to char;
# make the left half <= char and the right half > char;
# make the left half > char and the right half <= char.

# !枚举分界的字母
from collections import Counter
import string


class Solution:
    def solve(self, s):
        n = len(s)
        left = Counter(s[: n // 2])
        right = Counter(s[n // 2 :])
        res = n

        for char in string.ascii_lowercase:
            # all equal
            res = min(res, n - left[char] - right[char])

            # left <= char < right
            ok = sum(left[c] for c in left if c <= char)
            ok += sum(right[c] for c in right if c > char)
            res = min(res, n - ok)

            # right <= char < left
            ok = sum(left[c] for c in left if c > char)
            ok += sum(right[c] for c in right if c <= char)
            res = min(res, n - ok)

        return res


# 字符串长度为偶数
print(Solution().solve(s="aaabba"))
print(Solution().solve(s="bbbaaa"))
# f we change the last "a" to "b", then we can satisfy the s[i] < s[j] condition.

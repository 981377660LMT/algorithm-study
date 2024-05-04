# 3088. 使字符串反回文
# https://leetcode.cn/problems/make-string-anti-palindrome/description/
# 我们称一个长度为偶数的字符串 s 为 反回文 的，如果对于每一个下标 0 <= i < n ，s[i] != s[n - i - 1]。
# 给定一个字符串 s，你需要进行 任意 次（包括 0）操作使 s 成为 反回文。
# 在一次操作中，你可以选择 s 中的两个字符并且交换它们。
# 返回结果字符串。如果有多个字符串符合条件，返回 字典序最小的那个。如果它不能成为一个反回文，返回 "-1"。
# !排序，前半部分当基准，后半部分对应的每个字母要不一样。


class Solution:
    def makeAntiPalindrome(self, s: str) -> str:
        n = len(s)
        sb = sorted(s)
        m = n >> 1
        if sb[m] != sb[m - 1]:
            return "".join(sb)

        ptr = m
        while ptr < n and sb[ptr] == sb[m - 1]:
            ptr += 1

        # ![m,ptr) 前缀中的一些字符需要从后面[ptr,n) 交换过来
        for i in range(m, n):
            if sb[i] != sb[~i]:
                break
            if ptr >= n:
                return "-1"
            sb[i], sb[ptr] = sb[ptr], sb[i]
            ptr += 1

        return "".join(sb)


# "llml"
print(Solution().makeAntiPalindrome("lmmn"))

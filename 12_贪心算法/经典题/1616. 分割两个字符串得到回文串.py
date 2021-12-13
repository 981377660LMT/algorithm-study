# 1 <= a.length, b.length <= 105
# 将两个字符串都在 相同的下标 分割开。
# 请你判断 aprefix + bsuffix 或者 bprefix + asuffix 能否构成回文串。

# 两端相等，看中间是不是回文
# https://leetcode-cn.com/problems/split-two-strings-to-make-palindrome/solution/xiang-yu-wen-ti-on-by-tsmart/
class Solution:
    def checkPalindromeFormation(self, a: str, b: str) -> bool:
        isPalindrome = lambda s: s == s[::-1]

        def check(a: str, b: str):
            i, j = 0, len(a) - 1
            while i < j and a[i] == b[j]:
                i += 1
                j -= 1
            return isPalindrome(a[i : j + 1]) or isPalindrome(b[i : j + 1])

        return check(a, b) or check(b, a)


print(Solution().checkPalindromeFormation(a="ulacfd", b="jizalu"))
# 输出：true
# 解释：在下标为 3 处分割：
# aprefix = "ula", asuffix = "cfd"
# bprefix = "jiz", bsuffix = "alu"
# 那么 aprefix + bsuffix = "ula" + "alu" = "ulaalu" 是回文串。

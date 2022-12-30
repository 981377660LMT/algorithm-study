# 3 <= s.length <= 2000
# 给你一个字符串 s ，如果可以将它分割成`三个 非空 回文子字符串`，那么返回 true ，否则返回 false 。
# https://leetcode.cn/problems/palindrome-partitioning-iii/ 调用api

# 平方时间复杂度做法很简单,只需要预处理所有子串是否为回文串,然后枚举中间的串判断.
class Solution:
    def checkPartitioning(self, s: str) -> bool:
        if len(s) < 3:
            return False

        def isPalindrome(word):
            return word == word[::-1]

        n = len(s)
        for i in range(1, n):
            if not isPalindrome(s[:i]):
                continue
            for j in range(i + 1, n):
                if isPalindrome(s[i:j]) and isPalindrome(s[j:]):
                    return True
        return False


print(Solution().checkPartitioning(s="abcbdd"))
# 输出：true
# 解释："abcbdd" = "a" + "bcb" + "dd"，三个子字符串都是回文的。

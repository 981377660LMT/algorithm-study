# 3 <= s.length <= 2000
# 给你一个字符串 s ，如果可以将它分割成`三个 非空 回文子字符串`，那么返回 true ，否则返回 false 。
class Solution:
    def checkPartitioning(self, s: str) -> bool:
        if len(s) < 3:
            return False

        def is_backword(word):
            return word == word[::-1]

        n = len(s)
        for i in range(1, n):
            if not is_backword(s[:i]):
                continue
            for j in range(i + 1, n):
                if is_backword(s[i:j]) and is_backword(s[j:]):
                    return True
        return False


print(Solution().checkPartitioning(s="abcbdd"))
# 输出：true
# 解释："abcbdd" = "a" + "bcb" + "dd"，三个子字符串都是回文的。

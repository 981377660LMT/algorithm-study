class Solution(object):
    def validPalindrome(self, s: str):
        """给定一个非空字符串 s,最多`删除`一个字符。判断是否能成为回文字符串。"""

        def check(left: int, right: int, remain: int) -> bool:
            while left < right:
                if s[left] != s[right]:
                    if remain == 0:
                        return False
                    return check(left + 1, right, remain - 1) or check(
                        left, right - 1, remain - 1
                    )
                left, right = left + 1, right - 1
            return True

        return check(0, len(s) - 1, 1)

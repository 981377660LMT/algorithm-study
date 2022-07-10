class Solution:
    def makePalindrome(self, s: str) -> bool:
        """是否能正好进行1次或2次`替换`操作将原串变为回文串"""

        remain = 2
        left, right = 0, len(s) - 1
        while left < right:
            if s[left] != s[right]:
                if remain == 0:
                    return False
                remain -= 1
            left, right = left + 1, right - 1
        return True


print(Solution().makePalindrome(s="zbcfedcba"))

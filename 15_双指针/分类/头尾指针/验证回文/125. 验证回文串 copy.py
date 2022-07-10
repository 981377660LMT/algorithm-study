# 验证回文串
# 只考虑字母和数字字符，可以忽略字母的大小写
class Solution:
    def isPalindrome(self, s: str) -> bool:
        left, right = 0, len(s) - 1
        while left < right:
            while left < right and not s[left].isalnum():
                left += 1
            while left < right and not s[right].isalnum():
                right -= 1
            if not s[left].lower() == s[right].lower():
                return False
            left += 1
            right -= 1
        return True


print("a1".isalnum())  # 数字或字母

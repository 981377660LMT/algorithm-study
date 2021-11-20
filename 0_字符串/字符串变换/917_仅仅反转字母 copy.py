# 字母栈 倒序填充
class Solution:
    def reverseOnlyLetters(self, s: str) -> str:
        alpha = [char for char in s if char.isalpha()]
        print(alpha)
        return ''.join([char if not char.isalpha() else alpha.pop() for char in s])


print(Solution().reverseOnlyLetters("a-bC-dEf-ghIj"))

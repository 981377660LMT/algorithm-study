# 有一个长度为的字符串 ，你可以删除其中的 个字符，使剩余字符串的字典序最小，输出这个剩余字符串。
class Solution(object):
    def removeKdigits(self, num: str, k: int) -> str:
        """字典序最小 栈底肯定是最小的 维护单增的单调栈"""
        assert len(num) >= k
        stack = []
        for letter in num:
            while k and stack and stack[-1] > letter:
                stack.pop()
                k -= 1
            stack.append(letter)
        return ''.join(stack).lstrip('0') or '0'


print(Solution().removeKdigits("1432219", 3))
print(Solution().removeKdigits("10200", 1))

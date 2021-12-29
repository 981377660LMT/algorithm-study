# 若 s[i] 是小写字符，则 s[i+1] 不可以是相同的大写字符。
# 若 s[i] 是大写字符，则 s[i+1] 不可以是相同的小写字符。
# 栈 两个都消掉
class Solution:
    def makeGood(self, s: str) -> str:
        stack = []
        for char in s:
            if stack and ord(stack[-1]) ^ ord(char) == 32:
                stack.pop()
            else:
                stack.append(char)
        return ''.join(stack)


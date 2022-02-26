# 有一个长度为的字符串 ，你可以删除其中的 个字符，使剩余字符串的字典序最小，输出这个剩余字符串。
class Solution(object):
    def removeKdigits(self, num: str, k: int) -> str:
        assert len(num) >= k
        monotoneStack = []
        for letter in num:
            while k and monotoneStack and monotoneStack[-1] > letter:
                monotoneStack.pop()
                k -= 1
            monotoneStack.append(letter)
        return ''.join(monotoneStack).lstrip('0') or '0'


print(Solution().removeKdigits("1432219", 3))
print(Solution().removeKdigits("10200", 1))

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

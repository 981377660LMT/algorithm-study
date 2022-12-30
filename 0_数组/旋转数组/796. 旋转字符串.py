# 给定两个字符串, s 和 goal。
# 如果在若干次旋转操作之后，s 能变成 goal ，那么返回 true 。
# !旋转字符串/轮转字符串
class Solution:
    def rotateString(self, s1: str, s2: str) -> bool:
        if len(s1) != len(s2):
            return False
        return s2 in s1 + s1


assert Solution().rotateString(s1="abcde", s2="cdeab")

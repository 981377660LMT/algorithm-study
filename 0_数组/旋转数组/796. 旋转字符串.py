# 给定两个字符串, s 和 goal。如果在若干次旋转操作之后，s 能变成 goal ，那么返回 true 。
class Solution:
    def rotateString(self, s: str, goal: str) -> bool:
        return goal in (s + s)[1:-1] and len(s) == len(goal)


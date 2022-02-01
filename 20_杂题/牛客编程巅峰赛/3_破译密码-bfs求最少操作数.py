# 每次变换可以选择当前字符串中的一个位置，
# 然后剩下的三个位置的字符从左到右分别加上2，3，5，若是超出'z'，
# 则重新从'a'开始
# 求出需要破译的密码为从s1变换到s2`最少需要的变换次数`。
# bfs求最小深度
from collections import deque


class Solution:
    def solve(self, s1, s2):
        # write code here
        start = tuple(ord(char) - 97 for char in s1)
        target = tuple(ord(char) - 97 for char in s2)
        queue = deque([[0, start]])
        visited = set([start])

        while queue:
            step, cur = queue.popleft()
            if cur == target:
                return step
            a, b, c, d = cur
            next1 = ((a + 0) % 26, (b + 2) % 26, (c + 3) % 26, (d + 5) % 26)
            next2 = ((a + 2) % 26, (b + 0) % 26, (c + 3) % 26, (d + 5) % 26)
            next3 = ((a + 2) % 26, (b + 3) % 26, (c + 0) % 26, (d + 5) % 26)
            next4 = ((a + 2) % 26, (b + 3) % 26, (c + 5) % 26, (d + 0) % 26)
            for next in (next1, next2, next3, next4):
                if next not in visited:
                    visited.add(next)
                    queue.append([step + 1, next])
        return -1


print(Solution().solve("aaaa", "ccgk"))
# 返回值：
# 2

# 说明：
# 第一次变换选择第一个'a'，变成"acdf"，第二次变换选择第二个'c'，变成"ccgk"，故答案为2

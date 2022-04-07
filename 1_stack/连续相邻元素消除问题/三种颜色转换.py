# R G B 三种颜色的字符串
# 两个不同的颜色相邻时，可以合并成第三种颜色
# 求最后剩下元素个数最小值

mapping = {"R": 1, "G": 2, "B": 3}


class Solution:
    def solve(self, colors):
        n = len(colors)
        if len(set(colors)) == 1:
            return n
        if n <= 1:
            return n

        curXor = 0
        for i in range(n):
            curXor ^= mapping[colors[i]]
        return 2 if curXor == 0 else 1


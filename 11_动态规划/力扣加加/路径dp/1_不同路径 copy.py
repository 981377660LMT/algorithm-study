from math import factorial as f


# 从左上角移动到右下角，总共需要移动 m+n-2 次，如果选择其中m-1次向下走
class Solution:
    def uniquePaths(self, m: int, n: int) -> int:
        return int(f(n + m - 2) / (f(m - 1) * f(n - 1)))

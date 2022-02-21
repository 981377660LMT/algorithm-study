#
# 返回这样的数之和
# @param n int整型 数的长度
# @param m int整型 各个为之和
# @return long长整型
#
class Solution:
    def sum(self, n: int, m: int) -> int:
        def dfs(index: int, cur: int, digitSum: int) -> None:
            if index == n:
                if digitSum == m:
                    self.res += cur
                return
            for nextDigit in range(10):
                dfs(index + 1, cur * 10 + nextDigit, digitSum + nextDigit)

        self.res = 0
        for start in range(1, 10):
            dfs(1, start, start)
        return self.res

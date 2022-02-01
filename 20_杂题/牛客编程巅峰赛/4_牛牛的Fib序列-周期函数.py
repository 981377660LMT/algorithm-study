# 牛牛定义f(n) = f(n-1)+f(n+1); f(1)=a, f(2)=b， 现在给定初始值 a, b，现在求第n项f(n)%1000000007的值。
# 注意fn是周期为6n的周期函数
class Solution:
    def solve(self, a, b, n):
        # write code here
        if n == 1:
            return a
        if n == 2:
            return b
        res = [b - a, -a, -b, a - b, a, b]
        n -= 2
        return res[n % 6 - 1] % 1_000_000_007


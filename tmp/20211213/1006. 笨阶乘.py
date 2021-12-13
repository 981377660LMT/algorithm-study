# 1 <= N <= 10000
# 乘除加减循环添加操作符
# 给定一个整数 N，它返回 N 的笨阶乘。
# 乘除加减
class Solution:
    def clumsy(self, n: int) -> int:
        def recur(num: int, prefix=1):
            if num <= 2:
                return prefix * num
            if num == 3:
                return prefix * 6
            return (prefix * int(num * (num - 1) / (num - 2))) + (num - 3) + recur(num - 4, -1)

        return recur(n)


print(Solution().clumsy(10))
# 输出：12
# 解释：12 = 10 * 9 / 8 + 7 - 6 * 5 / 4 + 3 - 2 * 1


# python 的整数除法是向下取整，而不是向零取整，对于负数的除法会有问题。
# python3 的地板除 "//" 是整数除法， "-3 // 2 = -2" ；

# 而 C++/Java 中的整数除法是向零取整。
# C++/Java 中 "-3 / 2 = -1" .

# py中要使用int 才是向零取整

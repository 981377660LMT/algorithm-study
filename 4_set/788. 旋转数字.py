# 我们称一个数 X 为好数, 如果它的每位数字(液晶显示)逐个地被旋转 180 度后，我们仍可以得到一个有效的，且和 X 不同的数。要求每位数字都要被旋转。
# 现在我们有一个正整数 N, 计算从 1 到 N 中有多少个数 X 是好数？

# 每位都在(2, 5, 6, 9, 0, 1, 8)内，至少一位在(2, 5, 6, 9)内
class Solution:
    def rotatedDigits(self, n: int) -> int:
        s1 = set([1, 8, 0])
        s2 = set([1, 8, 0, 6, 9, 2, 5])

        def isGood(x):
            s = set([int(i) for i in str(x)])
            return s.issubset(s2) and not s.issubset(s1)

        return sum(isGood(i) for i in range(n + 1))


print(Solution().rotatedDigits(10))

# 输出: 4
# 解释:
# 在[1, 10]中有四个好数： 2, 5, 6, 9。
# 注意 1 和 10 不是好数, 因为他们在旋转之后不变。

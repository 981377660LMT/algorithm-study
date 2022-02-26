# 三步翻转法可以实现数组循环右移操作，其时间复杂度为 ，空间复杂度为 。
class Solution:
    def solve(self, n: int, k: int) -> int:
        """
        现在有一个长度为n的字符串，进行循环右移k位的操作，
        最少对这个字符串进行几次区间反转操作能实现循环右移k位呢。
        """
        k %= n
        # 不用反转
        if n == 1 or k == 0:
            return 0
        if n == 2:
            return 1
        # 如果k为1 省略第二步
        # 如果k为2 先反转前n-1个字符再反转后n-1个字符
        if k == 1 or n - k == 1 or k == 2 or n - k == 2:
            return 2
        return 3


# 反转[0,n-1]
# 反转[0,k-1]
# 反转[k,n-1]

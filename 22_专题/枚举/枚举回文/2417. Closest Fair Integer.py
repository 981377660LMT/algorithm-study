# 思路是直接while暴力查找+剪枝


class Solution:
    def closestFair(self, num: int) -> int:
        """找到大于等于num的最小整数使得其数位中偶数个数和奇数个数相等

        1 <= num <= 1e9

        时间复杂度 `O(sqrt(n) * log(n))`
        """
        while True:
            if sum(1 if int(char) & 1 else -1 for char in str(num)) == 0:  # !奇数个数和偶数个数相等
                return num
            num += 1
            len_ = len(str(num))
            if len_ & 1:  # !num的位数为奇数时，直接变为最近的偶数位数的数
                num = 10**len_

# 670. 最大交换
# 最大交换
# 给一串数字，最多交换其中两位，求最大值
# https://leetcode.cn/problems/maximum-swap/
# 交换数位
# 给定一个非负整数，你至多可以交换一次数字中的任意两位。
# 返回你能得到的最大值。
# 给定数字的范围是 [0, 1e8]
#
# 交换相邻两位：
# 100352. 交换后字典序最小的字符串
# https://leetcode.cn/problems/lexicographically-smallest-string-after-a-swap/description/


class Solution:
    def maximumSwap(self, num: int) -> int:
        """
        O(len(num))
        贪心:把右边大的数和尽量靠左的数交换
        从右往左扫描,记录当前最大值的索引
        """

        arr = list(str(num))
        n = len(arr)
        swap1, swap2 = -1, -1
        maxIndex = n - 1
        for i in range(n - 1, -1, -1):
            if arr[i] > arr[maxIndex]:
                maxIndex = i
            elif arr[i] < arr[maxIndex]:
                swap1, swap2 = i, maxIndex
        if swap1 == -1:
            return num
        arr[swap1], arr[swap2] = arr[swap2], arr[swap1]
        return int("".join(arr))


print(Solution().maximumSwap(2736))

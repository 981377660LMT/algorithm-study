# 1864. 构成交替字符串需要的最小交换次数
# https://leetcode.cn/problems/minimum-number-of-swaps-to-make-the-binary-string-alternating/description/
# 给你一个二进制字符串 s ，现需要将其转化为一个 交替字符串 。请你计算并返回转化所需的 最小 字符交换次数，如果无法完成转化，返回 -1 。
# 交替字符串 是指：相邻字符之间不存在相等情况的字符串。例如，字符串 "010" 和 "1010" 属于交替字符串，但 "0100" 不是。
# 任意两个字符都可以进行交换，不必相邻 。


class Solution:
    def minSwaps(self, s: str) -> int:
        n = len(s)
        ones = s.count("1")
        zeros = n - ones
        if (n & 1 == 1 and abs(ones - zeros) != 1) or (n & 1 == 0 and ones != zeros):
            return -1

        # Every swap reduces the mismatch by 2.
        def calc(start: str) -> int:
            mismatch = 0
            need = int(start)
            for c in s:
                if int(c) != need:
                    mismatch += 1
                need ^= 1
            return mismatch // 2

        if ones > zeros:
            return calc("1")
        if ones < zeros:
            return calc("0")
        return min(calc("1"), calc("0"))


print(Solution().minSwaps(s="111000"))
# 输出：1
# 解释：交换位置 1 和 4："111000" -> "101010" ，字符串变为交替字符串。

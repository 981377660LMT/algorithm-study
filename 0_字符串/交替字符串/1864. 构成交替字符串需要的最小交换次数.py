class Solution:
    def minSwaps(self, s: str) -> int:
        n = len(s)
        ones = s.count("1")
        zeros = n - ones
        if (n & 1 == 1 and abs(ones - zeros) != 1) or (n & 1 == 0 and ones != zeros):
            return -1

        # Every swap reduces the mismatch by 2.
        def countDiff(start: str) -> int:
            mismatch = 0
            for char in s:
                if char != start:
                    mismatch += 1
                start = '1' if start == '0' else '0'
            return mismatch // 2

        if ones > zeros:
            return countDiff('1')
        elif ones < zeros:
            return countDiff('0')
        else:
            return min(countDiff('1'), countDiff('0'))


print(Solution().minSwaps(s="111000"))
# 输出：1
# 解释：交换位置 1 和 4："111000" -> "101010" ，字符串变为交替字符串。

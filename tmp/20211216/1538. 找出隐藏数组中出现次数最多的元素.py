# """
# This is the ArrayReader's API interface.
# You should not implement it, or speculate about its implementation
# """


class ArrayReader(object):
    # Compares 4 different elements in the array
    # return 4 if the values of the 4 elements are the same (0 or 1).
    # return 2 if three elements have a value equal to 0 and one element has value equal to 1 or vice versa.
    # return 0 : if two element have a value equal to 0 and two elements have a value equal to 1.

    # 4 2 0 2 4
    def query(self, a: int, b: int, c: int, d: int) -> int:
        ...

    # Returns the length of the array
    def length(self) -> int:
        ...


# nums 中的所有整数都为 0 或 1
# 返回 nums 中出现次数最多的值的任意索引
# 若所有的值出现次数均相同，返回 -1。


# 我们可以使用两次 API 得到任意两个位置的数是否相等
# 因此我们可以固定一个位置，然后使用 2(n - 1) 次查询分别得到位置 q = 1, 2, 3, ..., n - 1 和位置 p 的数是否相等。这样我们就还原出了整个数组。


class Solution:
    def guessMajority(self, reader: 'ArrayReader') -> int:
        n = reader.length()

        res = 0

        # ----以0,1,2为基础。
        base = reader.query(0, 1, 2, 3)
        reader3 = 1  # reader[3]的个数
        cnt2 = 0
        diff_4 = False
        for i in range(4, n):
            if reader.query(0, 1, 2, i) == base:
                reader3 += 1
            else:
                res = i
                cnt2 += 1
                if i == 4:
                    diff_4 = True

        if base == reader.query(1, 2, 3, 4) ^ diff_4:
            reader3 += 1
        else:
            cnt2 += 1

        if base == reader.query(0, 2, 3, 4) ^ diff_4:
            reader3 += 1
        else:
            cnt2 += 1

        if base == reader.query(0, 1, 3, 4) ^ diff_4:
            reader3 += 1
        else:
            cnt2 += 2

        if reader3 > n / 2:
            return 3
        else:
            if cnt2 > n / 2:
                return res
            else:
                return -1


# 输入: nums = [0,0,1,0,1,1,1,1]
# 输出: 5
# 解释: API 的调用情况如下：
# reader.length() // 返回 8，因为隐藏数组中有 8 个元素。
# reader.query(0,1,2,3) // 返回 2，查询元素 nums[0], nums[1], nums[2], nums[3] 间的比较。
# // 三个元素等于 0 且一个元素等于 1 或出现相反情况。
# reader.query(4,5,6,7) // 返回 4，因为 nums[4], nums[5], nums[6], nums[7] 有相同值。
# 我们可以推断，最常出现的值在最后 4 个元素中。
# 索引 2, 4, 6, 7 也是正确答案。


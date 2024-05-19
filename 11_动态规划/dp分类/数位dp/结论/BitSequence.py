# nums    0 | 1 | 2  3    | 4  5     6     7       | 8  9     10    11       ...
# bitnums _ | 0 | 1  0  1 | 2  0  2  1  2  0  1  2 | 3  0  3  1  3  0  1  3  ...
#
# !bitnums 第i组的元素相当于在第 0 至 i−1 组的所有元素基础上添加了 2^i 个 i。


from typing import List, Tuple


class BitSequence:
    __slots__ = ("_maxLog", "_preCount", "_preSum")

    def __init__(self, maxLog: int) -> None:
        self._maxLog = maxLog
        self._preCount = [0] * (maxLog + 1)
        self._preSum = [0] * (maxLog + 1)
        pow2 = 1
        for i in range(1, maxLog + 1):
            self._preCount[i] = self._preCount[i - 1] * 2 + pow2
            self._preSum[i] = self._preSum[i - 1] * 2 + pow2 * (i - 1)
            pow2 *= 2

    def numPreCount(self, num: int) -> int:
        """`[0, num]`中所有数的二进制表示中1的个数之和.BitSequence 前缀长度."""
        res = 0
        while num:
            i = num.bit_length() - 1
            res += self._preCount[i]
            num -= 1 << i
            res += num + 1
        return res

    def numPreSum(self, num: int) -> int:
        """`[0, num]`中所有数的二进制表示中位为1的比特位之和.BitSequence 前缀和."""
        res = 0
        while num:
            i = num.bit_length() - 1
            res += self._preSum[i]
            num -= 1 << i
            res += (num + 1) * i
        return res

    def bitIndexPreSum(self, bitIndex: int) -> int:
        """BitSequence 闭区间`[0, bitIndex]`中所有数的之和."""
        return self.bitIndexPreSumFast(bitIndex)
        if bitIndex == 0:
            return 0
        num, pos = self.bitIndexToNum(bitIndex)
        res = self.numPreSum(num - 1)
        for _ in range(pos + 1):
            lowbit = num & -num
            res += lowbit.bit_length() - 1
            num ^= lowbit
        return res

    def bitIndexPreSumFast(self, bitIndex: int) -> int:
        """BitSequence 闭区间`[0, bitIndex]`中所有数的之和."""
        if bitIndex == 0:
            return 0
        res, n, preCount, preSum = 0, 0, 0, 0
        for i in range((bitIndex + 1).bit_length() - 1, 0, -1):
            c = (preCount << i) + (i << (i - 1))
            if c <= bitIndex:
                bitIndex -= c
                res += (preSum << i) + ((i * (i - 1) // 2) << (i - 1))
                preSum += i
                preCount += 1
                n |= 1 << i
        if preCount <= bitIndex:
            bitIndex -= preCount
            res += preSum
            n += 1
        for _ in range(bitIndex):
            lowbit = n & -n
            res += lowbit.bit_length() - 1
            n ^= lowbit
        return res

    def bitIndexToNum(self, bitIndex: int) -> Tuple[int, int]:
        """BitSequence 下标为 bitIndex 对应的(数值, 二进制表示中的第几位)."""
        return self.bitIndexToNumFast(bitIndex)
        if bitIndex == 0:
            return 0, 0
        left, right = 0, 1 << self._maxLog
        while left <= right:
            mid = (left + right) // 2
            if self.numPreCount(mid) < bitIndex:
                left = mid + 1
            else:
                right = mid - 1
        return left, bitIndex - self.numPreCount(left - 1) - 1

    def bitIndexToNumFast(self, bitIndex: int) -> Tuple[int, int]:
        """BitSequence 下标为 bitIndex 对应的(数值, 二进制表示中的第几位)."""
        if bitIndex == 0:
            return 0, 0
        bitIndex -= 1
        n, preCount = 0, 0
        for i in range((bitIndex + 1).bit_length() - 1, 0, -1):
            c = (preCount << i) + (i << (i - 1))
            if c <= bitIndex:
                bitIndex -= c
                preCount += 1
                n |= 1 << i
        if preCount <= bitIndex:
            bitIndex -= preCount
            n += 1
        return n, bitIndex


if __name__ == "__main__":
    # 3145. 大数组元素的乘积
    # https://leetcode.cn/problems/find-products-of-elements-of-big-array/description/
    class Solution:
        def findProductsOfElements(self, queries: List[List[int]]) -> List[int]:
            maxQ = max(max(q) for q in queries)
            S = BitSequence(maxQ.bit_length())
            return [
                pow(2, S.bitIndexPreSum(r + 1) - S.bitIndexPreSum(l), mod) for l, r, mod in queries
            ]

    S = BitSequence(20)
    print(S.numPreSum(6))
    print(sum(v.bit_count() for v in range(7)))
    for i in range(25):
        print("i", i, S.bitIndexToNumFast(i))

    def test_bit_sequence():
        from random import randint

        def cal_sum(v: int) -> int:
            res, bit = 0, 0
            while v:
                if v & 1:
                    res += bit
                bit += 1
                v >>= 1
            return res

        S = BitSequence(20)
        for _ in range(100):
            n = randint(0, int(1e4))
            assert S.numPreCount(n) == sum(v.bit_count() for v in range(n + 1))
            assert S.numPreSum(n) == sum(cal_sum(v) for v in range(n + 1))
            assert S.bitIndexPreSum(n) == S.bitIndexPreSumFast(n)
        print("PASSED")

    test_bit_sequence()

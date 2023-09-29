# 5-所有子数组异或的和
# 二进制拆位
# https://www.luogu.com.cn/blog/endlesscheng/post-ling-cha-ba-ti-ti-mu-lie-biao
# 相当于计算01数组中有多少个子数组异或和为1


from itertools import accumulate
from typing import List


def sumOfSubarrayXor(nums: List[int]) -> int:
    def solve(bit: int) -> int:
        arr = [(v >> bit) & 1 for v in nums]
        preXor = [0] + list(accumulate(arr, lambda x, y: x ^ y))
        ones = preXor.count(1)
        count = ones * (len(preXor) - ones)
        return count * (1 << bit)

    return sum(solve(i) for i in range(40))


if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    nums = list(map(int, input().split()))
    print(sumOfSubarrayXor(nums))

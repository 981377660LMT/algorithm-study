# 给定 n 个整数（可能重复），现在需要从中挑选任意个整数，使得选出整数的异或和最大。
# 请问，这个异或和的最大可能值是多少。
# n<=1e5
# si<=2^63-1

'''
思路：
把每一个数值看作63位的01向量，把所有向量求线性基，原问题里面选择的向量的异或和在线性基里面选向量做异或和是等价的，异或要最大，就需要
优先保证高位是1，刚好就是线性基里面每一个向量都取1个，就可以构造出一个最优解，这个最优解刚好所有能取1的高位都是1，线性基里面任何一个
向量取0个或者取多于1个，都不可能构造出更好的解
'''
# https://www.acwing.com/solution/content/47831/


from functools import reduce
from typing import List


def findBase(nums: List[int]) -> List[int]:
    """求线性基"""
    n = len(nums)
    nums = nums[:]

    row = 0  # 非0行计数
    for col in range(62, -1, -1):
        # 高斯消元，找到第一个1，把下面消成0
        for r in range(row, n):
            if (nums[r] >> col) & 1:
                nums[r], nums[row] = nums[row], nums[r]
                break

        # 说明当前这一位所有向量都是0，这一位没有携带任何信息，直接跳过这一位
        if (nums[row] >> col) & 1 == 0:
            continue

        for r in range(n):
            if r != row and (nums[r] >> col) & 1:
                nums[r] ^= nums[row]

        # 新组出来一个非0行向量，非0向量的个数是基的大小，不可能超过总的向量数
        row += 1
        if row == n:
            break

    return nums[:row]


if __name__ == '__main__':
    n = int(input())
    nums = list(map(int, input().split()))
    base = findBase(nums)
    print(reduce(lambda x, y: x ^ y, base))


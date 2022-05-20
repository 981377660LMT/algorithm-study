# 给定你由 N 个整数构成的整数序列，
# 你可以从中选取一些（至少一个）进行异或（xor）运算，
# 从而得到很多不同的结果。

# 请问，所有能得到的不同的结果中第 k 小的结果是多少。
# N<=1e4
# 1<=ki<=1e18


from typing import List


def findBase(nums: List[int]) -> List[int]:
    """求线性基"""
    n = len(nums)
    nums = nums[:]

    row = 0
    for col in range(62, -1, -1):
        for r in range(row, n):
            if (nums[r] >> col) & 1:
                nums[r], nums[row] = nums[row], nums[r]
                break

        if (nums[row] >> col) & 1 == 0:
            continue

        for r in range(n):
            if r != row and (nums[r] >> col) & 1:
                nums[r] ^= nums[row]

        row += 1
        if row == n:
            break

    return nums[:row]


def main(nums: List[int], queries: List[int]) -> int:
    base = findBase(nums)

    # 如果向量数超过基的大小，一定有线性相关的向量存在，原向量已经能通过线性组合组合出0向量
    hasZero = len(nums) > len(base)
    for q in queries:
        # 如果原先就能通过异或凑出来0，就考虑求线性基取非全0系数时候能够凑出来的第k-1小的向量
        q -= int(hasZero)

        '''
        基向量的选择情况和能够组合出来的向量的排序有如下关系:
        基向量0 基向量1 基向量2 ........ 基向量k-3 基向量k-2 基向量k-1
          0      0       0                0       0        1            能组合出的第1小的向量
          0      0       0                0       1        0            能组合出的第2小的向量
          0      0       0   ......       0       1        1            能组合出的第3小的向量
          ......

        由于每一个基向量要么取1个要么取0个，所以组合出来的向量的大小刚好就是符合二进制递增的规律
        '''
        res = 0
        if q == 0:
            res = 0
        elif q > 2 ** len(base) - 1:
            res = -1
        else:
            # q 的二进制表示 异或对应的线性基
            for bit in range(len(base)):
                if (q >> bit) & 1:
                    res ^= base[~bit]
        print(res)


# N,Q<=1e4
# k<=1e18
T = int(input())
for caseId in range(T):
    n = int(input())
    nums = list(map(int, input().split()))
    q = int(input())
    queries = list(map(int, input().split()))
    print(f"Case #{caseId+1}:")
    main(nums, queries)

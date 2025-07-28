# 数组分成两部分最大化异或和
# 要求将数组nums分成非空两部分.
# 记两部分的异或和分别为x和y，要求x+y最大

# !由于x+y=(x^y)+2*(x&y)，x^y的值是固定的，所以要使x+y最大，只需要使x&y最大即可
# !考虑到x&y≥0，因此题目中B和C非空的约束可以忽略。
# !考虑如果在原序列a中，第i比特位之和为奇数，那么x&y的第i比特一定是0。
# !因此我们将这些和为奇数的比特删除也不会影响x&y的值。
# !之后仅剩下为偶数的比特位，考虑到此时一定有x=y，因此我们只需要让x最大即可。
# 则问题变为：
# 删除原数组中和为奇数的比特位，在剩下的比特位中找到一个非空子集，使得异或和最大。


from typing import List

from LinearBaseVectorSpace import VectorSpace


def solve(nums: List[int]) -> int:
    if not nums:
        return 0

    xor_ = 0
    V1 = VectorSpace()
    for v in nums:
        xor_ ^= v
        V1.add(v)

    mask = ~xor_
    V2 = VectorSpace()
    for v in V1.bases:
        V2.add(v & mask)

    res = V2.getMax()
    return res + (xor_ ^ res)


if __name__ == "__main__":
    # F - Xor Sum 3
    # https://atcoder.jp/contests/abc141/tasks/abc141_f
    n = int(input())
    nums = list(map(int, input().split()))
    print(solve(nums))
    from random import randint

    def bruteForce(nums: List[int]) -> int:
        res = 0
        n = len(nums)
        for state in range(1 << n):
            group1, group2 = [], []
            for i in range(n):
                if state & (1 << i):
                    group1.append(nums[i])
                else:
                    group2.append(nums[i])
            if len(group1) == 0 or len(group2) == 0:
                continue
            xor1, xor2 = 0, 0
            for x in group1:
                xor1 ^= x
            for x in group2:
                xor2 ^= x
            if xor1 + xor2 > res:
                res = xor1 + xor2
        return res

    for _ in range(1000):
        nums = [randint(0, 1000) for _ in range(8)]
        if solve(nums) != bruteForce(nums):
            print(nums)
            print(solve(nums))
            print(bruteForce(nums))
            raise ValueError("error")
    print("done")

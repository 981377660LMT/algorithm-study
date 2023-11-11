"""LEQ(小于等于)
给定一个数组,求有多少长度 >=2 的`子序列subSeq`满足 subSeq[0] <= subSeq[-1]
n<=3e5 ai<=1e9

类似于逆序对 但是贡献不为1而是`贡献为函数` (关于(i,j)的函数)
!对于数对(i,j) 如果nums[i]<=nums[j] 那么对答案的贡献为2^(j-i-1) 即 2^(j-1)/2^i

遍历每个位置,当成j,找[1,nums[j]]有多少i满足A[i]≤A[j],每个i的贡献为2^(j-1)/2^i,
每次找区间元素符合条件的个数,单点修改,用树状数组维护
"""

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
inv2 = (MOD + 1) // 2

powinv2 = [1]
for _ in range(int(3e5 + 10)):
    powinv2.append(powinv2[-1] * inv2 % MOD)


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, 'index must be greater than 0'
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def __repr__(self) -> str:
        preSum = []
        for i in range(self.size):
            preSum.append(self.query(i))
        return str(preSum)

    def __len__(self) -> int:
        return self.size


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))

    bit = BIT1(int(1e9 + 10))
    res = 0
    for j, num in enumerate(nums):
        preSum = bit.query(num)
        if j >= 1:
            res = (res + pow(2, j - 1, MOD) * preSum) % MOD
        bit.add(num, powinv2[j])

    print(res)

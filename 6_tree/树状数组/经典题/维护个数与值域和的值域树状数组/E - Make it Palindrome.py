# https://atcoder.jp/contests/abc290/tasks/abc290_e
# 给定一个长为n(n<=2e5)的数组nums
# 记f(sub)为将sub变为回文串所需的最小`替换`次数
# 对于每个子串sub,求f(sub)的和

# !计算相同的数中每对pair的贡献
# 每一对的贡献为min(距离左侧的距离,距离右侧的距离)
# 遍历每个组,值域树状数组维护到右侧的距离之和和个数


from collections import defaultdict
from typing import List


def makeItPalindrome(nums: List[int]) -> int:
    n = len(nums)
    mp = defaultdict(list)
    for i, char in enumerate(nums):
        mp[char].append(i)

    res = 0
    for len_ in range(1, n + 1):
        res += (len_ // 2) * (n - len_ + 1)

    for group in mp.values():
        if len(group) == 1:
            continue
        dist1 = [num + 1 for num in group]
        dist2 = [n - num for num in group]
        bitSum, bitCount = BIT1(n + 10), BIT1(n + 10)
        for i, v in enumerate(dist2):
            bitSum.add(v, v)
            bitCount.add(v, 1)
        for i in range(len(group)):
            bitSum.add(dist2[i], -dist2[i])
            bitCount.add(dist2[i], -1)
            rigthBigger = bitCount.queryRange(dist1[i], n + 1)
            rightSmallSum = bitSum.queryRange(0, dist1[i] - 1)
            res -= rigthBigger * dist1[i] + rightSmallSum  # 右侧比自己距离左侧小的距离之和 + 右侧距离大于等于自己的位置个数*左侧距离
    return res


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


# 值域树状数组
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(makeItPalindrome(nums))

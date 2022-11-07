"""最大值的期望值"""

# 给定 n 个正整数 a1∼an，对于 k=1,2,…,n，你需要：
# 两次从 a1∼ak 中均等概率地选取一个取出，记录数值并放回。
# 你的得分为两次取出数值的最大值。
# 对于每个 k，求出得分的期望值(MOD 998244353)。
# 1<=n<=2e5
# 1<=ai<=2e5

# !注意到ai范围很小,因此需要用到值域这个条件
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
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


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))

    N = int(2e5 + 10)
    countBIT = BIT1(N)  # 记录每个数出现的次数的值域树状数组
    sumBIT = BIT1(N)  # 记录前缀和的值域树状数组

    # !每次多加入一个数a,有三种情况
    # 1.两次没选到最新位置 => 和为前一次计算结果
    # 2.第一次最新位置,第二次没最新位置(或者一二反过来) => 树状数组计算
    # 3.两次都是最新位置 => 和为a
    pre = 0  # 前一次计算结果
    for i, cur in enumerate(nums, start=1):
        sum1 = pre
        sum2 = 2 * (cur * countBIT.query(cur) + sumBIT.queryRange(cur + 1, N)) % MOD
        sum3 = cur
        pre = (sum1 + sum2 + sum3) % MOD
        print(pre * pow(i * i, MOD - 2, MOD) % MOD)
        countBIT.add(cur, 1)
        sumBIT.add(cur, cur)

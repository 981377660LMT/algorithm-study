# 三阶前缀和
# n <= 2e5
# 1.修改nums[i]变为v
# 2. 查询三阶前缀和D[i]

# !计算贡献 + 公式变形
# !因为要用树状数组维护 所以数组下标从1开始算 1<=i<=n
# i处的元素 在一阶前缀和中贡献为1
# i处的元素 在二阶前缀和中贡献为n-i+1
# i处的元素 在三阶前缀和中贡献为(n-i+1)(n-i+2)/2

# 因此三阶前缀和可表示为 ∑((n-i+1)(n-i+2)/2)*xi (i=1..x)
# !∑i^2*ai ∑i*ai ∑ai 三项系数 分别为 1/2 -(2*x+3)/2 (x+1)*(x+2)/2
# bit维护前缀和即可

# 类似`巫师的力量和`


from collections import defaultdict
import os
import sys


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353
INV2 = pow(2, MOD - 2, MOD)


class BIT1:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        index += 1
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


def main() -> None:
    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    bit1, bit2, bit3 = BIT1(n + 5), BIT1(n + 5), BIT1(n + 5)
    for i in range(n):
        bit1.add(i, nums[i])
        bit2.add(i, nums[i] * i)
        bit3.add(i, nums[i] * i * i)

    for _ in range(q):
        qType, *args = map(int, input().split())
        if qType == 1:
            index, target = args[0] - 1, args[1]
            delta = target - nums[index]
            bit1.add(index, delta)
            bit2.add(index, delta * index)
            bit3.add(index, delta * index * index)
            nums[index] = target
        else:
            index = args[0] - 1
            res = bit1.query(index) * (index + 1) * (index + 2)
            res -= bit2.query(index) * (2 * index + 3)
            res += bit3.query(index)
            res %= MOD
            res *= INV2
            res %= MOD
            print(res)


if os.environ.get("USERNAME", "") == "caomeinaixi":
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()

# !数组模m后是否存在绝对众数

# !随机算法
# 给定一个长度为n(n ≤ 5000)的数组A。
# 你需要选择一个数字M(3<= M <=1e9)。
# !执行以下操作Ai = Ai %M,使得变化过后的A有一个绝对众数(出现次数严格大于n//2)。

# 思路:
# !任意选择两个数a,b 这两个数模M后`正好都是绝对众数`的概率为1/4
# 此时M一定为abs(a-b)的约数
# !sqrt(a-b)求出约数 然后对每个M检查是否存在绝对众数
# !如果重复了100次都找不到 那么真的是找不到M了
# !时间复杂度O(n*sqrt(1e9)*尝试次数)

from random import sample
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))

    def check(f: int) -> bool:
        """检查模f后是否存在绝对众数"""
        counter = dict()
        for num in nums:
            counter[num % f] = counter.get(num % f, 0) + 1
            if counter[num % f] > n // 2:
                return True
        return False

    ROUND = 100
    for _ in range(ROUND):
        a, b = sample(nums, 2)
        dist = abs(a - b)
        for f in getFactors(dist):
            if not 3 <= f <= 10**9:
                continue
            if check(f):
                print(f)
                exit(0)

    print(-1)

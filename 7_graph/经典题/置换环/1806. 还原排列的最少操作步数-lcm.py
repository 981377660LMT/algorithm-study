# 1806. 还原排列的最少操作步数-lcm

# 给你一个偶数 n​​​​​​ ，已知存在一个长度为 n 的排列 perm ，其中 perm[i] == i​（下标 从 0 开始 计数）。

# 一步操作中，你将创建一个新数组 arr ，对于每个 i ：

# 如果 i % 2 == 0 ，那么 arr[i] = perm[i / 2]
# 如果 i % 2 == 1 ，那么 arr[i] = perm[n / 2 + (i - 1) / 2]
# 然后将 arr​​ 赋值​​给 perm 。

# 要想使 perm 回到排列初始值，至少需要执行多少步操作？返回最小的 非零 操作步数。


# !解法:
# 1. simulation. O(n^2) (the cycle length <=n).
# 2. find the cycle lengths and take the lcm. O(n).
# 3. index i will transform to 2*i mod (n-1),
#    so we can use Fermat's little theorem and compute phi(n-1) by factorization. O(factor(n)).


# !所有置换环的lcm:O(n)
# https://leetcode.cn/problems/minimum-number-of-operations-to-reinitialize-a-permutation/solution/qiu-tu-zhong-suo-you-huan-de-chang-du-de-91fs/

from math import lcm


class Solution:
    def reinitializePermutation1(self, n: int) -> int:
        """O(factor(n))寻找每个元素最后映射到的位置
        每个元素最后映射到的2*i mod (n-1)的位置

        !寻找最小的k使得2^k mod (n-1) = 1 => bsgs算法
        """
        return exbsgs(2, 1, n - 1)

    def reinitializePermutation2(self, n: int) -> int:
        """O(n)寻找置换环"""
        groups = []
        visited = [False] * n
        for i in range(n):
            if visited[i]:
                continue
            visited[i] = True
            group = [i]
            while True:
                next = i // 2 if i % 2 == 0 else n // 2 + (i - 1) // 2
                if visited[next]:
                    break
                visited[next] = True
                group.append(next)
                i = next
            groups.append(group)
        return lcm(*(len(group) for group in groups))


from math import ceil, gcd, sqrt
from typing import Tuple


def bsgs(base: int, target: int, p: int) -> int:
    """Baby-step Giant-step

    在base和p互质的情况下,求解 base^x ≡ target (mod p) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(p)))

    https://dianhsu.com/2022/08/27/template-math/#bsgs
    """
    mp = dict()
    t = ceil(sqrt(p))
    target %= p
    val = 1
    for i in range(t):
        tv = target * val % p
        mp[tv] = i
        val = val * base % p

    base, val = val, 1
    if base == 0:
        return 1 if target == 0 else -1

    for i in range(t + 1):
        tv = mp.get(val, -1)
        if tv != -1 and i * t - tv > 0:  # !注意这里取等号表示允许最小解为0
            return i * t - tv
        val = val * base % p
    return -1


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """
    求a, b最大公约数,同时求出裴蜀定理中的一组系数x, y,
    满足 x*a + y*b = gcd(a, b)

    ax + by = gcd_ 返回 `(gcd_, x, y)`
    """
    if b == 0:
        return a, 1, 0
    gcd_, x, y = exgcd(b, a % b)
    return gcd_, y, x - a // b * y


def exbsgs(base: int, target: int, p: int) -> int:
    """Extended Baby-step Giant-step

    求解 base^x ≡ target (mod p) 的最小解x,
    若不存在解则返回-1

    时间复杂度: O(sqrt(p)))

    https://dianhsu.com/2022/08/27/template-math/#exbsgs
    """
    base %= p
    target %= p

    # # !平凡解
    # if target == 1 or p == 1:
    #     return 0

    cnt = 0
    d, ad = 1, 1
    while True:
        d = gcd(base, p)
        if d == 1:
            break
        if target % d:
            return -1
        cnt += 1
        target //= d
        p //= d
        ad = ad * (base // d) % p
        if ad == target:
            return cnt

    gcd_, x, _y = exgcd(ad, p)
    inv = x % p
    res = bsgs(base, target * inv % p, p)
    if res != -1:
        res += cnt
    return res


# assert Solution().reinitializePermutation(2) == 1
print(Solution().reinitializePermutation1(int(1e9)))
# print(Solution().reinitializePermutation2(int(1e5)))

from bisect import bisect_left
from typing import Callable, List


def countLIS(nums: List[int], strict=True) -> int:
    lis = []
    countPreSum = []
    for v in nums:
        target = v if strict else v + 1
        pos = search(len(lis), lambda i: lis[i][-1] >= target)
        count = 1
        if pos > 0:
            tmp1 = lis[pos - 1]
            k = search(len(tmp1), lambda k: tmp1[k] < target)
            tmp2 = countPreSum[pos - 1]
            count = tmp2[-1] - tmp2[k]
        if pos == len(lis):
            lis.append([v])
            countPreSum.append([0, count])
        else:
            lis[pos].append(v)
            tmp = countPreSum[pos]
            tmp.append(tmp[-1] + count)
    return countPreSum[-1][-1]


def search(n: int, f: Callable[[int], bool]) -> int:
    left, right = 0, n
    while left < right:
        mid = (left + right) // 2
        if f(mid):
            right = mid
        else:
            left = mid + 1
    return left


MOD = int(1e9 + 7)


def countLISWithLength(nums: List[int], length: int, strict=True) -> int:
    nums = nums[:]
    sorted_ = sorted(nums)
    for i, v in enumerate(nums):
        nums[i] = bisect_left(sorted_, v) + 2

    n = len(nums)

    def add(i: int, val: int):
        while i < n + 2:
            tree[i] = (tree[i] + val) % MOD
            i += i & -i

    def sum_(i: int) -> int:
        res = 0
        while i > 0:
            res = (res + tree[i]) % MOD
            i &= i - 1
        return res

    dp = [[0] * n for _ in range(length + 1)]
    for i in range(1, length + 1):
        tree = [0] * (n + 2)
        tmp1, tmp2 = dp[i - 1], dp[i]
        if i == 1:
            add(1, 1)
        if strict:
            for j, v in enumerate(nums):
                tmp2[j] = sum_(v - 1)
                add(v, tmp1[j])
        else:
            for j, v in enumerate(nums):
                tmp2[j] = sum_(v)
                add(v, tmp1[j])

    return sum(dp[length]) % MOD


if __name__ == "__main__":

    def check(nums: List[int], strict=True) -> int:
        res = 0
        lis = 0
        for state in range(1 << len(nums)):
            tmp = []
            for i in range(len(nums)):
                if state & (1 << i):
                    tmp.append(nums[i])
            for a, b in zip(tmp, tmp[1:]):
                if strict and a >= b:
                    break
                if not strict and a > b:
                    break
            else:
                if len(tmp) > lis:
                    lis = len(tmp)
                    res = 1
                elif len(tmp) == lis:
                    res += 1
        return res

    import random

    for _ in range(100):
        n = random.randint(1, 15)
        nums = [random.randint(-5, 5) for _ in range(n)]
        res1 = countLIS(nums)
        res2 = countLIS(nums, False)
        res3 = check(nums)
        res4 = check(nums, False)
        if res1 != res3 or res2 != res4:
            print(nums)
            print(res1, res2)
            print(res3, res4)
            break
    else:
        print("OK")

    # T = int(input())
    # for i in range(1, T + 1):
    #     n, m = map(int, input().split())
    #     nums = list(map(int, input().split()))
    #     print("Case #{}: {}".format(i, countLISWithLength(nums, m, True)))

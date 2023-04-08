# https://mugen1337.github.io/procon/tips/JOI-Candies.hpp
# https://www.ioi-jp.org/camp/2018/2018-sp-tasks/day4/candies.pdf
# https://atcoder.jp/contests/abc218/tasks/abc218_h

# 从数组不相邻选择k个数,最大化和
# n<=2e5


from heapq import heapify, heappop, heappush
from typing import List


INF = int(1e18)


def select(nums: List[int]) -> List[int]:
    """对0<=k<=ceil(n/2)的每个k,求出从数组不相邻选择k个数的最大和."""
    nums = nums[:]
    n = len(nums)
    k = (n + 1) // 2
    dp = [0] * (k + 1)
    res = 0
    left, right = [i - 1 for i in range(n)], [i + 1 for i in range(n)]
    pq = [(-nums[i], i) for i in range(n)]
    heapify(pq)
    for k_ in range(1, k + 1):
        x, i = heappop(pq)
        x = -x
        while nums[i] != x:
            x, i = heappop(pq)
            x = -x

        res += x
        dp[k_] = res
        f = left[i] >= 0 and right[i] < n
        nums[i] = -nums[i]
        if left[i] >= 0:
            nums[i] += nums[left[i]]
            nums[left[i]] = -INF
            left[i] = left[left[i]]
            if left[i] >= 0:
                right[left[i]] = i
        if right[i] < n:
            nums[i] += nums[right[i]]
            nums[right[i]] = -INF
            right[i] = right[right[i]]
            if right[i] < n:
                left[right[i]] = i
        if f:
            heappush(pq, (-nums[i], i))
        else:
            nums[i] = -INF

    return dp


def rob(nums: List[int]) -> int:
    """打家劫舍."""
    res = select(nums)
    return max(res)


if __name__ == "__main__":
    nums, expcted = [3, 5, 1, 7, 6], [0, 7, 12, 10]
    assert select(nums) == expcted

    nums = [
        623239331,
        125587558,
        908010226,
        866053126,
        389255266,
        859393857,
        596640443,
        60521559,
        11284043,
        930138174,
        936349374,
        810093502,
        521142682,
        918991183,
        743833745,
        739411636,
        276010057,
        577098544,
        551216812,
        816623724,
    ]
    expected = [
        0,
        936349374,
        1855340557,
        2763350783,
        3622744640,
        4439368364,
        5243250666,
        5982662302,
        6605901633,
        7183000177,
        7309502029,
    ]
    assert select(nums) == expected

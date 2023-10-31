from collections import Counter, deque
from typing import List

MOD = int(1e9 + 7)


def boundedKnapsackDp(
    values: List[int], weights: List[int], counts: List[int], maxCapacity: int
) -> int:
    """
    单调队列优化多重背包问题.
    选择物品使得价值最大,总重量不超过maxCapacity.求出最大价值.
    时间复杂度O(n*maxCapacity).
    """
    dp = [0] * (maxCapacity + 1)
    for i, count in enumerate(counts):
        v, w = values[i], weights[i]
        for rem in range(w):
            j = 0
            queue = deque()
            while True:
                pos = j * w + rem
                if pos > maxCapacity:
                    break
                cand = dp[pos] - j * v
                while queue and queue[-1][0] <= cand:
                    queue.pop()
                queue.append((cand, j))
                dp[pos] = queue[0][0] + j * v
                if j - queue[0][1] == count:
                    queue.popleft()
                j += 1
    return dp[maxCapacity]


def boundedKnapsackDpBinary(
    values: List[int], weights: List[int], counts: List[int], maxCapacity: int
) -> int:
    """多重背包二进制优化."""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    dp = [0] * (maxCapacity + 1)
    for i, count in enumerate(counts):
        v, w = values[i], weights[i]
        remain = count
        step = 1
        while remain > 0:
            k = step if step < remain else remain
            for j in range(maxCapacity, k * w - 1, -1):
                dp[j] = max(dp[j], dp[j - k * w] + k * v)
            remain -= k
            step *= 2
    return dp[maxCapacity]


def boundedKnapsackDpNaive(
    values: List[int], weights: List[int], counts: List[int], maxCapacity: int
) -> int:
    """多重背包朴素解法."""

    def max(a: int, b: int) -> int:
        return a if a > b else b

    n = len(counts)
    dp = [[0] * (maxCapacity + 1) for _ in range(n + 1)]
    for i, count in enumerate(counts):
        v, w = values[i], weights[i]
        tmp1, tmp2 = dp[i + 1], dp[i]
        for j in range(maxCapacity + 1):
            # 枚举选了 k=0,1,2,...num 个第 i 种物品
            k = 0
            while k <= count and k * w <= j:
                tmp1[j] = max(tmp1[j], tmp2[j - k * w] + k * v)
                k += 1
    return dp[n][maxCapacity]


def boundedKnapsackDPCountWays(values: List[int], counts: List[int]) -> List[int]:
    """
    多重背包求方案数(分组前缀和优化).
    dp[i] 表示总价值为 i 的方案数.
    O(n*sum(values[i]*counts[i]))
    """
    n = len(values)
    allSum = 0
    count0 = 0
    for i in range(n):
        count = counts[i]
        value = values[i]
        if value == 0:
            count0 += count
            continue
        allSum += count * value
    dp = [0] * (allSum + 1)
    dp[0] = count0 + 1

    maxJ = 0
    for i in range(n):
        value = values[i]
        if value == 0:
            continue
        count = counts[i]
        maxJ += value * count
        for j in range(value, maxJ + 1):
            dp[j] = (dp[j] + dp[j - value]) % MOD
        for j in range(maxJ, value * (count + 1) - 1, -1):
            dp[j] = (dp[j] - dp[j - value * (count + 1)]) % MOD
    return dp


def boundedKnapsackDpCountWaysWithUpper(
    values: List[int], counts: List[int], upper: int
) -> List[int]:
    """O(n*upper)."""
    n = len(values)
    count0 = 0
    for i in range(n):
        count = counts[i]
        value = values[i]
        if value == 0:
            count0 += count
            continue
    dp = [0] * (upper + 1)
    dp[0] = count0 + 1

    maxJ = 0
    for i in range(n):
        value = values[i]
        if value == 0:
            continue
        count = counts[i]
        maxJ += value * count
        if maxJ > upper:
            maxJ = upper
        for j in range(value, maxJ + 1):
            dp[j] = (dp[j] + dp[j - value]) % MOD
        for j in range(maxJ, value * (count + 1) - 1, -1):
            dp[j] = (dp[j] - dp[j - value * (count + 1)]) % MOD
    return dp


if __name__ == "__main__":
    import sys

    input = sys.stdin.readline

    # n, maxCapacity = map(int, input().split())
    # values = [0] * n
    # weights = [0] * n
    # counts = [0] * n
    # for i in range(n):
    #     values[i], weights[i], counts[i] = map(int, input().split())
    # print(boundedKnapsackDp(values, weights, counts, maxCapacity))

    n, target = map(int, input().split())
    counts = list(map(int, input().split()))
    dp = boundedKnapsackDPCountWays([1] * n, counts)
    print(dp[target])

    # 2585. 获得分数的方法数
    # https://leetcode.cn/problems/number-of-ways-to-earn-points/description/
    # 考试中有 n 种类型的题目。给你一个整数 target 和一个下标从 0 开始的二维整数数组 types ，
    # 其中 types[i] = [counti, marksi] 表示第 i 种类型的题目有 counti 道，每道题目对应 marksi 分。
    # !返回你在考试中恰好得到 target 分的方法数。由于答案可能很大，结果需要对 1e9 +7 取余。
    # !注意，同类型题目无法区分。
    # target<=1000
    # n<=50
    # counti<=50
    # !O(n*target) 按模分组前缀和优化dp
    # dp[i][j]表示前i种题目恰好得到j分的方法数
    # !ndp[j] = sum(dp[j-k*mark] for k in range(count+1) if j-k*mark>=0
    # 这是一个按模分组的前缀和
    class Solution:
        def waysToReachTarget(self, target: int, types: List[List[int]]) -> int:
            values = [v[0] for v in types]
            counts = [v[1] for v in types]
            return boundedKnapsackDpCountWaysWithUpper(values, counts, target)[target]

    # 2902. 和带限制的子多重集合的数目
    # https://leetcode.cn/problems/count-of-sub-multisets-with-bounded-sum/description/
    # 给你一个下标从 0 开始的非负整数数组 nums 和两个整数 l 和 r 。
    # 请你返回 nums 中子多重集合的和在闭区间 [l, r] 之间的 子多重集合的数目 。
    # 由于答案可能很大，请你将答案对 109 + 7 取余后返回。
    # 子多重集合 指的是从数组中选出一些元素构成的 无序 集合，每个元素 x 出现的次数可以是 0, 1, ..., occ[x] 次，其中 occ[x] 是元素 x 在数组中的出现次数。
    class Solution2:
        def countSubMultisets(self, nums: List[int], l: int, r: int) -> int:
            counter = Counter(nums)
            values = list(counter)
            counts = list(counter.values())
            dp = boundedKnapsackDpCountWaysWithUpper(values, counts, r)
            return sum(dp[l : r + 1]) % MOD

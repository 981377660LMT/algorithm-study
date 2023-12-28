"""
1425. 带限制的子序列和-相邻项不超过k的最大子序列和
https://leetcode.cn/problems/constrained-subsequence-sum/

# !请你返回 非空 子序列元素和的最大值，子序列需要满足：
# !子序列中每两个 相邻 的整数 nums[i] 和 nums[j] 它们在原数组中的下标 i 和 j 满足 i < j 且 j - i <= k 。
"""

from MonoQueue import MonoQueue

from typing import List, Tuple

INF = int(1e18)


def constrainedSubsetSum(nums: List[int], k: int) -> int:
    """
    - dp[i] 表示前 i 个元素中，以第 i 个元素结尾的子序列元素和的最大值(0<=i<n)
    - dp[i] = max(dp[i], max(dp[i - k] ,..., dp[i-1], 0) + nums[i])
    - res = max(dp)
    """
    n = len(nums)
    queue = MonoQueue[Tuple[int, int]](lambda x, y: x[0] > y[0])  # (dp[i], i)
    dp = [-INF] * n
    for i, num in enumerate(nums):
        while queue and i - queue.head()[1] > k:  # 1.不在窗口内的元素出队
            queue.popleft()
        preMax = max(0, queue.head()[0]) if queue else 0
        dp[i] = max(dp[i], preMax + num)  # 2.更新dp
        queue.append((dp[i], i))  # 3.入队
    return max(dp)


if __name__ == "__main__":
    # nums = [10,-2,-10,-5,20], k = 2
    print(constrainedSubsetSum(nums=[10, -2, -10, -5, 20], k=2))

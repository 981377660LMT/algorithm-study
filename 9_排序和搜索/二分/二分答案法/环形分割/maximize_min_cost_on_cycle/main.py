from typing import Callable


def maximize_min_cost_on_cycle_dp(
    n: int, cost: Callable[[int, int], int], k: int, lower: int, upper: int
) -> int:
    """
    给定一个n个点的环形数组, [start, end) 的代价为 cost(start, end), 且 cost 满足单调性.
    将环形数组分成k个非空连续子数组, 最大化这k个子数组的代价的最小值.
    返回这个最大的最小值.
    O(n*log(upper-lower))时间复杂度.

    Args:
        n (int): 环形数组的点数
        cost (Callable[[int, int], int]): 代价函数
        k (int): 分割的段数
        lower (int): 代价的下界
        upper (int): 代价的上界

    Returns:
        int: 最大的最小值
    """
    if k > n:
        raise ValueError("k must be not greater than n")

    def check(mid: int) -> bool:
        # 先求解链上的问题，剪枝
        count = 0
        left = 0
        for right in range(n):
            if cost(left, right + 1) >= mid:
                count += 1
                left = right + 1
        if count >= k:
            return True
        if count <= k - 2:
            return False

        # 预处理 next 数组：对每个 left 找到最小的 right 满足 cost(left, right) >= mid
        # [left, right) 为满足条件的左闭右开区间.
        next = [-1] * n
        right = 0
        for left in range(n):
            while right < n and cost(left, right) < mid:
                right += 1
            next[left] = right if cost(left, right) >= mid else -1

        # dp[i] = (count, next)
        # 意味着从位置 i 出发，最多能获得 count 段，且最终结束位置为 next
        dp = [(0, i) for i in range(n + 1)]
        for i in range(n - 1, -1, -1):
            if next[i] != -1:
                cnt, nxt = dp[next[i]]
                dp[i] = (cnt + 1, nxt)

        # 尝试在环上拆分：选取一个断点 i，把环拆成链，额外检查首尾部分能否合并
        for i in range(n):
            cnt, end = dp[i]
            if cnt <= k - 2:
                break
            # 检查最后一段是否满足条件，把环的前后部分合并
            cnt += cost(0, i) + cost(end, n) >= mid
            if cnt >= k:
                return True
        return False

    left, right = lower, upper
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    return right


if __name__ == "__main__":
    from typing import List

    # 3464. 正方形上的点之间的最大距离
    # https://leetcode.cn/problems/maximize-the-distance-between-points-on-a-square/description/
    class Solution:
        def maxDistance(self, side: int, points: List[List[int]], k: int) -> int:
            def trans(x: int, y: int) -> int:
                if x == 0:
                    return y
                if y == side:
                    return side + x
                if x == side:
                    return 3 * side - y
                return 4 * side - x

            n = len(points)
            nums = [trans(x, y) for x, y in points]
            nums.sort()
            nums.append(4 * side + nums[0])
            res = maximize_min_cost_on_cycle_dp(n, lambda i, j: nums[j] - nums[i], k, 1, 2 * side)
            return res

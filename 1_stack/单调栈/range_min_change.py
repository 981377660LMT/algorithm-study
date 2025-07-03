from typing import List, Tuple

INF = int(1e18)


def range_min_change(nums: List[int], default_val=INF) -> List[List[Tuple[int, int, int, int]]]:
    """
    维护区间最小值的变化历史。
    返回：res[i]，表示右端点r=i+1时，所有受影响区间[l,r)的最小值变化记录：(l, r, old_min, new_min)
    每次右端点推进，所有被当前元素“刷新”最小值的区间都会被记录下来，适用于区间DP、单调栈优化等场景。
    """
    n = len(nums)
    res = [[] for _ in range(n)]
    stack = []
    for i, v in enumerate(nums):
        res[i].append((i, i + 1, default_val, v))
        ptr = i
        while stack:
            l, r, x = stack[-1]
            if x <= v:
                break
            res[i].append((l, r, x, v))
            ptr = l
            stack.pop()
        stack.append((ptr, i + 1, v))
        res[i].reverse()
    return res


if __name__ == "__main__":
    # Example usage
    # Changes for index 0: [(0, 1, 1000000000000000000, 5)]
    # Changes for index 1: [(0, 1, 5, 3), (1, 2, 1000000000000000000, 3)]
    # Changes for index 2: [(2, 3, 1000000000000000000, 4)]
    # Changes for index 3: [(0, 2, 3, 2), (2, 3, 4, 2), (3, 4, 1000000000000000000, 2)]
    # Changes for index 4: [(4, 5, 1000000000000000000, 6)]
    A = [5, 3, 4, 2, 6]
    result = range_min_change(A)
    for i, changes in enumerate(result):
        print(f"Changes for index {i}: {changes}")

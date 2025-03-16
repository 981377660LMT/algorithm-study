from bisect import bisect_left
from typing import List, Sequence, TypeVar


T = TypeVar("T")


def longest_common_subsequence(seq1: Sequence[T], seq2: Sequence[T]) -> List[T]:
    """
    求两个无重复元素序列的最长公共子序列，时间复杂度 O(n log n).
    """

    mp = {v: i for i, v in enumerate(seq2)}

    # 在 seq1 中筛选出也在 seq2 中出现的元素，同时将其映射为 seq2 中的下标
    filtered = []
    mapped = []  # mapped[i] 对应 filtered[i] 在 seq2 中的下标
    for val in seq1:
        if val in mp:
            filtered.append(val)
            mapped.append(mp[val])

    # 如果没有公共元素，则返回空列表
    if not mapped:
        return []

    # 使用二分查找方法求 mapped 的最长递增子序列（LIS）
    # dp 存储的是构造出的最长递增子序列中各长度末尾的 mapped 值
    dp = []  # 存储递增序列的数值（mapped 值）
    dp_idx = []  # 存储 dp 数组中对应的 mapped 序列在 filtered/mapped 数组中的索引
    prev = [-1] * len(mapped)  # prev[i] 用来记录 mapped[i] 的前驱索引，用于恢复路径

    for i, val in enumerate(mapped):
        # 找到 dp 中第一个不小于 val 的位置
        pos = bisect_left(dp, val)
        if pos == len(dp):
            dp.append(val)
            dp_idx.append(i)
        else:
            dp[pos] = val
            dp_idx[pos] = i

        # 记录前驱节点索引：如果 pos > 0，则前驱为 dp_idx[pos-1]
        if pos > 0:
            prev[i] = dp_idx[pos - 1]

    # dp_idx[-1] 存储最长递增序列末尾在 mapped 中的索引
    lis_length = len(dp)
    lis = [0] * lis_length
    k = dp_idx[-1]
    for j in range(lis_length - 1, -1, -1):
        lis[j] = filtered[k]
        k = prev[k]

    return lis  # type: ignore


# 示例：
if __name__ == "__main__":
    # 假设两个序列无重复元素
    seq1 = [3, 9, 8, 4, 5, 7, 6]
    seq2 = [3, 4, 9, 5, 7, 8, 6]
    result = longest_common_subsequence(seq1, seq2)
    print("最长公共子序列:", result)

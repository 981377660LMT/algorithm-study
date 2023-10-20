# https://atcoder.jp/contests/arc147/tasks/arc147_c


from typing import List, Tuple


def minDiffSum(intervals: List[Tuple[int, int]]) -> int:
    """
    数轴上有n个人,每个人可以呆的区间为[l, r].
    求所有人的距离和的最小值。
    2<=n<=3e5
    1<=l<=r<=1e9

    将区间排序, R[i]越小的越靠左排, L[i]越大的越靠右排。
    """
    n = len(intervals)
    left = [l for l, _ in intervals]
    left.sort(reverse=True)
    right = [r for _, r in intervals]
    right.sort()
    res = 0
    for i in range(n):
        if right[i] > left[i]:
            break
        res += (n - 2 * i - 1) * (left[i] - right[i])
    return res


if __name__ == "__main__":
    n = int(input())
    intervals = []
    for _ in range(n):
        l, r = map(int, input().split())
        intervals.append((l, r))
    print(minDiffSum(intervals))

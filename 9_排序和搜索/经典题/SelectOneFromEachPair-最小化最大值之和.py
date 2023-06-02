# SelectOneFromEachPair-最小化最大值之和
# https://zhuanlan.zhihu.com/p/268630329

# 给定n个对(ai,bi),对每个对,需要选择ai加入集合A,或者选择bi加入集合B
# !最小化`max(A)+max(B)`


from typing import List, Tuple


def selectOneFromEachPairMinimizeMaxSum(pairs: List[Tuple[int, int]]) -> int:
    """
    给定n个对(ai,bi),对每个对,需要选择ai加入集合A,或者选择bi加入集合B.
    !最小化`max(A)+max(B)`的值.
    """
    if not pairs:
        return 0
    pairs.sort()
    sufMax = [b for _, b in pairs] + [0]
    for i in range(len(sufMax) - 2, -1, -1):
        if sufMax[i] < sufMax[i + 1]:
            sufMax[i] = sufMax[i + 1]
    res = sufMax[0]
    for i in range(len(pairs)):
        res = min(res, pairs[i][0] + sufMax[i + 1])
    return res


if __name__ == "__main__":
    assert (
        selectOneFromEachPairMinimizeMaxSum([(1, 4), (1, 8), (2, 3), (2, 7), (3, 1), (4, 2)])
    ) == 4

    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    # CF1408D
    n, m = map(int, input().split())
    robber = [tuple(map(int, input().split())) for _ in range(n)]
    searchLight = [tuple(map(int, input().split())) for _ in range(m)]

    pairs = []  # 每个强盗摆脱每个探照灯的(x方向需要移动的距离, y方向需要移动的距离)
    for rx, ry in robber:
        for sx, sy in searchLight:
            if rx <= sx and ry <= sy:
                pairs.append((sx - rx + 1, sy - ry + 1))
    print(selectOneFromEachPairMinimizeMaxSum(pairs))

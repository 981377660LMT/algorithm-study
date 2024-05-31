from bisect import bisect_left, bisect_right
from typing import List, Sequence, Tuple


INF = int(2e18)


def LISOfSequnce(seq: List[int], strict=True) -> List[int]:
    """每个元素是否在LIS中.
    1: 不在LIS中
    2: 在某些LIS中
    3: 在所有LIS中
    """
    seq = seq[:]
    lis, dp1 = lisDp(seq, strict)
    seq = seq[::-1]
    for i in range(len(seq)):
        seq[i] = -seq[i]
    _, dp2 = lisDp(seq, strict)
    dp2 = dp2[::-1]

    n = len(seq)
    counter = [0] * n
    for i in range(n):
        if dp1[i] + dp2[i] == lis:
            counter[dp1[i]] += 1
    res = [0] * n
    for i in range(n):
        if dp1[i] + dp2[i] < lis:
            res[i] = 1
        elif counter[dp1[i]] == 1:
            res[i] = 3
        else:
            res[i] = 2
    return res


def lisDp(seq: Sequence[int], strict=True) -> Tuple[int, Sequence[int]]:
    """返回每个位置为结尾的LIS长度(不包括自身)."""
    n = len(seq)
    dp = [INF] * n
    lis = 0
    lisRank = [0] * n
    f = bisect_left if strict else bisect_right
    for i, v in enumerate(seq):
        pos = f(dp, v)
        dp[pos] = v
        lisRank[i] = pos
        if lis < pos:
            lis = pos
    return lis, lisRank


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    res = LISOfSequnce(nums)
    print("".join(map(str, res)))

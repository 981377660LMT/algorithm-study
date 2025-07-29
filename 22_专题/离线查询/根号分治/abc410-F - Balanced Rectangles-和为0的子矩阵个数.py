# abc410-F - Balanced Rectangles-二维矩阵
# 给定 T 个测试用例，每个用例是一个 H×W 的网格，由字符 # 和 . 构成。
# 我们要统计子矩形区域 (u,d,l,r)（1≤u≤d≤H，1≤l≤r≤W），使得该区域内 # 与 . 的格子数相等。
#
# !所有数据 H*W<=3e5
#
# 令 # 记作 +1，. 记作 -1。一个子矩形和为 0 当且仅当 # 与 . 数目相等。
# !我们选择枚举“较小维度”的区间对，转化成一维“和为 0 的子数组”计数：
# !O(min(H,W)^2 * max(H,W))
#
# !和为零的子矩阵


from typing import List


def solve(h: int, w: int, grid: List[List[int]]) -> int:
    if h > w:
        h, w = w, h
        grid = list(zip(*grid))  # type: ignore

    res = 0
    for u in range(h):
        colSum = [0] * w
        for d in range(u, h):
            for j, v in enumerate(grid[d]):
                colSum[j] += v
            presum = dict({0: 1})
            cursum = 0
            for v in colSum:
                cursum += v
                res += presum.get(cursum, 0)
                presum[cursum] = presum.get(cursum, 0) + 1
    return res


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        H, W = map(int, input().split())
        grid = [list(map(lambda x: 1 if x == "#" else -1, input().strip())) for _ in range(H)]
        print(solve(H, W, grid))

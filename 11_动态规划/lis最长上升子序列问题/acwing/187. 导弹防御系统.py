# 1≤n≤50
# dfs爆搜 求最小深度 要记录全局最小值
# 为什么不用bfs?:空间会爆炸
up = [0] * 55
down = [0] * 55

# n2n
def dfs(index: int, upSeq: int, downSeq: int) -> None:
    """当前元素，上升子序列个数，下降子序列个数"""
    global res, up, down, n
    if upSeq + downSeq >= res:
        return
    if index == n:
        res = min(res, upSeq + downSeq)
        return
    # 当前数放到上升子序列

    # 当前数放到下降子序列


while True:
    n = int(input())
    if n == 0:
        break
    nums = list(map(int, input().split()))
    res = n
    dfs(0, 0, 0)
    print(res)

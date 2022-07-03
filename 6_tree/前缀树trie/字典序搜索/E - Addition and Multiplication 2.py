# 字典序搜索 十叉树
# 最初为0 每次可以选择一个分支 即 x => 10*x + i
# 求最终的最大值
# 边权之和<=1e6

# !答案会超出longlong 输出字符串
# !1. 不能背包dp 因为 字符串太长了
# !2. 贪心
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    k = int(input())
    costs = [int(1e18)] + list(map(int, input().split()))

    # !位数要多 多出来的用来把大的数位放在前面
    min_ = min(costs)
    dLen = k // min_
    remain = k - dLen * min_
    res: list[str] = []
    while len(res) < dLen:
        for select in range(9, 0, -1):
            if remain >= costs[select] - min_:
                res.append(str(select))
                remain -= costs[select] - min_
                break

    print("".join(res))


while True:
    try:
        main()
    except (EOFError, ValueError):
        break

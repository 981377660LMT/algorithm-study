# 买卖股票 现金无限
# 已知n天股票的价格
# 每天可以 买/卖/不做
# 求最大收益
# n<=2e5
# 1<=pi<=1e9

# !O(n^2)的解法
# !dp[i][j] 表示第i天，有j股股票时的最大收益

# !dp过不了 考虑贪心
# !想象现在有全是右括号的非法的序列 (-1 0 1)
# !现在要以最小代价改变变为合法的序列(变为左括号/没有括号 每个位置都有两次改变机会)
# !求最小的代价 满足每个前缀和>=0(变为合法的括号序列)
# 优先队列


from heapq import heappop, heappush
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    nums = list(map(int, input().split()))
    res, pq = 0, []
    for num in nums:
        res += num  # 卖出
        heappush(pq, num)  # 不操作
        heappush(pq, num)  # 买入
        res -= heappop(pq)
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":

        while True:
            main()
    else:
        main()

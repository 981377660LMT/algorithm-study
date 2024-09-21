# abc360-E - Random Swaps of Balls-等于不等等于两种状态
# https://atcoder.jp/contests/abc360/tasks/abc360_e
# 随机交换
#
# n个球，其中第一个是黑球，其余是白球。进行以下操作k次：
# 随机两次独立选取i,j，交换第i个球和第j个球。问黑球位置的期望值。
# 假设最终黑球位于第i个球的概率是pi，根据期望定义，答案就是Σi*pi。
# !容易发现第2,3,...,n个球没有本质区别，其概率都是一样的，因此最终我们需要求的就两个数
# !p1: 黑球位于第1个球的概率
# !p2: 黑球位于非第1个球的概率
# 容易发现经过第k次操作的概率，仅依赖于第k−1次的情况，因此可以迭代球
# dp[k][0/1]表示经过k次操作，黑球位于第1个球/非第1个球的概率

MOD = 998244353


def randomSwapsOfBalls(n: int, k: int) -> int:
    ...


if __name__ == "__main__":
    n, k = map(int, input().split())
    print(randomSwapsOfBalls(n, k))

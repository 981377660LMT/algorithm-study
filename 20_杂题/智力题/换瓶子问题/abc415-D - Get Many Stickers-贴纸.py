# https://atcoder.jp/contests/abc415/tasks/abc415_d
#
# 有一个神奇的可乐店，只能用空瓶换新可乐和纪念贴纸。
# 初始有 N 瓶可乐，0 空瓶。每次可以：
#
# 喝掉一瓶可乐，空瓶+1，可乐-1。
# 选择某个 i（1≤i≤M），用 Ai 个空瓶换 Bi 瓶可乐和 1 枚贴纸（需空瓶≥Ai，且 Bi<Ai）。 问最多能获得多少枚贴纸。
#
# !每次优先用“消耗空瓶最少”的方案（即 Ai-Bi 最小的可行方案），这样能最大化贴纸数。


def solve():
    N, M = map(int, input().split())
    A, B = [0] * M, [0] * M
    for i in range(M):
        A[i], B[i] = map(int, input().split())

    data = [(a - b, a, b) for a, b in zip(A, B)]
    data.sort()

    res = 0
    remain = N
    for d, a, _ in data:
        if remain < a:
            continue
        count = 1 + (remain - a) // d
        remain -= count * d
        res += count

    print(res)


if __name__ == "__main__":
    solve()

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class FenwickTree:
    def __init__(self, size):
        self.N = size + 2
        self.tree = [0] * self.N

    def update(self, index, value):
        index += 1
        while index < self.N:
            self.tree[index] += value
            index += index & -index

    def query(self, index):
        index += 1
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res


if __name__ == "__main__":
    N, M = map(int, sys.stdin.readline().split())
    A = list(map(int, sys.stdin.readline().split()))

    prefix = [0] * (N + 1)
    for i in range(N):
        prefix[i + 1] = (prefix[i] + A[i]) % M

    ft_count = FenwickTree(M)
    ft_sum = FenwickTree(M)

    ft_count.update(0, 1)
    ft_sum.update(0, 0)

    sum_y = 0
    total = 0

    for r in range(1, N + 1):
        x_r = prefix[r]
        cnt_less_eq_x = ft_count.query(x_r)
        sum_less_eq_x = ft_sum.query(x_r)
        total += x_r * r - sum_y + M * (r - cnt_less_eq_x)
        ft_count.update(x_r, 1)
        ft_sum.update(x_r, x_r)
        sum_y += x_r

    print(total)

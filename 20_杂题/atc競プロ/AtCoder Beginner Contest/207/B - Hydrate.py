import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    a, b, c, d = map(int, input().split())

    def check(mid: int) -> bool:
        return (a + mid * b) <= c * mid * d

    left, right = 0, int(1e18)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1

    print(left if check(left) else -1)  # !如果不能,则输出-1 : 只需最后验证一下left即可

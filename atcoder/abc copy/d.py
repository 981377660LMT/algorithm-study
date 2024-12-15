import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    N, S = map(int, input().split())
    A = list(map(int, input().split()))
    cycleSum = sum(A)
    if cycleSum == S:
        print("Yes")
        exit(0)
    if min(A) > S:
        print("No")
        exit(0)

    def check(target):
        arr = A + A
        l = 0
        cur = 0
        for r in range(len(arr)):
            cur += arr[r]
            while cur > target:
                cur -= arr[l]
                l += 1
            if cur == target:
                return True
        return False

    if cycleSum > S:
        print("Yes" if check(S) else "No")
    else:
        k = S // cycleSum
        remainder = S - k * cycleSum
        if remainder == 0:
            print("Yes")
        else:
            print("Yes" if check(remainder) else "No")

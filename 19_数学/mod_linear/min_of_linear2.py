# 一次函数取模最小值
# n,mod<=1e9 0<=a,b<mod
def min_of_linear(n: int, mod: int, a: int, b: int, is_min=True, p=1, q=1) -> int:
    """
    ```
    min((a*x+b)%mod for x in range(n))
    ```
    """
    if a == 0:
        return b
    if is_min:
        if b >= a:
            t = (mod - b + a - 1) // a
            c = (t - 1) * p + q
            if n <= c:
                return b
            n -= c
            b += a * t - mod
        b = a - 1 - b
    else:
        if b < mod - a:
            t = (mod - b - 1) // a
            c = t * p
            if n <= c:
                return a * ((n - 1) // p) + b
            n -= c
            b += a * t
        b = mod - 1 - b
    d = mod // a
    c = min_of_linear(n, a, mod % a, b, not is_min, (d - 1) * p + q, d * p + q)
    return a - 1 - c if is_min else mod - 1 - c


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    T = int(input())
    for _ in range(T):
        n, mod, a, b = map(int, input().split())
        print(min_of_linear(n, mod, a, b))

# 求 a 乘 b 对 p 取模的值。
# a,b,p<=int(1e18)


def main():
    a = int(input())
    n = int(input())
    MOD = int(input())

    res = 0
    while n:
        if n & 1:
            res += a
            res %= MOD

        a **= 2
        a %= MOD
        n >>= 1

    print(res)


if __name__ == "__main__":
    main()

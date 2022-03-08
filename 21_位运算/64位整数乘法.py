# 求 a 乘 b 对 p 取模的值。
# 1≤a,b,p≤1018
def main():
    a = int(input())
    b = int(input())
    MOD = int(input())

    res = 0
    while b:
        if b & 1:
            res += a
            res %= MOD
        a *= 2
        a %= MOD
        b >>= 1
    print(res)


if __name__ == "__main__":
    main()

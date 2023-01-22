"""
多重背包
1. 二进制加速
2. bit加速 (int 来代替 set)
"""

# 高桥君有n种硬币，每种ai元的硬币有bi个，现在要用这些硬币凑成x元，问是否可以做到？
# n<=50 ai<=100 bi<=50 x<=1e4

if __name__ == "__main__":
    n, x = map(int, input().split())
    coins = [list(map(int, input().split())) for _ in range(n)]  # (cost, count)

    newCoins = []  # 每个硬币的面额(cost)
    for cost, count in coins:
        cur = 1
        while cur <= count:
            newCoins.append(cost * cur)
            count -= cur
            cur *= 2
        if count:
            newCoins.append(cost * count)

    dp = 1 << x
    for cost in newCoins:
        dp |= dp >> cost
    print("Yes" if dp & 1 else "No")

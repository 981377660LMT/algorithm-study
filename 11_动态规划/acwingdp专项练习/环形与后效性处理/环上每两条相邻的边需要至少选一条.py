# 环上每两条相邻的边需要至少选一条
# 有多少种选法

MOD = int(1e9 + 7)


def calCycle(n: int) -> int:
    """
    环上每两条相邻的边需要`至少`选一条 有多少种选法

    n表示环上的点数
    """
    if n == 1:
        return 1

    # 不选第一个
    dp1 = [[0, 0] for _ in range(n)]
    dp1[0][0] = 1
    for i in range(1, n):
        dp1[i][0] = dp1[i - 1][1]
        dp1[i][1] = (dp1[i - 1][0] + dp1[i - 1][1]) % MOD

    # 选第一个
    dp2 = [[0, 0] for _ in range(n)]
    dp2[0][1] = 1
    for i in range(1, n):
        dp2[i][0] = dp2[i - 1][1]
        dp2[i][1] = (dp2[i - 1][0] + dp2[i - 1][1]) % MOD

    return (dp1[n - 1][1] + dp2[n - 1][0] + dp2[n - 1][1]) % MOD


if __name__ == "__main__":
    print(calCycle(1))
    print(calCycle(2))
    print(calCycle(3))

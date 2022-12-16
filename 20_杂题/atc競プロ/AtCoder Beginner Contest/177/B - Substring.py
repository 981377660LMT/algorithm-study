# 成为子串的最少修改次数 - 定长滑窗
# n<=1000


def minModify(long: str, short: str) -> int:
    n, m = len(long), len(short)
    res = m
    for start in range(n - m + 1):
        diff = 0
        for i in range(m):
            if long[start + i] != short[i]:
                diff += 1
        res = min(res, diff)
    return res


if __name__ == "__main__":
    s = input()
    t = input()
    print(minModify(s, t))

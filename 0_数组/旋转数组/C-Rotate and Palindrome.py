import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)

# 给定一个长度为n的字符串s。
# 有两种操作可以对s进行：
# 1. 支付a元,将s的左端字符串移动到右端。即将s1s2…sN变为s2…sNs1。 (轮转)
# 2. 支付b元,选择1≤i≤n的整数i,将si替换为任意小写英文字母。
# 问:将s变为回文串至少需要多少元?
if __name__ == "__main__":
    n, a, b = map(int, input().split())
    s = input()

    s += s  # 轮转字符串加倍
    res = INF
    for i in range(n):  # 枚举轮转次数
        cur = i * a
        for j in range(n // 2):
            if s[i + j] != s[i + n - 1 - j]:
                cur += b
        res = min(res, cur)

    print(res)

# 京东-相邻都不相等的最小变化次数
# 只有r、e、d的字符串相邻相同字符如dd变成一个字符，最终相邻都不相等的最小变化次数

# 分类讨论
# !1.只有连续两个重复，那这两个变一次就行
# !2.超过三个连续的话最小次数就是隔一个变化后两个，也就是以三个为一组拆分
def solve(s: str) -> int:
    res = 0
    n = len(s)
    i = 0
    while i < n - 1:
        if s[i] == s[i + 1]:
            if i + 2 < n and s[i] == s[i + 2]:
                i += 2
                res += 1
            else:
                i += 1
                res += 1
        else:
            i += 1
    return res


if __name__ == "__main__":
    print(solve("redreedreeeed"))

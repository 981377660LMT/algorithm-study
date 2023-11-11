# 对字符串s 函数f(s)的值为
# !S の先頭の文字を削除し末尾に追加する操作を任意の回数行うことによって作ることのできる文字列の種類数
# !也就是字符串的最小周期
# 可以把任意k个字符替换成任意字符
# 求替换后的字符串的最小周期 因为n很小 所以直接枚举n的因子检查 O(n^4/3)
# n,k<=2000
# !枚举周期

from collections import defaultdict
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    s = input()

    for period in range(1, n):
        if n % period != 0:
            continue
        counter = defaultdict(lambda: defaultdict(int))  # 按mod统计每组个数
        for i in range(n):
            mod = i % period
            counter[mod][s[i]] += 1
        same = 0
        for mod in range(period):
            same += max(counter[mod].values(), default=0)  # 变为每个组的最大频率的字符
        if n - same <= k:
            print(period)
            exit(0)

    print(n)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()

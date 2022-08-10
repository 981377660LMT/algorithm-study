# https://garnacha.techblog.jp/archives/39854885.html
# !求s的字典序第k小的子序列 (子序列要去重)
# 不存在则输出"Eel"
# len(s)<=1e6
# k<=1e18


# dp复原
# dp[i][j - 'a'] ：文字列sの先頭からi文字目以降で、jという文字から始まる部分列の数

# ps:
# !如果是字典序第k小的子串 怎么做
# !二分出排名为K的子串是哪一个后缀的第几个未被计算过的前缀(每个后缀贡献子串数是这个后缀的长度减去其LCP)


import sys
from typing import Optional


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def solve(s: str, k: int, BASE=97) -> Optional[str]:
    """求s的字典序第k小的子序列 (子序列要去重) s由小写字母组成

    https://atcoder.jp/contests/tdpc/submissions/20602930
    """
    n = len(s)
    ords = [ord(c) - BASE for c in s]
    mins = [0] * n + [1, 0]
    nexts1 = [n + 1] * 26
    nexts2 = [0] * n
    for i in range(n - 1, -1, -1):
        num = ords[i]
        mins[i] = min(mins[i + 1] * 2 - mins[nexts1[num]], k + 1)
        nexts2[i] = nexts1[num]
        nexts1[num] = i + 1

    if mins[0] <= k:
        return None

    res = []
    index = 0
    while k:
        k -= 1
        count = 0
        j = 0
        while j < len(nexts1):
            v = nexts1[j]
            tmp = mins[v]
            if count + tmp > k:
                break
            count += tmp
            j += 1
        res.append(j)
        nextIndex = nexts1[j]
        for i in range(index, nextIndex):
            nexts1[ords[i]] = nexts2[i]
        index = nextIndex
        k -= count

    return "".join([chr(num + BASE) for num in res])


s = input()
k = int(input())

res = solve(s, k)
if res is None:
    print("Eel")
else:
    print(res)

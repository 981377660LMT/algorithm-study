"""
所有长为n的字符串中
RLE编码后比原来短的字符串有多少个
n<=3000 暗示O(n^2)的dp

dp[i][j] :i文字目まで決めて、RLEの長さがjの場合の文字列の通り数
遷移：ブロックを一つずつ追加
eg: 3 => 12 => 1 对应 +2 +3 +2 方案数对应 26*25*25
dp[i+k][j+getLRELen(k)]+=dp[i][j]*25 (1≤k≤N,k表示新追加的bloack的长度)

1-9 => 2
10-99 => 3
100-999 => 4
1000-9999 => 5

这个dp是O(n^3)的
用前缀和优化到O(n^2logn)
"""
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def getLRELen(blockLen: int) -> int:
    return len(str(blockLen)) + 1


def main() -> None:
    n, MOD = map(int, input().split())


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()


# TODO

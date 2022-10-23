from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# あなたはアメーバの観察記録をつけました。

# 最初 1 匹のアメーバがおり、番号は 1 です。

# 観察記録は時系列順に N 個あり、i 番目の観察記録は「番号 A
# i
# ​
#   のアメーバが分裂して消滅し、新たに 2 匹のアメーバが生まれ、それらにそれぞれ 2i,2i+1 と番号をつけた」というものです。
# このとき、アメーバ A
# i
# ​
#   を アメーバ 2i,2i+1 の親と呼びます。

# 各 k=1,…,2N+1 について、アメーバ k から何代親を遡るとアメーバ 1 になるか求めてください。
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    depth = defaultdict(int)
    for i in range(1, n + 1):
        cur = nums[i - 1]
        left, right = 2 * i, 2 * i + 1
        depth[left] = depth[cur] + 1
        depth[right] = depth[cur] + 1

    for i in range(1, 2 * n + 2):
        print(depth[i])

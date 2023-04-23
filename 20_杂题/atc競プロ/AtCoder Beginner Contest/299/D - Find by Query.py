from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# この問題は インタラクティブな問題（あなたが作成したプログラムとジャッジプログラムが標準入出力を介して対話を行う形式の問題）です。

# ジャッジが
# 0 と
# 1 のみからなる長さ
# N の文字列
# S=S
# 1
# ​
#  S
# 2
# ​
#  …S
# N
# ​
#   を持っています。 文字列
# S は、
# S
# 1
# ​
#  =0 および
# S
# N
# ​
#  =1 を満たします。

# あなたには
# S の長さ
# N が与えられますが、
# S 自体は与えられません。 その代わり、あなたはジャッジに対して以下の質問を
# 20 回まで行うことができます。

# 1≤i≤N を満たす整数
# i を選び、
# S
# i
# ​
#   の値を尋ねる。
# 1≤p≤N−1 かつ
# S
# p
# ​

# 
# =S
# p+1
# ​
#   を満たす整数
# p を
# 1 個出力してください。
# なお、本問題の条件下でそのような整数
# p が必ず存在することが示せます。


if __name__ == "__main__":
    n = int(input())
    left, right = 0, n - 2
    mp = defaultdict(int)
    mp[n - 1] = 1
    mp[0] = 0
    while left <= right:
        mid = (left + right) // 2
        print(f"? {mid+1}", flush=True)
        res = int(input())
        mp[mid] = res
        if res == 0:
            left = mid + 1
        else:
            right = mid - 1
    keys = sorted(mp.keys())
    for k1, k2 in zip(keys, keys[1:]):
        if k1 + 1 == k2 and mp[k1] ^ mp[k2]:
            print(f"! {k1+1}", flush=True)
            exit()

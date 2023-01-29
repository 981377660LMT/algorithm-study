import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 英小文字と ? からなる文字列
# S,T が与えられます。ここで、
# ∣S∣>∣T∣ が成り立ちます(文字列
# X に対し、
# ∣X∣ で
# X の長さを表します)。

# また、
# ∣X∣=∣Y∣ を満たす文字列
# X,Y は、次の条件を満たすとき及びそのときに限りマッチするといいます。

# X,Y に含まれる ? をそれぞれ独立に好きな英小文字に置き換えることで
# X と
# Y を一致させることができる
# x=0,1,…,∣T∣ に対して次の問題を解いてください。

# S の先頭の
# x 文字と末尾の
# ∣T∣−x 文字を順番を保ったまま連結することで得られる長さ
# ∣T∣ の文字列を
# S
# ′
#   とする。
# S
# ′
#   と
# T がマッチするならば Yes と、そうでなければ No と出力せよ。

if __name__ == "__main__":
    s = input()
    t = input()
    diff = 0
    lastT = s[-len(t) :]
    for a, b in zip(lastT, t):
        if a != b and a != "?" and b != "?":
            diff += 1
    print("Yes" if diff == 0 else "No")

    for x in range(1, len(t) + 1):
        # remove lastT[x-1]
        removed = lastT[x - 1]
        if removed != "?" and t[x - 1] != "?" and removed != t[x - 1]:
            diff -= 1
        # add s[x-1]
        added = s[x - 1]
        if added != "?" and t[x - 1] != "?" and added != t[x - 1]:
            diff += 1
        print("Yes" if diff == 0 else "No")

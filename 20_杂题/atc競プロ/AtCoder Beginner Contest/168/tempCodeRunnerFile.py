import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# いろはちゃんは、人気の日本製ゲーム「ÅtCoder」で遊びたい猫のすぬけ君のために日本語を教えることにしました。

# 日本語で鉛筆を数えるときには、助数詞として数の後ろに「本」がつきます。この助数詞はどんな数につくかで異なる読み方をします。具体的には、999 以下の正の整数 N について、「N 本」と言うときの「本」の読みは

# N の 1 の位が 2,4,5,7,9 のとき hon
# N の 1 の位が 0,1,6,8 のとき pon
# N の 1 の位が 3 のとき bon
# です。

# N が与えられるので、「N 本」と言うときの「本」の読みを出力してください。
if __name__ == "__main__":
    n = input()
    if n[0] in "24679":
        print("hon")
    elif n[0] in "0168":
        print("pon")
    else:
        print("bon")

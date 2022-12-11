from math import ceil
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

# みかんがたくさんあります。
# どのみかんの重さも A グラム以上 B グラム以下であることがわかっています。（みかんの重さは整数とは限りません。）
# この中からいくつかのみかんを選んだところ、選んだみかんの重さの合計がちょうど W キログラムになりました。
# 選んだみかんの個数として考えられる最小値と最大値を求めてください。ただし、このようなことが起こり得ないなら、かわりにそのことを報告してください。

# 選んだみかんの個数としてありえる最小値と最大値を空白区切りでこの順に出力せよ。
# ただし、与えられた条件に合うような個数が存在しない場合、かわりに UNSATISFIABLE と出力せよ。

# !A*n <= 1000*W <= B*n 求n的范围
if __name__ == "__main__":
    A, B, w = map(int, input().split())
    w *= 1000
    upper = w // A
    lower = ceil(w / B)
    if lower > upper:
        print("UNSATISFIABLE")
    else:
        print(lower, upper)

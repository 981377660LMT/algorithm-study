# N個の石があります。各ターンでは、今残っている石の数をa個とするとき、1個以上a/2個以下の石を取らなければなりません。また、初めて石を取れなくなったほうが負けです。
# 両者が最善を尽くした時、先手と後手どちらが勝つかを求めるプログラムを作成してください。
# 1<=N<=1e18
# !两个人每次的和凑成2^(n-1)

N = int(input())
flag = (N + 1) & N  # (N+1)是否为2的幂
if flag == 0:
    print("Second")
else:
    print("First")


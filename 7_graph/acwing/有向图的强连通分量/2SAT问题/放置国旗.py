# 1-N号旗设置位置(放置国旗)
# 第i号旗可以设置在xi位置或者yi位置
# !任意两面旗距离需要大于D
# 是否可以设置旗子

# 1≤N≤1000
# D,Xi,Yi<=1e9

# !定义 Ai 表示「可以选 Yi」，
# !这样若两个旗子 i j 满足 |Xi-Xj|<D 时，
# !就相当于 Ai与Aj至少满足一个，即 Ai或Aj为真，就可用addLimit api 加边
# https://atcoder.jp/contests/practice2/tasks/practice2_h

from TwoSat import TwoSat


n, d = map(int, input().split())
xy = [tuple(map(int, input().split())) for _ in range(n)]


ts = TwoSat(n)
for i in range(n):
    x1, y1 = xy[i]
    for j in range(i + 1, n):
        x2, y2 = xy[j]
        if abs(x1 - x2) < d:
            ts.addLimit(i, True, j, True)
            # ts.addEdge(i + n, j)  # !i为假可以推出j为真 (选x1就必须选y2)
        if abs(x1 - y2) < d:
            ts.addLimit(i, True, j, False)
            # ts.addEdge(i + n, j + n)  # !i为假可以推出j为假 (选x1就不能选y2)
        if abs(y1 - x2) < d:
            ts.addLimit(i, False, j, True)
            # ts.addEdge(i, j)  # !i为真可以推出j为真 (选y1就不能选x2)
        if abs(y1 - y2) < d:
            ts.addLimit(i, False, j, False)
            # ts.addEdge(i, j + n)  # !i为真可以推出j为假 (选y1就不能选y2)


ts.build()
if ts.check():
    print("Yes")
    res = ts.work()
    res = [xy[i][t] for i, t in enumerate(res)]  # true:命题为真，即可以选yi
    print(*res, sep="\n")
else:
    print("No")

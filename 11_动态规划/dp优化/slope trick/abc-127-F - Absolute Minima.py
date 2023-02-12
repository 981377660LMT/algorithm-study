# 开始时函数f(x)=0
# q个查询
# !1 a b: f(x)变为 g(x)=f(x)+|x-a|+b
# !2 求f(x)取到最小值时的x(如果有多个，输出最小的) 以及最小值
# q<=2e5

import sys
from SlopeTrick import SlopeTrick

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    q = int(input())
    st = SlopeTrick()
    for _ in range(q):
        t, *rest = map(int, input().split())
        if t == 1:
            a, b = rest
            st.add_abs(a)
            st.add_all(b)
        else:
            min_, leftTop, _ = st.query()
            print(leftTop, min_)

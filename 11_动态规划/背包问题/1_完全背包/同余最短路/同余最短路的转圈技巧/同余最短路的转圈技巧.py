# TODO
# https://www.cnblogs.com/alex-wei/p/17531487.html
# https://zhuanlan.zhihu.com/p/672216458
# 理解：体积很大的多重背包问题
#
# 有趣的是，同余最短路不应该从最短路的角度考虑。
# !其本质上是根据单调性值域定义域互换后将完全背包转化为体积模m意义下的完全背包。
# 普通完全背包的转移是有向无环图，而环上完全背包转移成环，这让我们想到最短路。
# !但因为一个点不会经过它自己，对应原问题就是对于一个物品，不会使得它的总体积为基准物品体积的倍数，
# 所以，我们可以将完全背包转化为类多重背包问题。

# constexpr ll MAXN = 5e5 + 5;
# ll n, l, r, f[MAXN];
# vector<int> a;
# il void solver_main() {
#   read(n, l, r);
#   For(i, 1, n) {
#     int x;
#     read(x);
#     if (x)
#       a.emplace_back(x);
#   }
#   if (a.empty())
#     return puts("0"), void();
#   sort(a.begin(), a.end());

#   int M = a[0], len = a.size(); // 选取最小值做 M，减小常数
#   fill(f, f + M, 1e18 + 5);
#   f[0] = 0;
#   For(i, 1, len - 1) {
#     int lim = __gcd(M, a[i]) - 1; // 环的数量
#     For(j, 0, lim) {
#       for (int cur = j, cnt = 0; cnt < 2; cnt += cur == j) { // 转两圈
#         int nxt = (cur + a[i]) % M;
#         f[nxt] = min(f[nxt], f[cur] + a[i]), cur = nxt;
#       }
#     }
#   }

#   ll ans = 0;
#   For(i, 0, M - 1) {
#     if (r >= f[i])
#       ans += (r - f[i]) / M + 1;
#     if (l > f[i])
#       ans -= (l - f[i] - 1) / M + 1;
#   }


#   cout << ans << endl;
# }


# P2371 [国家集训队] 墨墨的等式
# https://www.luogu.com.cn/problem/P2371
# 给定n个系数coeffs和上下界lower,upper
# !对于 lower<=k<=upper 求有多少个k能够满足
# !a0*x0+a1*x1+...+an*xn=k
# n<=12 0<=ai<=5e5 1<=lower<=upper<=2^63-1
# !时间复杂度：O(n*ai)
def p2371() -> None:
    n, lower, upper = map(int, input().split())
    coeffs = list(map(int, input().split()))
    coeffs = [v for v in coeffs if v != 0]
    if not coeffs:
        print(0)
        return


# P2662 牛场围栏(求最大的不能被线性表出的数)
# https://www.luogu.com.cn/problem/P2662

# P3403 跳楼机
# https://www.luogu.com.cn/problem/P3403
# 给定x,y,z,h
# !对于 1<=k<=h 求有多少个k能够满足
# !a*x+b*y+c*z=k
# (0<=x,y,z,1<=a,b,c<=1e5,h<=2^63-1)


# P9140 [THUPC 2023 初赛] 背包
# https://www.luogu.com.cn/problem/P9140

# D - Small Multiple
# https://atcoder.jp/contests/abc077/tasks/arc084_b
# k的倍数中 求各位数字之和的最小值
# k<=1e5.

if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    p2371()

# https://www.luogu.com.cn/problem/P3403
# 给定n个系数coeffs和上下界lower,upper
# !对于 lower<=k<=upper 求有多少个k能够满足
# !a0*x0+a1*x1+...+an*xn=k
# n<=12 0<=ai<=5e5 1<=lower<=upper<=2^63-1


from ModShortestPath import modShortestPath

INF = int(1e18)
n, lower, upper = map(int, input().split())
coeffs = sorted(map(int, input().split()))


def cal(start: int, step: int, left: int, right: int) -> int:
    """[left,right]闭区间中有多少形如start+step*k的数(k>=0)"""
    if left > right or start > right:
        return 0
    res1 = max(0, 1 + (right - start) // step)
    res2 = max(0, 1 + (left - 1 - start) // step)
    return res1 - res2


# !每个剩余类的贡献 dist[i],dist[i]+base,...,dist[i]+k*base 都能到达
base, dist = modShortestPath(coeffs)
print(sum(cal(start, base, lower, upper) for start in dist))

# https://www.luogu.com.cn/problem/P3403
# 给定x,y,z,h
# !对于 1<=k<=h 求有多少个k能够满足
# !a*x+b*y+c*z=k
# (0<=x,y,z,1<=a,b,c<=1e5,h<=2^63-1)


from ModShortestPath import modShortestPath

INF = int(1e18)
h = int(input())
a, b, c = sorted(map(int, input().split()))


base, dist = modShortestPath([a, b, c])
dist = [v + 1 for v in dist]  # !从1楼开始上跳


def cal(start: int, step: int, left: int, right: int) -> int:
    """[left,right]闭区间中有多少形如start+step*k的数(k>=0)"""
    if left > right or start > right:
        return 0
    res1 = max(0, 1 + (right - start) // step)
    res2 = max(0, 1 + (left - 1 - start) // step)
    return res1 - res2


# !每个剩余类的贡献 dist[i],dist[i]+base,...,dist[i]+k*base 都能到达
print(sum(cal(start, base, 1, h) for start in dist))

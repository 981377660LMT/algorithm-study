import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# D - Congruence Points
# 点集相等/点集重合
# 给定两个二维平面,各有 n 个点,问,
# 第一个平面上的点经过若干次整体平移或整体旋转(任意角度),能否和第二个平面的点重合
# n<=100

# 三维重建的点云匹配
# 什么是点云配准
# https://zhuanlan.zhihu.com/p/397926700

# !解法1:重心+枚举旋转度数
# https://blog.csdn.net/weixin_43989731/article/details/118309878
# !如果这两个点集可以通过旋转匹配，那么，这两个点集中的点相对重心的位置分布是相同的。
# !重心:每个点加权要一样
# 因此，我们可以先求出这两个点集中的点相对其重心的相对坐标，这样我们就能够忽略掉平移这个操作所带来的影响。
# 下面就要考虑能否通过旋转使两个点集匹配。
# 从点集S 中选一个x 坐标不为0的点p ，这样可以保证有tan值。
# !由于这题范围很小，因此直接暴力枚举点集T中的所有点转到点p所需的角度θ
# 然后，判断点集S的剩余点转θ能否与点集T中的某一个点匹配即可。
# !旋转:复数乘法就是逆时针旋转

EPS = 1e-6
if __name__ == "__main__":
    n = int(input())
    if n == 1:
        print("Yes")
        exit(0)

    points1 = [complex(*map(int, input().split())) for _ in range(n)]
    points2 = [complex(*map(int, input().split())) for _ in range(n)]

    # 重心平移到原点之后 只需判断能否旋转得到
    center1, center2 = sum(points1) / n, sum(points2) / n
    points1 = [point - center1 for point in points1]
    points2 = [point - center2 for point in points2]

    p0 = [p for p in points1 if p][0]  # 选一个坐标不为(0,0)(重心)的点
    for p1 in points2:  # 假设p0对应p1
        for p2 in points1:
            for p3 in points2:
                if abs(p2 * p1 - p3 * p0) < EPS and abs(abs(p2) - abs(p3)) < EPS:  # p0=>p1 p2=>p3
                    break
            else:
                break
        else:
            print("Yes")
            exit(0)

    print("No")

# !解法2:枚举距离对
# https://kanpurin.hatenablog.com/entry/2021/06/26/224005
# !解法3
# O(nlogn)的解法
# https://atcoder.jp/contests/abc207/editorial/2167

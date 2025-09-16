# C - Lock All Doors(关门开锁)
# https://atcoder.jp/contests/abc423/tasks/abc423_c
# 有 N+1 个房间，编号 0 到 N，排成一列。房间之间有 N 个门，编号 1 到 N。门 i 连接房间 i-1 和 i。
# 每个门 i 有一个状态 L_i：0 表示开着，1 表示锁着。
# 高桥君初始在房间 R。他只能通过开着的门在相邻房间移动。当他位于房间 i-1 或 i 时，可以对门 i 进行“开/关操作”，这会翻转门的状态（开变关，关变开）。
# 目标是使所有门都变为锁上状态。请求出所需的最少“开/关操作”次数。
#
# !问题转化为：为了能够到达所有初始为“开”的门的位置，高桥君最少需要临时打开几扇已锁上的门？


if __name__ == "__main__":
    N, S = map(int, input().split())
    A = list(map(int, input().split()))
    pos0 = [i for i, v in enumerate(A) if v == 0]
    if not pos0:
        print(0)
        exit(0)

    L = min(S, pos0[0])
    R = max(S, pos0[-1] + 1)
    print(sum(A[L:R]) + (R - L))

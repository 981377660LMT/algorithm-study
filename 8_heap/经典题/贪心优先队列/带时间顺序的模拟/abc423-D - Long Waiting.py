# D - Long Waiting(餐厅排队、优先队列模拟、带时间顺序的模拟)
# https://atcoder.jp/contests/abc423/tasks/abc423_d
#
# 有一家最多可同时容纳 K 位顾客的餐厅。餐厅外有一条唯一的排队队列。
# 在时刻 0，店内没有客人，队列也是空的。
# 今天预计有 N 个团体客人光顾，按到达顺序编号为 1 到 N。团体 i 有 C_i 人，在时刻 A_i 到达并加入队尾。入店后，该团体将在 B_i 个单位时间后离店。
# 每个团体将在满足以下两个条件的 最早时刻 离开队列并进入餐厅：
# 该团体位于队列的最前端。
# 该团体的人数，加上店内所有已在客人（包括同一时刻刚入店的，不包括刚离店的）的人数总和，不超过餐厅容量 K。
# 请为每个团体计算他们入店的时刻。

from heapq import heapify, heappop, heappush


if __name__ == "__main__":
    N, K = map(int, input().split())
    A, B, C = [0] * N, [0] * N, [0] * N
    for i in range(N):
        A[i], B[i], C[i] = map(int, input().split())

    people = 0
    waits = []  # index
    events = [(A[i], 0, i) for i in range(N)]  # (time, in/out, index)
    heapify(events)
    res = []
    while events:
        curTime, inout, index = heappop(events)
        if inout == 0:
            heappush(waits, index)
        else:
            people -= C[index]

        while waits:
            index = waits[0]
            if people + C[index] <= K:
                res.append(curTime)
                people += C[index]
                heappush(events, (curTime + B[index], 1, index))
                heappop(waits)
            else:
                break

    print("\n".join(map(str, res)))

# C - New Skill Acquired
# https://atcoder.jp/contests/abc424/tasks/abc424_c
#
# 问题描述
#
# 高桥君正在玩一个游戏。游戏中有 N 个技能，编号从 1 到 N。
#
# 给定 N 对整数 (A_1, B_1), ..., (A_N, B_N)。
#
# 如果 (A_i, B_i) = (0, 0)，表示高桥君已经学会了技能 i。
# 否则，高桥君当且仅当已经学会了技能 A_i 和技能 B_i 中的至少一个时，才能学会技能 i。
# 请计算包括已经学会的技能在内，高桥君最终能够学会的技能总数。


from collections import deque


if __name__ == "__main__":
    N = int(input())
    A, B = [0] * N, [0] * N
    for i in range(N):
        A[i], B[i] = map(int, input().split())
        A[i] -= 1
        B[i] -= 1

    adjList = [[] for _ in range(N)]
    queue = deque()
    for i, (a, b) in enumerate(zip(A, B)):
        if a == -1 and b == -1:
            queue.append(i)
        else:
            adjList[a].append(i)
            adjList[b].append(i)

    visited = [False] * N
    while queue:
        cur = queue.popleft()
        if visited[cur]:
            continue
        visited[cur] = True
        for next_ in adjList[cur]:
            queue.append(next_)

    print(sum(visited))

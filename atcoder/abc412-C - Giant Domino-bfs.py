# C - Giant Domino
# https://atcoder.jp/contests/abc412/tasks/abc412_c
# 有编号 1…N 的 N 块多米诺骨牌，第 i 块骨牌的大小为 S[i]。
# 选出至少 2 块骨牌，按某个顺序一字排开，最左放编号 1，最右放编号 N。
# 当仅推倒骨牌 1 时，若当前倒下的骨牌 i 右侧紧邻的骨牌 j 满足 S[j] ≤ 2⋅S[i]，
# 则它也会倒下。要求推倒最终能倒下骨牌 N。问是否存在这样的排法？
# 若存在，最少要排多少块骨牌，否则输出 −1。
# 共有 T 个测试，∑N≤2⋅10^5。


from collections import deque


if __name__ == "__main__":

    def solve():
        N = int(input())
        S = list(map(int, input().split()))

        pairs = sorted((v, i) for i, v in enumerate(S))
        visited = [False] * N
        dist = [0] * N

        visited[0] = True
        dist[0] = 1
        queue = deque([0])

        ptr = 0
        while queue:
            cur = queue.popleft()
            if cur == N - 1:
                print(dist[cur])
                return
            limit = 2 * S[cur]
            while ptr < N and pairs[ptr][0] <= limit:
                j = pairs[ptr][1]
                ptr += 1
                if not visited[j]:
                    visited[j] = True
                    dist[j] = dist[cur] + 1
                    queue.append(j)

        print(-1)

    T = int(input())
    for _ in range(T):
        solve()

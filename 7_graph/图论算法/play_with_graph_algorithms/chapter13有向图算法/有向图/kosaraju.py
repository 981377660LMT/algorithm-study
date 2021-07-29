# https://www.pythonf.cn/read/147039
N = 7
graph, rgraph = [[] for _ in range(N)], [[] for _ in range(N)]
used = [False for _ in range(N)]
popped = []


# 建图
def add_edge(u, v):
    graph[u].append(v)
    rgraph[v].append(u)


# 正向遍历
def dfs(u):
    used[u] = True
    for v in graph[u]:
        if not used[v]:
            dfs(v)
    popped.append(u)


# 反向遍历
def rdfs(u, scc):
    used[u] = True
    scc.append(u)
    for v in rgraph[u]:
        if not used[v]:
            rdfs(v, scc)


# 建图，测试数据
def build_graph():
    add_edge(1, 3)
    add_edge(1, 2)
    add_edge(2, 4)
    add_edge(3, 4)
    add_edge(3, 5)
    add_edge(4, 1)
    add_edge(4, 6)
    add_edge(5, 6)


if __name__ == "__main__":
    build_graph()
    for i in range(1, N):
        if not used[i]:
            dfs(i)

    used = [False for _ in range(N)]
    # 将第一次dfs出栈顺序反向
    popped.reverse()
    for i in popped:
        if not used[i]:
            scc = []
            rdfs(i, scc)
            print(scc)

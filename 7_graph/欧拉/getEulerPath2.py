# 构建有向图的欧拉路径


from typing import List, Optional, Tuple


def hierholzer(n: int, edges: List[Tuple[int, int]]) -> Optional[List[int]]:
    m = len(edges)
    G = [[] for _ in range(n)]
    deg = [0] * n
    rdeg = [0] * n
    for a, b in edges:
        deg[a] += 1
        rdeg[b] += 1
        G[a].append(b)

    # find starting and ending vertices
    s = t = u = -1
    for i in range(n):
        if deg[i] == rdeg[i] == 0:
            continue
        df = deg[i] - rdeg[i]
        if not -1 <= df <= 1:
            return None
        if df == 1:
            if s != -1:
                return None
            s = i
        elif df == -1:
            if t != -1:
                return None
            t = i
        else:
            u = i
    v0 = s if s != -1 else u

    # find an Eulerian path (or circuit)
    res = []
    it = [0] * n
    st = [v0]
    (*it,) = map(iter, G)
    while st:
        v = st[-1]
        w = next(it[v], -1)
        if w == -1:
            res.append(v)
            st.pop()
            continue
        st.append(w)
    res.reverse()
    if len(res) != m + 1:
        return None
    return res


E = [
    (2, 0),
    (0, 3),
    (0, 1),
    (3, 2),
    (1, 2),
]
print(hierholzer(4, E))
# => "[0, 3, 2, 0, 1, 2]"

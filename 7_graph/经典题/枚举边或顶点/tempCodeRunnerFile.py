            deg[u] += 1
            deg[v] += 1

        for u, v in edges:
            u, v = sorted((u, v), key=lambda x: (deg[x], x))
            adjMap[u].append(v)

        res = 0
        for p1 in range(1, n + 1):
            for p2 in adjMap[p1]:
                for p3 in adjMap[p2]:
                    if p3 in adjMap[p1]:
                        res += 1
        return res
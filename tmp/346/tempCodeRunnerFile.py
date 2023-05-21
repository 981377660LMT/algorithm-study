            # check
            newAdjList = [[] for _ in range(n)]
            for u, v, w in edges:
                newAdjList[u].append((v, w))
                newAdjList[v].append((u, w))
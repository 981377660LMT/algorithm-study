
        (2, 0),
        (0, 3),
        (0, 1),
        (3, 2),
        (1, 2),
    ]

    allVertex = {0, 1, 2, 3}
    adjMap = defaultdict(set)
    for a, b in E:
        adjMap[a].add(b)
        adjMap[b].add(a)

    print(getEulerPath(allVertex, adjMap, isDirected=False))

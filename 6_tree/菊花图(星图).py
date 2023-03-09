# generate a tree with 1e5 nodes.
# eg: edges: (0,1),(0,2),(0,3)....(0,1e5-1)


# 生成菊花图(星图)数据
with open("starg.txt", "w") as f:
    n = int(1e5)
    edges = [[0, i] for i in range(1, n)]
    guess = [[0, i] for i in range(1, n)]
    k = 1

    sb = []
    for e in edges:
        sb.append(str(e))
    f.write("[" + ",".join(sb) + "]\n")

    sb = []
    for g in guess:
        sb.append(str(g))
    f.write("[" + ",".join(sb) + "]\n")

    f.write(str(k))

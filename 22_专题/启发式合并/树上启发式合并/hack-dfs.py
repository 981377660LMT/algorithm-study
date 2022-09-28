# hack启发式合并
n = int(3e4)
# 14999,14888,...,1,0,0,1,2,...,14999
vals = [*range(14999, -1, -1), *range(15000)]
edges = []
for i in range(n - 1):
    edges.append([i, i + 1])

with open("big-head-center-tail-with-symmetric-vals-chain.txt", "w") as f:
    f.write(str(vals) + "\n")
    f.write(str(edges) + "\n")

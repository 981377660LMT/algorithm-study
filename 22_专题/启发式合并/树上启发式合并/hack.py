n = int(3e4)
vals = []
for i in range(7500, -1, -1):  # 7500 7499 ... 1 0
    vals.append(i)
vals.append(30000)
for i in range(7501):  # 0 1 ... 7498 7500
    vals.append(i)
while len(vals) < n:
    vals.append(vals[-1] + 1)

edges = []
for i in range(n):
    edges.append(i)

########################################
# x: edge y:vals
import numpy as np
import matplotlib.pyplot as plt

x = np.array(edges)
y = np.array(vals)

# draw x y
plt.plot(x, y, "o")
plt.show()

#################################################################################
# print vals and edges to a file like
# [1,1,2,2,3]
# [[0,1],[1,2],[2,3],[2,4]]
with open("big-head-center-tail-with-symmetric-vals-chain.txt", "w") as f:
    f.write(str(vals) + "\n")
    f.write(str(edges) + "\n")

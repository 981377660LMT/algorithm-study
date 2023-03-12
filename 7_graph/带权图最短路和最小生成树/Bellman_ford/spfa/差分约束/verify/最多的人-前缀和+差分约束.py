# 水平的列车上有n个座位，从左到右座位号为1,2...n。
# !现在有m条规定，每条规定的形式如下∶从座位l到座位r，不多于x个人乘坐。
# 在满足所有规定的前提下，该列车最多能坐多少人?

# Sright - Sleft-1 <= w
# Si - Si-1 >= 0
# Si - Si-1 <= 1


from 差分约束 import DualShortestPath


n, m = map(int, input().split())
D = DualShortestPath(n + 10, min=False)
for _ in range(m):
    left, right, w = map(int, input().split())
    D.addEdge(right, left - 1, w)

# !前缀和满足的约束
for i in range(1, n + 5):  # 多加一点
    D.addEdge(i - 1, i, 0)
    D.addEdge(i, i - 1, 1)


res, ok = D.run()
print(res[n])

# 10 3
# 1 4 2
# 3 6 2
# 10 10 1
# 输出8

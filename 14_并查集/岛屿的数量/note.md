并查集的做法：
遍历矩阵:对每个 1 ,union 左侧和上方的 1

```Python
        # 初始化并查集
        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    for x, y in [[-1, 0], [0, -1]]:
                        tmp_i = i + x
                        tmp_j = j + y
                        if 0 <= tmp_i < row and 0 <= tmp_j < col and grid[tmp_i][tmp_j] == "1":
                            union(tmp_i * row + tmp_j, i * row + j)
        # print(f)
        res = set()
        for i in range(row):
            for j in range(col):
                if grid[i][j] == "1":
                    res.add(find((i * row + j)))
        return len(res)

```

`复杂827. 最大人工岛`

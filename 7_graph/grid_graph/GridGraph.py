# 多维坐标转换


class GridGraph:
    __slots__ = ("dims", "dim", "volume", "partition")

    def __init__(self, *dims):
        """
        多维坐标转换.

        dims : 各维度的大小
        """
        self.dims = tuple(dims)
        self.dim = len(self.dims)
        r = [1]
        for a in self.dims[::-1]:
            r.append(r[-1] * a)
        self.volume, *self.partition = r[::-1]

    def number_to_position(self, n):
        """将编号为 n 的点转换为坐标."""
        assert 0 <= n < self.volume

        pos = [0] * self.dim
        for i in range(self.dim):
            pos[i], n = n // self.partition[i], n % self.partition[i]
        return pos

    def position_to_number(self, *pos):
        """将坐标转换为编号."""
        assert len(pos) == self.dim

        res = 0
        for i in range(self.dim):
            assert 0 <= pos[i] < self.dims[i]
            res += self.partition[i] * pos[i]
        return res

    def number_neighborhood_yielder(self, n):
        """生成编号为 n 的点的邻居."""
        r = n
        for i in range(self.dim):
            q, r = r // self.partition[i], r % self.partition[i]

            if 0 < q:
                yield n - self.partition[i]

            if q < self.dims[i] - 1:
                yield n + self.partition[i]

    def position_neighborhood_yielder(self, *pos):
        """生成坐标为 pos 的点的邻居."""
        assert self.dim == len(pos)
        pos = list(pos)

        for i in range(self.dim):
            if 0 < pos[i]:
                pos[i] -= 1
                yield pos
                pos[i] += 1

            if pos[i] < self.dims[i] - 1:
                pos[i] += 1
                yield pos
                pos[i] -= 1


if __name__ == "__main__":

    def demo():
        # 示例：创建一个 3x4 的网格（二维网格）
        grid = GridGraph(3, 4)
        print("网格维度:", grid.dim)
        print("网格总体积:", grid.volume)
        print("每个维度大小:", grid.dims)
        print("每个维度的 partition 权重:", grid.partition)

        # 将编号转换为位置
        N = 7
        pos = grid.number_to_position(N)
        print(f"编号 {N} 对应的位置: {pos}")

        # 将位置转换为编号（验证与上面的编号是否一致）
        N2 = grid.position_to_number(*pos)
        print(f"位置 {pos} 对应的编号: {N2}")

        # 输出编号 N 的邻居编号（在一维编号空间中）
        print(f"编号 {N} 的邻居编号:")
        for neighbor in grid.number_neighborhood_yielder(N):
            print(neighbor, "-> 位置", grid.number_to_position(neighbor))

        # 输出位置 pos 的邻居位置（直接操作多维坐标）
        print(f"位置 {pos} 的邻居位置:")
        for neighbor_pos in grid.position_neighborhood_yielder(*pos):
            print(neighbor_pos, "-> 对应编号", grid.position_to_number(*neighbor_pos))

    demo()

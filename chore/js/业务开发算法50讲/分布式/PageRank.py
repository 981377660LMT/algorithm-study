# 基于马尔可夫链（Markov Chain）实现 PageRank 算法

import numpy as np


class PageRank:
    def __init__(self, adjacency_list, damping_factor=0.85, tolerance=1e-8, max_iterations=100):
        """
        初始化 PageRank 模型。

        :param adjacency_list: 有向图的邻接表，字典形式。例如：{0: [1, 2], 1: [2], 2: [0], 3: [2, 7], ...}
        :param damping_factor: 阻尼系数，通常为 0.85，表示用户有 85% 的概率继续点击链接，15% 的概率跳转到任意网页。
        :param tolerance: 收敛阈值
        :param max_iterations: 最大迭代次数
        """
        self.adjacency_list = adjacency_list
        self.d = damping_factor
        self.tolerance = tolerance
        self.max_iterations = max_iterations
        self.N = len(adjacency_list)
        self.transition_matrix = self.build_transition_matrix()

    def build_transition_matrix(self):
        """
        构建转移概率矩阵。

        :return: 转移概率矩阵的 NumPy 数组
        """
        P = np.zeros((self.N, self.N))
        for node in self.adjacency_list:
            outgoing = self.adjacency_list[node]
            if outgoing:
                prob = self.d / len(outgoing)
                for target in outgoing:
                    P[target][node] += prob  # 注意转置
            else:
                # !死胡同，随机跳转到所有页面，解决悬空链接(dangling links)问题
                prob = self.d / self.N
                P[:, node] += prob
        # !加上随机跳转的部分，解决蜘蛛陷阱(spider traps)问题
        random_jump = (1 - self.d) / self.N
        P += random_jump
        return P.T  # 转置以符合向量左乘的形式

    def compute_pagerank(self):
        """
        计算 PageRank 值。

        :return: PageRank 的 NumPy 数组
        """
        # 初始化 PageRank 向量
        PR = np.ones(self.N) / self.N
        for iteration in range(self.max_iterations):
            PR_new = PR @ self.transition_matrix
            # 检查收敛性
            if np.linalg.norm(PR_new - PR, ord=1) < self.tolerance:
                print(f"PageRank 在 {iteration + 1} 次迭代后收敛。")
                return PR_new
            PR = PR_new
        print("警告：PageRank 未在最大迭代次数内收敛。")
        return PR


# 示例使用
if __name__ == "__main__":
    # 定义有向图的邻接表
    # 例如：网页 0 链接到 1 和 2，网页 1 链接到 2，网页 2 链接到 0，网页 3 链接到 2 和 7 等等
    adjacency_list = {0: [1, 2], 1: [2], 2: [0], 3: [2, 7], 4: [0, 2], 5: [1], 6: [1], 7: [0]}

    # 实例化 PageRank 模型
    pagerank_model = PageRank(adjacency_list)

    # 计算 PageRank
    pagerank_scores = pagerank_model.compute_pagerank()
    print(sum(pagerank_scores))

    # 输出结果
    for node, score in enumerate(pagerank_scores):
        print(f"网页 {node} 的 PageRank 值：{score:.6f}")

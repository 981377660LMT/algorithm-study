import numpy as np
from collections import Counter
import math
import bisect


def calculate_cost(colors, counts, start, end, p=2):
    """
    计算将灰度值[start:end]映射到同一颜色的最小惩罚

    参数:
        colors: 排序后的不同灰度值列表
        counts: 每个灰度值的出现次数
        start, end: 灰度值范围的索引
        p: p范数的参数

    返回:
        (min_penalty, optimal_color): 最小惩罚值和对应的最优颜色
    """
    if start > end:
        return 0, 0

    total_pixels = sum(counts[i] for i in range(start, end + 1))

    # 对于不同的p范数，最优颜色的选择策略不同
    if p == 1:  # L1范数: 最优值是中位数
        cumsum = 0
        for i in range(start, end + 1):
            cumsum += counts[i]
            if cumsum >= total_pixels / 2:
                optimal_color = colors[i]
                break
    elif p == 2:  # L2范数: 最优值是加权平均值
        weighted_sum = sum(colors[i] * counts[i] for i in range(start, end + 1))
        optimal_color = round(weighted_sum / total_pixels)
    else:  # 一般p范数: 通过遍历找最优值
        min_penalty = float("inf")
        optimal_color = 0
        for candidate_idx in range(start, end + 1):
            candidate = colors[candidate_idx]
            penalty = sum(
                counts[i] * (abs(colors[i] - candidate) ** p) for i in range(start, end + 1)
            )
            if penalty < min_penalty:
                min_penalty = penalty
                optimal_color = candidate

    # 计算最优颜色对应的惩罚值
    min_penalty = sum(
        counts[i] * (abs(colors[i] - optimal_color) ** p) for i in range(start, end + 1)
    )

    return min_penalty, optimal_color


def find_optimal_palette_monotone_dp(colors, counts, K, p=2):
    """
    使用具有决策单调性的动态规划找出最优的K个调色板颜色

    参数:
        colors: 排序后的不同灰度值列表
        counts: 每个灰度值的出现次数
        K: 调色板大小
        p: p范数的参数

    返回:
        (min_penalty, optimal_palette): 最小惩罚值和最优调色板
    """
    C = len(colors)

    # 预计算所有cost(i,j)和对应的最优颜色
    cost = [[0] * C for _ in range(C)]
    opt_color = [[0] * C for _ in range(C)]

    for i in range(C):
        for j in range(i, C):
            cost[i][j], opt_color[i][j] = calculate_cost(colors, counts, i, j, p)

    # dp[i][j]: 使用i个颜色覆盖前j+1个灰度值的最小惩罚
    dp = [[float("inf")] * C for _ in range(K + 1)]
    # 决策点: 上一个颜色的结束位置
    prev = [[0] * C for _ in range(K + 1)]

    # 基础情况: 只用1个颜色
    for j in range(C):
        dp[1][j] = cost[0][j]

    # 使用决策单调性优化的DP
    for i in range(2, K + 1):
        # 分治DP优化函数
        def solve(j_start, j_end, k_start, k_end):
            if j_start > j_end:
                return

            j_mid = (j_start + j_end) // 2

            # 找到dp[i][j_mid]的最优决策点
            best = k_start
            best_val = float("inf")

            for k in range(k_start, min(k_end + 1, j_mid + 1)):
                val = dp[i - 1][k - 1] + cost[k][j_mid]
                if val < best_val:
                    best_val = val
                    best = k

            dp[i][j_mid] = best_val
            prev[i][j_mid] = best

            # 递归处理左半部分和右半部分
            solve(j_start, j_mid - 1, k_start, best)
            solve(j_mid + 1, j_end, best, k_end)

        solve(0, C - 1, 0, C - 1)

    # 回溯找出最优调色板
    palette = []
    j = C - 1
    for i in range(K, 0, -1):
        k = prev[i][j]
        if i > 1:
            palette.append(opt_color[k][j])
            j = k - 1
        else:
            palette.append(opt_color[0][j])

    return dp[K][C - 1], sorted(palette)


def compress_image(image, palette):
    """
    使用给定调色板压缩图像

    参数:
        image: 原始灰度图像
        palette: 调色板（K个灰度值）

    返回:
        (compressed_image, indices): 压缩后的图像和映射索引
    """
    H, W = image.shape
    compressed = np.zeros_like(image)
    indices = np.zeros((H, W), dtype=np.uint8)

    # 对于每个像素，找到最接近的调色板颜色
    for i in range(H):
        for j in range(W):
            pixel = image[i, j]
            # 二分查找最接近的调色板颜色
            idx = bisect.bisect_left(palette, pixel)
            if idx == 0:
                best_idx = 0
            elif idx == len(palette):
                best_idx = len(palette) - 1
            else:
                if abs(palette[idx] - pixel) < abs(palette[idx - 1] - pixel):
                    best_idx = idx
                else:
                    best_idx = idx - 1

            compressed[i, j] = palette[best_idx]
            indices[i, j] = best_idx

    return compressed, indices


def quantize_grayscale_image(image, K, p=2):
    """
    将灰度图像量化为K个颜色

    参数:
        image: 原始灰度图像，HxW的numpy数组
        K: 要使用的颜色数
        p: p范数参数

    返回:
        (compressed, palette, indices): 压缩后的图像、调色板和索引
    """
    # 提取图像中的所有灰度值和频率
    flat_image = image.flatten()
    color_counts = Counter(flat_image)

    colors = sorted(color_counts.keys())
    counts = [color_counts[c] for c in colors]

    # 找出最优调色板
    min_penalty, optimal_palette = find_optimal_palette_monotone_dp(colors, counts, K, p)

    # 使用最优调色板压缩图像
    compressed, indices = compress_image(image, optimal_palette)

    return compressed, optimal_palette, indices


def main():
    """演示函数"""
    # 生成测试图像
    H, W = 100, 100
    Q = 8  # 8位灰度 (0-255)
    np.random.seed(42)
    image = np.random.randint(0, 2**Q, size=(H, W), dtype=np.uint8)

    # 量化到K个颜色
    K = 8
    p = 2  # L2范数

    print(f"正在将{H}x{W}的{Q}位灰度图像压缩为{K}种颜色...")
    compressed, palette, indices = quantize_grayscale_image(image, K, p)

    # 计算压缩效率
    original_size = H * W * Q
    palette_size = K * Q
    index_size = H * W * math.ceil(math.log2(K))
    compressed_size = index_size + palette_size

    print(f"原始图像大小: {original_size} bits")
    print(f"压缩图像大小: {compressed_size} bits")
    print(f"压缩率: {compressed_size/original_size:.4f}")
    print(f"最优调色板: {palette}")

    # 计算失真度
    mse = np.mean((image - compressed) ** 2)
    psnr = 10 * math.log10((2**Q - 1) ** 2 / mse) if mse > 0 else float("inf")
    print(f"均方误差: {mse:.2f}")
    print(f"峰值信噪比: {psnr:.2f} dB")


if __name__ == "__main__":
    main()

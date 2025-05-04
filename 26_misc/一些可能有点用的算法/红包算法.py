import random


def distribute_red_packet(N, K):
    """
    将N元红包公平分配给K个人

    参数:
        N (int): 红包金额（单位：元）
        K (int): 人数

    返回:
        list: 每个人获得的金额列表（单位：元）
    """
    # 将金额转换为分
    N_cents = int(N * 100)

    # 使用组合选择算法选择K-1个分隔点
    # !在N*100+K-1个位置中选择K-1个作为分隔点
    sep_points = random.sample(range(N_cents + K - 1), K - 1)

    # 添加边界点
    sep_points = [-1] + sorted(sep_points) + [N_cents + K - 2]

    # 计算每个人获得的金额（分）
    amounts_cents = [sep_points[i + 1] - sep_points[i] - 1 for i in range(K)]

    # 转换回元为单位，保留2位小数
    amounts = [amount / 100 for amount in amounts_cents]

    return amounts


# 测试代码
def test_distribution():
    """测试红包分配的均匀性"""
    N, K = 1, 3  # 1元分给3个人
    total_cents = N * 100

    # 统计不同分配方案的出现次数
    distribution = {}
    trials = 50000

    for _ in range(trials):
        result = distribute_red_packet(N, K)
        # 将结果转换为整数分并排序，作为键
        key = tuple(sorted([int(x * 100) for x in result]))
        distribution[key] = distribution.get(key, 0) + 1

    # 计算理论上的总组合数
    import math

    total_combinations = math.comb(total_cents + K - 1, K - 1)

    print(f"N={N}元, K={K}人")
    print(f"理论组合数: {total_combinations}")
    print(f"实际观察到的不同组合数: {len(distribution)}")

    # 检查分布是否均匀
    expected = trials / total_combinations
    max_deviation = (
        max(abs(freq - expected) / expected for freq in distribution.values())
        if distribution
        else 0
    )
    print(f"期望每种组合出现次数: {expected:.2f}")
    print(f"最大偏差率: {max_deviation:.2%}")

    # 打印前10种组合及其出现次数
    print("\n前10种组合及其出现次数:")
    for i, (key, count) in enumerate(
        sorted(distribution.items(), key=lambda x: x[1], reverse=True)
    ):
        if i >= 10:
            break
        # 转换回元显示
        amounts = [x / 100 for x in key]
        print(f"{amounts}: {count}次 ({count/trials*100:.2f}%)")


if __name__ == "__main__":
    # 示例用法
    print("示例1: 10元分给3人")
    result1 = distribute_red_packet(10, 3)
    # 格式化为2位小数显示
    formatted1 = [f"{amount:.2f}" for amount in result1]
    print(formatted1)

    print("\n示例2: 1元分给5人")
    result2 = distribute_red_packet(1, 5)
    formatted2 = [f"{amount:.2f}" for amount in result2]
    print(formatted2)

    print("\n示例3: 0.88元分给4人")
    result3 = distribute_red_packet(0.88, 4)
    formatted3 = [f"{amount:.2f}" for amount in result3]
    print(formatted3)

    # 如需检验分布均匀性，取消下行注释
    test_distribution()

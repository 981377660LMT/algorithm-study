# random.sample

import random


def random_set(N, K):
    """
    从0到N-1中随机选出K个不同的元素，每种选法等概率.
    O(K)时间复杂度.

    参数:
        N (int): 数字范围上限（0到N-1）
        K (int): 要选择的不同元素数量，K <= N

    返回:
        list: 包含K个不同整数的列表
    """
    if K > N:
        raise ValueError("K不能大于N")

    if K == 0:
        return []

    # 策略选择点：当2K > N时，洗牌算法更优；否则随机选择算法更优
    if 2 * K > N:
        # 算法1：洗牌算法 (Fisher-Yates) - O(N)
        nums = list(range(N))
        # 只需要洗前K个位置
        for i in range(K):
            # 随机选择i到N-1之间的一个位置并与i位置交换
            j = random.randint(i, N - 1)
            nums[i], nums[j] = nums[j], nums[i]

        return nums[:K]
    else:
        # 算法2：随机选择 - O(K)期望时间
        result = set()
        while len(result) < K:
            # 随机选择0到N-1之间的一个数
            num = random.randint(0, N - 1)
            result.add(num)

        return list(result)


# 测试代码
def test_distribution(N, K, trials=10000):
    """测试随机选择的分布"""
    count = {}

    for _ in range(trials):
        # 排序以便相同的组合映射到相同的键
        result = tuple(sorted(random_set(N, K)))
        count[result] = count.get(result, 0) + 1

    # 计算理论概率
    total_combinations = 1
    for i in range(K):
        total_combinations *= (N - i) / (i + 1)
    total_combinations = int(total_combinations)

    # 检查分布是否均匀
    expected = trials / total_combinations
    max_deviation = max(abs(freq - expected) / expected for freq in count.values())

    print(f"N={N}, K={K}")
    print(f"总组合数: {total_combinations}")
    print(f"不同组合出现数: {len(count)}")
    print(f"期望每种组合出现次数: {expected:.2f}")
    print(f"最大偏差: {max_deviation:.2%}")

    return max_deviation < 0.1  # 假设偏差小于10%是可接受的


# 运行测试
if __name__ == "__main__":
    # 测试小数据
    print(random_set(5, 2))  # 例如[0, 3]
    print(random_set(10, 7))  # 例如[0, 1, 3, 5, 7, 8, 9]

    # 如果需要验证分布均匀性（可选）
    test_distribution(5, 2, trials=10000)
    test_distribution(10, 3, trials=10000)

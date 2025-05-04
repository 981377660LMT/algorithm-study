# 桌游Rating系统实现

以下是一个基于权重衰减模型的桌游Rating系统的Python实现，满足所有提出的要求：

```python
class GameRatingSystem:
    def __init__(self, T=10):
        """
        初始化桌游Rating系统

        参数:
            T (int): 权重计算中的常数参数，影响权重衰减速度
        """
        self.T = T                    # 权重计算常数
        self.match_results = []       # 存储比赛结果 (1 表示胜利, -1 表示失败)
        self.ratings_history = []     # 存储每次计算后的rating
        self.current_rating = 0       # 当前rating

    def add_match_result(self, won):
        """
        添加新的比赛结果

        参数:
            won (bool): 是否胜利

        返回:
            float: 更新后的rating
        """
        # 转换为1(胜利)或-1(失败)
        result = 1 if won else -1

        # 添加到比赛结果列表的开头(最新的在前)
        self.match_results.insert(0, result)

        # 重新计算rating
        self._update_rating()

        return self.current_rating

    def _update_rating(self):
        """
        根据所有历史比赛更新rating
        """
        N = len(self.match_results)
        total_rating = 0

        # 计算所有比赛结果的加权和
        for i in range(N):
            weight = 1 / (self.T + i + 1)  # 权重计算，i+1是因为i从0开始
            total_rating += weight * self.match_results[i]

        self.current_rating = total_rating
        self.ratings_history.append(total_rating)

    def get_current_rating(self):
        """
        获取当前rating

        返回:
            float: 当前rating
        """
        return self.current_rating

    def get_ratings_history(self):
        """
        获取历史rating变化

        返回:
            list: rating历史记录
        """
        return self.ratings_history

    def get_match_count(self):
        """
        获取已记录的比赛场次

        返回:
            int: 比赛场次
        """
        return len(self.match_results)

    def reset(self):
        """
        重置系统
        """
        self.match_results = []
        self.ratings_history = []
        self.current_rating = 0


# 使用示例
def demo_rating_system():
    # 创建rating系统实例
    rating_system = GameRatingSystem(T=5)

    # 模拟一系列比赛结果
    results = [True, False, True, True, False, True, False]

    print("模拟比赛结果 (最近的在前): ", ["胜" if r else "负" for r in results])
    print("\n比赛记录与Rating变化:")
    print("-" * 50)
    print("场次 | 结果 | Rating | 变化")
    print("-" * 50)

    prev_rating = 0
    for i, won in enumerate(results):
        new_rating = rating_system.add_match_result(won)
        change = new_rating - prev_rating
        change_symbol = "+" if change >= 0 else ""

        result_str = "胜利" if won else "失败"
        print(f"{i+1:2d}   | {result_str} | {new_rating:.4f} | {change_symbol}{change:.4f}")

        prev_rating = new_rating

    print("-" * 50)
    print(f"最终Rating: {rating_system.get_current_rating():.4f}")

    # 验证性质3: 最近一场胜利rating不降，最近一场失败rating不升
    print("\n验证Rating属性:")

    # 测试最近一场胜利
    test_system1 = GameRatingSystem(T=5)
    for r in results[1:]:  # 除去第一场
        test_system1.add_match_result(r)
    old_rating = test_system1.get_current_rating()
    test_system1.add_match_result(True)  # 添加一场胜利
    new_rating = test_system1.get_current_rating()
    print(f"添加一场胜利: {old_rating:.4f} -> {new_rating:.4f}, 变化: {new_rating - old_rating:.4f}")

    # 测试最近一场失败
    test_system2 = GameRatingSystem(T=5)
    for r in results[1:]:  # 除去第一场
        test_system2.add_match_result(r)
    old_rating = test_system2.get_current_rating()
    test_system2.add_match_result(False)  # 添加一场失败
    new_rating = test_system2.get_current_rating()
    print(f"添加一场失败: {old_rating:.4f} -> {new_rating:.4f}, 变化: {new_rating - old_rating:.4f}")


if __name__ == "__main__":
    demo_rating_system()
```

## 算法设计说明

1. **权重分配公式**:

   - 对第i场比赛(从最近开始计数)，权重为 wi = 1/(T+i+1)
   - T是可调参数，控制权重衰减速度
   - 满足 w1 ≥ w2 ≥ ... ≥ wN ≥ 0

2. **Rating计算**:

   - Rating = ∑(wi \* xi)，其中xi是比赛结果(1表示胜利，-1表示失败)
   - 每次添加新比赛后重新计算所有权重和总和

3. **满足的三个关键属性**:

   - ✓ 最近比赛的影响较大(权重随i增大而减小)
   - ✓ 早期比赛的影响较小(随时间推移，早期比赛的权重不断减小)
   - ✓ 最近一场胜利Rating不降；最近一场失败Rating不升(已在代码中验证)

4. **数据存储设计**:
   - 比赛结果倒序存储(最新的在前)，使权重计算更直观
   - 保存完整的Rating历史记录，便于分析和展示

这个实现不仅满足了所有要求，还易于扩展。例如，可以轻松加入新的功能如：根据对手强度调整胜负权重、添加季节性重置、计算信心区间等。

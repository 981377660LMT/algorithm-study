import numpy as np

def moe_routing_optimization(n, m, p, k, expert_probs):
    # step 1：将 n个专家平均分配到 m个NPU上，确定每个组的专家编号范围
    experts_per_npu = n // m
    groups = []
    for i in range(m):
        start = i * experts_per_npu
        end = (i + 1) * experts_per_npu if i != m - 1 else n
        groups.append((i, start, end))  # (NPU编号, 起始专家编号, 结束专家编号)

    # 打印组和编号
    print('每个组的专家编号范围:{}\n'.format(groups))
    
    # step 2：计算每个组的最大概率作为本组代表，选择概率最大的p个组
    group_max_probs = []
    for npu_id, start, end in groups:
        group_prob = expert_probs[start:end]
        max_group_prob = np.max(group_prob)
        group_max_probs.append((npu_id, max_group_prob))

    # 按组的最大概率降序排序，选择前p个组
    group_max_probs.sort(key=lambda x: x[1], reverse=True)
    selected_groups = group_max_probs[:p]

    # step 3：选择 p 个 NPU 对应的所有专家
    selected_experts = []
    for npu_id, _ in selected_groups:
        start, end = next((start, end) for _, start, end in groups if _ == npu_id)
        npu_experts_probs = expert_probs[start:end]
        # 获取专家编号（相对于全局的编号）
        expert_indices = np.arange(start, end)
        # 按概率降序排序
        sorted_indices = expert_indices[np.argsort(npu_experts_probs)[::-1]]
        sorted_probs = npu_experts_probs[np.argsort(npu_experts_probs)[::-1]]
        selected_experts.extend(list(zip(sorted_indices, sorted_probs)))

    # step 4: 从所有选中组的专家中，再选 k个概率最大的
    selected_experts.sort(key=lambda x: x[1], reverse=True)
    final_experts = selected_experts[:k]

    return final_experts

# 测试
n = 8  # 专家总数
m = 2  # NPU数量
p = 1  # 选择的组数量
k = 2  # 最终选择的专家数量
# 每个专家对应的概率，这里随机生成示例
expert_probs = np.random.rand(n)
print('随机生成每个专家的概率:\n{}\n\n'.format(expert_probs))

# 测试代码
final_routing = moe_routing_optimization(n, m, p, k, expert_probs)

print("最终路由的专家编号及概率：")
for expert_id, prob in final_routing:
    print(f"专家编号：{expert_id},概率:{prob:.4f}")
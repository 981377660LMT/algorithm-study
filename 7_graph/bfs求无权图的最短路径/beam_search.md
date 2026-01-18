Beam Search 可以通过多种方式进行优化，主要集中在**性能提升（速度/内存）**和**生成质量改善**两个方面。

以下是几个关键的优化方向及其代码实现示例：

### 1. 长度归一化 (Length Normalization)

原始分数是累加的对数概率（全为负数），序列越长，分数越低。这会导致模型倾向于生成极短的序列。我们需要引入长度惩罚来平衡。

### 2. 提前停止 (Early Stopping with EOS)

目前的实现是固定跑满 `max_steps`。实际上，如果生成的序列遇到了 `match <EOS>`（结束符），该序列就不应再继续扩展，而是直接移入完成列表。

### 3. Log-Sum-Exp 技巧 (数值稳定性)

如果在深度学习框架中使用，通常直接操作 logits，但在纯 Python 实现中，保持使用 Log Probability（对数概率）相加而不是概率相乘可以防止下溢。之前的代码已经使用了加法，这是对的。

### 4. 限制重复 (N-gram Blocking)

为了防止生成重复的短语（如 "I went to the the the..."），可以禁止已经出现过的 n-gram 再次生成。

---

### 优化后的 Beam Search 代码

这是一个集成了**长度归一化**和**EOS 处理**的改进版本：

```python
import math

def optimized_beam_search(predict_next_step, start_sequence, beam_width=3, max_steps=10, eos_token='<EOS>', alpha=0.7):
    """
    Args:
        alpha (float): 长度惩罚系数 (0.0 - 1.0)。0.6-0.8 是常见值。
                       分数计算变为: score / (len(seq) ** alpha)
    """
    # 存储最终完成的序列
    completed_sequences = []

    # 当前正在进行的 candidate: [(cumulative_score, sequence)]
    running_candidates = [[0.0, start_sequence]]

    for step in range(max_steps):
        if not running_candidates:
            break

        all_candidates = []

        for score, seq in running_candidates:
            # 获取下一步预测
            next_steps = predict_next_step(seq)

            for step_score, next_token in next_steps:
                new_score = score + step_score
                new_seq = seq + [next_token]

                # OPTIMIZATION 2: 处理结束符 EOS
                if next_token == eos_token:
                    # 序列结束，不需要继续扩展，直接加入完成列表
                    completed_sequences.append([new_score, new_seq])
                else:
                    # 序列未结束，加入候选列表等待下一轮扩展
                    all_candidates.append([new_score, new_seq])

        # 按照当前分数排序并截断 (为了下一轮的预测，暂时只看累积概率)
        # 注意：这里我们只保留 active 的序列，finished 序列已经在 completed_sequences 中
        ordered = sorted(all_candidates, key=lambda tup: tup[0], reverse=True)
        running_candidates = ordered[:beam_width]

    # 将剩余还在 running 的序列也视为完成（虽然是被迫截断的）
    completed_sequences.extend(running_candidates)

    # OPTIMIZATION 1: 长度归一化 (Length Normalization)
    # 最终排序时，依据归一化后的分数
    # 公式: score / (length^alpha)
    def normalize_score(candidate):
        score, seq = candidate
        # 减去起始符长度，避免除以0或偏差
        length = max(1, len(seq) - 1)
        return score / (length ** alpha)

    # 最终排序
    final_result = sorted(completed_sequences, key=normalize_score, reverse=True)

    # 返回 Top K
    return final_result[:beam_width]

# --- 验证优化 ---

def dummy_predict_advanced(seq):
    last_token = seq[-1]
    # 模拟：较短的路径分更高，但有了长度惩罚，长路径可能会胜出
    if last_token == '<BOS>':
        return [(-0.5, 'Quick'), (-2.0, 'Long')]
    elif last_token == 'Quick':
        return [(-0.1, '<EOS>')] # 总长2，总分 -0.6
    elif last_token == 'Long':
        return [(-0.5, 'Path')]
    elif last_token == 'Path':
        return [(-0.1, '<EOS>')] # 总长3，总分 -2.6
    return []

start_seq = ['<BOS>']

# 不带 alpha (alpha=0), 短序列优势巨大
print("--- Alpha = 0.0 (无惩罚) ---")
res_raw = optimized_beam_search(dummy_predict_advanced, start_seq, beam_width=2, alpha=0.0)
for s, seq in res_raw:
    print(f"Seq: {seq}, Raw Score: {s:.4f}")

# 带 alpha, 这里的逻辑会让长序列的分数除以更大的数（因为原分数为负数，除以大数变大接近0），
# 注意：对于负数 log 概率，归一化通常是 score / length。
# 因为 score 是负数，length 越大，结果绝对值越小（越接近0），排名越高。
print("\n--- Alpha = 1.0 (标准长度惩罚) ---")
res_norm = optimized_beam_search(dummy_predict_advanced, start_seq, beam_width=2, alpha=1.0)
for s, seq in res_norm:
    # 重新计算一下归一化分数展示
    norm_score = s / (len(seq)-1)**1.0
    print(f"Seq: {seq}, Norm Score: {norm_score:.4f} (Raw: {s})")
```

### 总结关键点

1.  **分离 `completed` 和 `running` 队列**：不要浪费计算资源去扩展已经生成 `<EOS>` 的序列。
2.  **后处理归一化**：不要在每一步循环中做长度归一化，只在最后排序时做。每一步扩展只应该基于概率（greedy assumption），否则会导致搜索方向偏离最可能的路径。
3.  **返回结构**：确保即使到达 `max_steps` 没有生成 `<EOS>`，这些长序列也能被返回。

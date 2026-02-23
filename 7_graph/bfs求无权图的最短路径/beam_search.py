def beam_search(predict_next_step, start_sequence, beam_width=3, max_steps=10):
    """
    实现 Beam Search 算法。

    Args:
        predict_next_step (func): 一个函数，接受当前序列，返回可能的下一步候选项及其对数概率 [(score, token), ...]。
        start_sequence (list): 初始序列 (例如 [<BOS>])。
        beam_width (int): 保留的候选序列数量 (K)。
        max_steps (int): 最大生成步数。

    Returns:
        list: 排序后的最终序列列表，格式为 [(cumulative_score, sequence), ...]。
    """

    sequences = [[0.0, start_sequence]]
    for _ in range(max_steps):
        all_candidates = []
        for score, seq in sequences:
            next_steps = predict_next_step(seq)
            for step_score, next_token in next_steps:
                candidate = [score + step_score, seq + [next_token]]
                all_candidates.append(candidate)
        ordered = sorted(all_candidates, key=lambda tup: tup[0], reverse=True)
        sequences = ordered[:beam_width]
    return sequences


# --- 示例用法 (模拟) ---
def dummy_predict(seq):
    """
    一个模拟的预测函数。
    实际上这里会调用你的神经网络或概率模型。
    """
    last_token = seq[-1]
    # 简单的模拟词汇表和概率
    if last_token == "<BOS>":
        return [(-0.1, "A"), (-0.3, "B"), (-1.0, "C")]
    elif last_token == "A":
        return [(-0.1, "B"), (-0.5, "C")]
    elif last_token == "B":
        return [(-0.2, "C"), (-0.4, "A")]
    else:  # 'C'
        return [(-0.1, "<EOS>")]


# 运行 Beam Search
start_seq = ["<BOS>"]
result = beam_search(dummy_predict, start_seq, beam_width=2, max_steps=3)

print("最终结果 (Top K):")
for score, seq in result:
    print(f"Score: {score:.4f}, Sequence: {seq}")

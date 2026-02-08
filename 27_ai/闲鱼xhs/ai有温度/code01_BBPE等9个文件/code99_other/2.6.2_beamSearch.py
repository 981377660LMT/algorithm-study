import torch

def beam_search(decoder, 
                sos_token, 
                eos_token, 
                vocab_size, 
                hidden_size, 
                k=3, 
                max_length=50):
    """
    Beam Search（束搜索）算法实现，用于序列生成任务
    与贪心搜索只保留最佳候选不同，束搜索保留top-k个最佳候选，平衡了性能和计算成本
    
    参数:
        decoder: 解码器模型，输入为(prev_words, hidden)，输出为(outputs, hidden)
        sos_token: 序列开始符号的ID
        eos_token: 序列结束符号的ID
        vocab_size: 词汇表大小
        hidden_size: 解码器隐藏层维度大小
        k: beam size（束大小），每次保留的候选序列数量
        max_length: 最大序列长度，防止生成无限长序列
    """
    # 初始化：每个beam序列都以SOS（开始符号）开头
    # 创建形状为(k, 1)的张量，填充sos_token
    k_prev_words = torch.full((k, 1), sos_token, dtype=torch.long)  # (k, 1)
    # 初始化序列集合，一开始每个序列都只包含SOS
    seqs = k_prev_words  # (k, 1)
    
    # 初始化每个候选序列的分数为0
    top_k_scores = torch.zeros(k, 1)
    
    # 存储已完成的序列（包含EOS的序列）及其对应的分数
    complete_seqs = []
    complete_seqs_scores = []
    
    # 记录当前解码步骤
    step = 1
    # 初始化解码器隐藏状态，形状为(1, k, hidden_size)
    # 第一个维度是LSTM层数*方向数，这里用1表示单层单向
    hidden = torch.zeros(1, k, hidden_size)
    
    # 解码循环，直到所有序列完成或达到最大长度
    while True:
        # 将上一步的词和隐藏状态输入解码器，得到输出和新的隐藏状态
        # outputs形状: (k, seq_len, vocab_size)
        # hidden形状: (1, k, hidden_size)
        outputs, hidden = decoder(k_prev_words, hidden)
        
        # 取最后一个时间步的输出，即当前步的预测结果
        # 形状变为(k, vocab_size)
        next_token_logits = outputs[:, -1, :]
        
        # 将输出转换为对数概率，避免数值下溢
        log_probs = torch.log_softmax(next_token_logits, dim=1)
        
        if step == 1:
            # 第一步特殊处理：初始时只有一个有效序列（全部是SOS）
            # 从第一个序列的预测结果中取top-k个词
            top_k_scores, top_k_words = log_probs[0].topk(k, dim=0, largest=True, sorted=True)
            # 所有候选都来自第0个初始序列
            prev_word_inds = torch.zeros(k, dtype=torch.long)
        else:
            # 后续步骤：需要合并历史分数和当前分数
            # 累加log概率（相当于概率相乘）
            cumulative_scores = top_k_scores + log_probs
            # 展平分数张量，以便全局选取top-k
            flat_scores = cumulative_scores.view(-1)
            # 选取全局top-k个最高分及其索引
            top_k_scores, top_k_indices = flat_scores.topk(k, 0, True, True)
            
            # 计算这些候选来自哪个beam（前一个词的索引）
            prev_word_inds = (top_k_indices // vocab_size).long()  # (k)
            # 计算当前预测的词索引
            top_k_words = (top_k_indices % vocab_size).long()      # (k)
        
        # 确保索引在有效范围内，防止越界错误
        prev_word_inds = torch.clamp(prev_word_inds, 0, seqs.size(0) - 1)
        
        # 更新序列：将前序序列与新预测的词拼接起来
        # 形状从(k, step)变为(k, step+1)
        seqs = torch.cat([seqs[prev_word_inds], top_k_words.unsqueeze(1)], dim=1)
        
        # 区分已完成和未完成的序列
        incomplete_inds = []  # 未完成序列的索引
        complete_inds = []    # 已完成序列的索引
        
        # 检查每个新预测的词是否是EOS
        for ind, next_word in enumerate(top_k_words):
            if next_word == eos_token:
                # 如果是EOS，标记为已完成序列
                complete_inds.append(ind)
            else:
                # 否则，继续保留为未完成序列
                incomplete_inds.append(ind)
        
        # 处理已完成的序列：加入到完成序列列表中
        if len(complete_inds) > 0:
            # 将完成的序列转换为列表并添加到集合中
            complete_seqs.extend(seqs[complete_inds].tolist())
            # 保存对应的分数
            complete_seqs_scores.extend(top_k_scores[complete_inds])
        
        # 更新beam size：减去已完成的序列数量
        k -= len(complete_inds)
        
        # 终止条件：所有序列都已完成，或达到最大长度
        if k == 0 or step >= max_length:
            break
        
        # 准备下一轮迭代的数据：只保留未完成的序列
        seqs = seqs[incomplete_inds]
        # 更新隐藏状态，只保留未完成序列对应的隐藏状态
        hidden = hidden[:, prev_word_inds[incomplete_inds], :]
        # 更新分数，只保留未完成序列的分数
        top_k_scores = top_k_scores[incomplete_inds].unsqueeze(1)
        # 更新下一轮的输入词
        k_prev_words = top_k_words[incomplete_inds].unsqueeze(1)
        
        # 步骤加1
        step += 1
    
    # 如果没有任何完成的序列（都没生成EOS）
    if not complete_seqs:
        if len(seqs) > 0:
            # 从剩余序列中选择分数最高的一个
            max_score_idx = torch.argmax(top_k_scores).item()
            return seqs[max_score_idx].tolist()
        else:
            # 极端情况：返回一个最小序列[SOS, EOS]
            return [sos_token, eos_token]
    
    # 从完成的序列中选择分数最高的序列作为结果
    max_score_idx = complete_seqs_scores.index(max(complete_seqs_scores))
    return complete_seqs[max_score_idx]


# 示例解码器类（用于测试beam search）
class SimpleDecoder:
    def __call__(self, input_ids, hidden):
        """
        简单的解码器实现，仅用于测试
        在实际应用中，这应该是真实的神经网络解码器
        """
        batch_size, seq_len = input_ids.shape
        # 简单起见，假设hidden_size等于vocab_size
        vocab_size = hidden.size(-1)
        
        # 生成随机输出（实际中应该是模型的预测结果）
        outputs = torch.randn(batch_size, seq_len, vocab_size)
        # 简单更新隐藏状态（实际中应该是LSTM/Transformer的输出）
        new_hidden = hidden + torch.randn_like(hidden) * 0.1
        
        return outputs, new_hidden


# 使用示例
if __name__ == "__main__":
    # 配置参数
    SOS_TOKEN = 0    # 开始符号ID
    EOS_TOKEN = 1    # 结束符号ID
    VOCAB_SIZE = 100 # 词汇表大小
    HIDDEN_SIZE = 128# 隐藏层大小
    BEAM_SIZE = 5    # beam大小
    MAX_LENGTH = 20  # 最大序列长度
    
    # 初始化解码器（实际应用中替换为真实模型）
    decoder = SimpleDecoder()
    
    # 执行beam search生成序列
    generated_sequence = beam_search(
        decoder=decoder,
        sos_token=SOS_TOKEN,
        eos_token=EOS_TOKEN,
        vocab_size=VOCAB_SIZE,
        hidden_size=HIDDEN_SIZE,
        k=BEAM_SIZE,
        max_length=MAX_LENGTH
    )
    
    print("生成的序列:", generated_sequence)
    
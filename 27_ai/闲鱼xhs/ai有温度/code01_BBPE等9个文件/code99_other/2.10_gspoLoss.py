def gspo_update(prompts, responses, rewards):
    # 计算token级log概率
    log_probs = policy.log_prob_per_token(prompts, responses)  # [batch, seq_len]
    ref_log_probs = ref_policy.log_prob_per_token(prompts, responses)  # [batch, seq_len]

    # 群组归一化优势（与GRPO相同）
    group_mean = rewards.mean()
    group_std = rewards.std() + 1e-8
    advantages = (rewards - group_mean) / group_std  # [batch]

    # GSPO关键：计算序列级重要性比率
    negative_approx_kl = log_probs - ref_log_probs  # [batch, seq_len]
    seq_lengths = response_mask.sum(dim=-1).clamp(min=1)  # [batch]

    # 序列级重要性比率：s_i(θ) = exp((1/|y_i|) * Σ_t log(π_θ/π_θ_old))
    negative_approx_kl_seq = (negative_approx_kl * response_mask).sum(dim=-1) / seq_lengths  # [batch]

    # 构造token级序列重要性比率（用于梯度计算）
    log_seq_ratio = (log_probs - log_probs.detach() +
                     negative_approx_kl_seq.detach().unsqueeze(-1))  # [batch, seq_len]
    seq_importance_ratio = torch.exp(log_seq_ratio)  # [batch, seq_len]

    # PPO裁剪损失
    clipped_ratio = torch.clamp(seq_importance_ratio, 1-epsilon_low, 1+epsilon_high)
    loss1 = -advantages.unsqueeze(-1) * seq_importance_ratio
    loss2 = -advantages.unsqueeze(-1) * clipped_ratio
    pg_losses = torch.maximum(loss1, loss2)  # 取最大值而非最小值

    # 序列级聚合损失
    policy_loss = (pg_losses * response_mask).sum() / response_mask.sum()

    # 更新策略
    optimizer.zero_grad()
    policy_loss.backward()
    optimizer.step()

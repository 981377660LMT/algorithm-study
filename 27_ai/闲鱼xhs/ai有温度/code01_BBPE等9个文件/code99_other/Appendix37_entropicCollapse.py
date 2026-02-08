import torch  # 假设使用PyTorch框架

# Clip-Higher策略参数
clip_higher_args = {
    "select_ratio": 0.2,        # 选择20%的token进行特殊处理
    "method": "clip_higher",    # 指定策略类型
    "epsilon_high": 0.28,       # 低概率token提升上限（1+28%）
    "epsilon_low": 0.2,         # 高概率token降低下限（1-20%）
    "cov_control": True         # 启用高协方差token梯度控制
}

# KL-Cov策略参数
kl_cov_args = {
    "select_ratio": 0.2,        # 选择20%的高协方差token
    "method": "kl_cov",         # 指定策略类型
    "kl_coef": 0.1              # KL惩罚系数
}


def compute_policy_loss(old_log_prob, log_prob, advantages, select_ratio, method, **args):
    """
    计算策略损失，支持Clip-Higher和KL-Cov两种策略
    
    参数:
        old_log_prob: 旧策略的对数概率 (tensor)
        log_prob: 当前策略的对数概率 (tensor)
        advantages: 优势函数值 (tensor)
        select_ratio: 筛选token的比例
        method: 策略类型 ("clip_higher" 或 "kl_cov")
       ** args: 策略专用参数
    """
    # 计算新旧策略的概率比
    ratio = torch.exp(log_prob - old_log_prob)
    # 基础策略损失（未裁剪/未加惩罚）
    pg_losses1 = -ratio * advantages

    # 计算token级别的中心化协方差
    covs = (log_prob - log_prob.mean()) * (advantages - advantages.mean())
    select_num = int(select_ratio * len(pg_losses1))

    if method == "clip_higher":
        # 获取非对称裁剪参数
        epsilon_high = args["epsilon_high"]
        epsilon_low = args["epsilon_low"]
        
        # 计算裁剪上下界（非对称）
        clip_lb = 1 - epsilon_low  # 高概率token的降低限制
        clip_ub = 1 + epsilon_high  # 低概率token的提升限制
        
        # 应用裁剪的损失
        pg_losses2 = -torch.clamp(ratio, clip_lb, clip_ub) * advantages
        
        # 对高协方差token进行梯度控制（可选）
        if args.get("cov_control", False):
            # 筛选高协方差token的索引
            high_cov_idx = torch.topk(covs, k=select_num, largest=True).indices
            # 截断这些token的梯度
            pg_losses1[high_cov_idx] = pg_losses1[high_cov_idx].detach()
            pg_losses2[high_cov_idx] = pg_losses2[high_cov_idx].detach()
        
        # 取裁剪前后损失的最大值（PPO风格更新）
        pg_loss = torch.max(pg_losses1, pg_losses2).mean()

    elif method == "kl_cov":
        # 获取KL惩罚系数
        kl_coef = args["kl_coef"]
        # 计算KL散度相关惩罚项（用绝对值近似）
        kl_penalty = torch.abs(log_prob - old_log_prob)
        # 筛选高协方差token
        select_idx = torch.topk(covs, k=select_num, largest=True).indices
        # 仅对高协方差token添加KL惩罚
        pg_losses1[select_idx] += kl_coef * kl_penalty[select_idx]
        pg_loss = pg_losses1.mean()

    else:
        raise ValueError(f"不支持的策略方法: {method}")

    return pg_loss

# 调用示例
if __name__ == "__main__":
    # 生成示例数据（模拟策略输出）
    batch_size = 128
    old_log_prob = torch.randn(batch_size)  # 旧策略对数概率
    log_prob = torch.randn(batch_size)     # 当前策略对数概率
    advantages = torch.randn(batch_size)   # 优势函数值

    # 1. 使用Clip-Higher策略计算损失
    loss_clip_higher = compute_policy_loss(
        old_log_prob=old_log_prob,
        log_prob=log_prob,
        advantages=advantages,
        **clip_higher_args
    )
    print(f"Clip-Higher策略损失: {loss_clip_higher.item()}")

    # 2. 使用KL-Cov策略计算损失
    loss_kl_cov = compute_policy_loss(
        old_log_prob=old_log_prob,
        log_prob=log_prob,
        advantages=advantages,** kl_cov_args
    )
    print(f"KL-Cov策略损失: {loss_kl_cov.item()}")

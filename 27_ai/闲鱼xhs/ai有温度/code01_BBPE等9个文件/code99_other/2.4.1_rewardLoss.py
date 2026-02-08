import torch
import torch.nn as nn
import torch.nn.functional as F

class PairWiseLoss(nn.Module):
    """
    Pairwise Loss 实现 Reward Model 的成对损失计算
    核心逻辑：最大化 chosen response（优响应）与 rejected response（劣响应）的奖励差值
    公式对应：Loss = -log(sigmoid(r_c - r_r - margin)) （有margin时）或 -log(sigmoid(r_c - r_r)) （无margin时）
    
    参数说明：
    - margin: 可选的边际值，用于拉开优/劣响应的奖励差距
    前向传播参数：
    - chosen_reward: 优响应的奖励值，shape: [batch_size] 或 [batch_size, seq_len]
    - reject_reward: 劣响应的奖励值，shape 需与 chosen_reward 一致
    - margin: 边际值（优先级：forward传参 > 初始化传参）
    """
    def __init__(self, margin: float = None):
        super().__init__()
        self.margin = margin

    def forward(self, chosen_reward: torch.Tensor, reject_reward: torch.Tensor, margin: float = None):
        current_margin = margin if margin is not None else self.margin
        reward_diff = chosen_reward - reject_reward
        if current_margin is not None:
            loss = -F.logsigmoid(reward_diff - current_margin)
        else:
            loss = -F.logsigmoid(reward_diff)
        return loss.mean()


if __name__ == "__main__":
    # 实例化损失函数（两种方式）
    loss_fn = PairWiseLoss()
    chosen_rewards = torch.tensor([1.2, 0.8, 1.5, 0.9])
    reject_rewards = torch.tensor([0.5, 0.3, 0.7, 0.2])
    loss_no_margin = loss_fn(chosen_rewards, reject_rewards)
    print(f"无margin的损失值: {loss_no_margin.item():.4f}")

    # 计算损失（有margin，值为0.1）
    loss_with_margin = loss_fn(chosen_rewards, reject_rewards, margin=0.1)
    print(f"有margin(0.1)的损失值: {loss_with_margin.item():.4f}")
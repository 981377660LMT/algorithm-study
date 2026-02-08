import torch
import torch.nn as nn
from collections import *
import torch.functional as F

class PolicyLoss(nn.Module):
    """
    PPO（Proximal Policy Optimization）算法中的策略损失计算
    
    实现了PPO的剪辑损失（clipped surrogate objective），用于限制策略更新的幅度，
    提高训练的稳定性。通过比较新策略与旧策略的比例，取两种可能损失的最小值。
    """

    def __init__(self, clip_eps: float = 0.2) -> None:
        """
        初始化策略损失计算器
        
        参数:
            clip_eps: 剪辑范围的epsilon值，控制策略更新的最大幅度，默认0.2
        """
        super().__init__()
        self.clip_eps = clip_eps

    def forward(
        self,
        log_probs: torch.Tensor,      # 新策略下的动作对数概率
        old_log_probs: torch.Tensor,  # 旧策略下的动作对数概率
        advantages: torch.Tensor,     # 优势函数值，衡量动作的相对好坏
        action_mask: Optional[torch.Tensor] = None,  # 动作掩码，用于过滤无效动作
    ) -> torch.Tensor:
        """
        计算PPO的策略损失
        
        步骤:
        1. 计算新旧策略的概率比（ratio）
        2. 计算两种可能的替代损失（surrogate loss）
        3. 取两种损失的最小值作为最终损失
        4. 应用掩码并计算平均值
        
        返回:
            计算得到的策略损失值
        """
        # 计算新旧策略的概率比：exp(new_log_prob - old_log_prob) = new_prob / old_prob
        ratio = (log_probs - old_log_probs).exp()
        
        # 第一种替代损失：直接使用概率比乘以优势值
        surr1 = ratio * advantages
        
        # 第二种替代损失：将概率比限制在[1-ε, 1+ε]范围内，再乘以优势值
        surr2 = ratio.clamp(1 - self.clip_eps, 1 + self.clip_eps) * advantages
        
        # PPO的剪辑损失：取两种替代损失的最小值，并取负（因为是最大化问题，转为最小化）
        loss = -torch.min(surr1, surr2)
        
        # 应用掩码（如果有），计算指定维度的平均值，最后再求总体平均
        loss = masked_mean(loss, action_mask, dim=-1).mean()
        return loss

def masked_mean(tensor, mask, dim):
    """
    计算带掩码的张量平均值
    
    对于被掩码标记为无效的元素，在计算平均值时会被忽略。
    常用于处理序列数据中长度不一致的情况（如padding部分）。
    
    参数:
        tensor: 输入张量
        mask: 掩码张量，与tensor形状相同，有效元素为1，无效为0
        dim: 需要计算平均值的维度
        
    返回:
        按指定维度计算的带掩码的平均值
    """
    if mask is None:
        # 如果没有掩码，直接计算指定维度的平均值
        return tensor.mean(axis=dim)
    # 计算掩码区域内的元素和，再除以有效元素的数量
    return (tensor * mask).sum(axis=dim) / mask.sum(axis=dim)

class ValueLoss(nn.Module):
    """
    PPO算法中的价值函数损失计算
    
    用于训练价值网络（Critic），价值网络用于估计状态的预期收益。
    支持对价值更新进行剪辑，类似于策略损失的剪辑机制。
    """

    def __init__(self, clip_eps: float = None) -> None:
        """
        初始化价值损失计算器
        
        参数:
            clip_eps: 剪辑范围的epsilon值，若为None则不进行剪辑
        """
        super().__init__()
        self.clip_eps = clip_eps

    def forward(
        self,
        values: torch.Tensor,         # 当前价值网络输出的状态价值
        old_values: torch.Tensor,     # 旧价值网络输出的状态价值
        returns: torch.Tensor,        # 实际的累积回报（目标值）
        action_mask: Optional[torch.Tensor] = None,  # 掩码，用于过滤无效位置
    ) -> torch.Tensor:
        """
        计算价值函数损失
        
        步骤:
        1. 如果启用剪辑，计算剪辑后的价值和两种可能的平方损失
        2. 取两种损失的最大值作为最终损失（剪辑机制）
        3. 如果不启用剪辑，直接计算MSE损失
        4. 应用掩码并计算平均值，乘以0.5是为了求导方便
        
        返回:
            计算得到的价值损失值
        """
        # 对价值函数也应用剪辑操作，原理与策略损失类似
        if self.clip_eps is not None:
            # 计算剪辑后的价值：限制新价值与旧价值的差异在[-ε, ε]范围内
            values_clipped = old_values + (values - old_values).clamp(-self.clip_eps, self.clip_eps)
            
            # 计算两种平方损失：剪辑后价值与回报的平方差、原始价值与回报的平方差
            surr1 = (values_clipped - returns) ** 2
            surr2 = (values - returns) ** 2
            
            # 取两种损失的最大值作为价值损失（确保不更新幅度过大）
            loss = torch.max(surr1, surr2)
        else:
            # 不使用剪辑时，直接计算MSE损失
            loss = (values - returns) ** 2

        # 应用掩码，计算指定维度的平均值，最后乘以0.5（MSE损失的常规做法）
        loss = masked_mean(loss, action_mask, dim=-1).mean()
        return 0.5 * loss
    
class PairWiseLoss(nn.Module):
    """
    奖励模型（Reward Model）中的成对损失计算
    
    用于训练奖励模型，通过比较两个候选响应（chosen和reject）的奖励值，
    使模型学会区分好坏响应。常见于强化学习从人类反馈中学习（RLHF）流程。
    """
    def forward(self, chosen_reward, reject_reward, margin):
        """
        计算成对损失
        
        基于log-sigmoid函数，鼓励chosen响应的奖励高于reject响应的奖励。
        可加入margin参数增加区分难度，要求chosen奖励至少比reject高margin。
        
        参数:
            chosen_reward: 被选中的好响应的奖励值
            reject_reward: 被拒绝的差响应的奖励值
            margin:  margin值，若不为None，则要求奖励差至少为margin
            
        返回:
            计算得到的成对损失的平均值
        """
        if margin is not None:
            # 带margin的损失计算：鼓励chosen_reward - reject_reward > margin
            loss = -F.logsigmoid(chosen_reward - reject_reward - margin)
        else:
            # 不带margin的损失计算：仅要求chosen_reward > reject_reward
            loss = -F.logsigmoid(chosen_reward - reject_reward)
        # 返回损失的平均值
        return loss.mean()
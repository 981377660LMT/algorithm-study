import torch
import torch.nn as nn
import torch.nn.functional as F
from typing import Tuple

class DPOLoss(nn.Module):
    """
    DPO（Direct Preference Optimization）损失函数
    
    直接偏好优化算法，用于训练语言模型以符合人类偏好。
    无需显式中间奖励模型，直接通过偏好数据（选中的响应vs被拒绝的响应）优化策略。
    """

    def __init__(self, beta: float, label_smoothing: float = 0.0, ipo: bool = False) -> None:
        """
        初始化DPO损失计算器
        
        参数:
            beta: 温度参数，控制策略更新的强度
            label_smoothing: 标签平滑系数，用于缓解过拟合，默认0.0（不使用平滑）
            ipo: 是否使用IPO（Implicit Preference Optimization）变体，默认False
        """
        super().__init__()
        self.beta = beta
        self.label_smoothing = label_smoothing
        self.ipo = ipo

    def forward(
        self,
        policy_chosen_logps: torch.Tensor,      # 策略模型对选中响应的对数概率
        policy_rejected_logps: torch.Tensor,    # 策略模型对被拒绝响应的对数概率
        reference_chosen_logps: torch.Tensor,   # 参考模型对选中响应的对数概率
        reference_rejected_logps: torch.Tensor, # 参考模型对被拒绝响应的对数概率
    ) -> Tuple[torch.Tensor, torch.Tensor, torch.Tensor]:
        """
        计算DPO损失
        
        步骤:
        1. 计算策略模型和参考模型对选中/拒绝响应的对数概率比
        2. 计算两者的差值作为logits
        3. 根据是否使用IPO计算不同的损失函数
        4. 计算选中和拒绝响应的奖励估计值
        
        返回:
            损失值、选中响应的奖励、被拒绝响应的奖励
        """
        # 计算策略模型下选中响应与被拒绝响应的对数概率比
        pi_logratios = policy_chosen_logps - policy_rejected_logps
        # 计算参考模型下选中响应与被拒绝响应的对数概率比
        ref_logratios = reference_chosen_logps - reference_rejected_logps
        # 计算logits：策略比与参考比的差值，衡量策略相对参考模型的改进
        logits = pi_logratios - ref_logratios

        if self.ipo:
            # 使用IPO损失（平方损失），来自论文https://arxiv.org/pdf/2310.12036v2.pdf公式17
            # 鼓励logits接近1/(2*beta)，实现更稳定的优化
            losses = (logits - 1 / (2 * self.beta)) ** 2  
        else:
            # 标准DPO损失，来自论文https://ericmitchell.ai/cdpo.pdf公式3
            # 当label_smoothing=0时，等价于原始DPO公式（https://arxiv.org/pdf/2305.18290.pdf公式7）
            # 结合标签平滑，平衡选中和拒绝样本的损失权重
            losses = (
                -F.logsigmoid(self.beta * logits) * (1 - self.label_smoothing)  # 选中样本的损失分量
                - F.logsigmoid(-self.beta * logits) * self.label_smoothing      # 拒绝样本的损失分量（平滑项）
            )

        # 计算批次平均损失
        loss = losses.mean()
        
        # 计算选中和拒绝响应的奖励估计（基于策略与参考模型的差异）
        # detach()确保不影响梯度计算
        chosen_rewards = self.beta * (policy_chosen_logps - reference_chosen_logps).detach()
        rejected_rewards = self.beta * (policy_rejected_logps - reference_rejected_logps).detach()

        return loss, chosen_rewards, rejected_rewards
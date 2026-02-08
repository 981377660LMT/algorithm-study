import torch
import torch.nn as nn

class RMSNorm(nn.Module):
    def __init__(self, hidden_size, eps=1e-6):
        super().__init__()
        self.weight = nn.Parameter(torch.ones(hidden_size))  # 缩放参数
        self.eps = eps

    def forward(self, x):
        # 计算平方的均值（RMS核心）
        mean_square = x.pow(2).mean(-1, keepdim=True)
        # 归一化并应用缩放
        return self.weight * x / torch.sqrt(mean_square + self.eps)
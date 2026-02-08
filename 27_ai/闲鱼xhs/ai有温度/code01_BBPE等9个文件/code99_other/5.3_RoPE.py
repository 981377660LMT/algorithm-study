import torch
import matplotlib.pyplot as plt
import numpy as np

def compute_rope(x, theta_base=10000):
    """
    整合参数计算的RoPE实现：直接在旋转函数中计算正弦余弦参数
    
    参数:
        x: 输入张量，形状为(batch_size, num_heads, seq_len, head_dim)
        theta_base: 用于计算频率的基础值（原论文中使用10000）
    
    返回:
        x_rotated: 应用旋转后的张量，形状与x相同
    """
    # 获取输入张量的维度信息
    batch_size, num_heads, seq_len, head_dim = x.shape
    assert head_dim % 2 == 0, "Head dimension must be even"

    # 计算逆频率：1 / (theta_base ^ (2i / head_dim))
    inv_freq = 1.0 / (theta_base ** (torch.arange(0, head_dim, 2).float() / head_dim))
    
    # 生成位置索引（0到seq_len-1，根据输入序列长度动态计算）
    positions = torch.arange(seq_len)
    
    # 计算角度：position * inv_freq
    angles = positions[:, None] * inv_freq[None, :]  # 形状: (seq_len, head_dim // 2)
    
    # 扩展角度到与head_dim匹配
    angles = torch.cat([angles, angles], dim=1)  # 形状: (seq_len, head_dim)
    
    # 计算正弦和余弦参数
    cos = torch.cos(angles)
    sin = torch.sin(angles)

    # 将输入向量按维度分成两半
    x1 = x[..., : head_dim // 2]  # 前半部分维度
    x2 = x[..., head_dim // 2 :]  # 后半部分维度

    # 调整cos和sin的形状，添加batch和head维度以匹配输入
    cos = cos.unsqueeze(0).unsqueeze(0)  # 形状: (1, 1, seq_len, head_dim)
    sin = sin.unsqueeze(0).unsqueeze(0)

    # 构造旋转项并应用旋转公式
    rotated = torch.cat((-x2, x1), dim=-1)
    x_rotated = (x * cos) + (rotated * sin)

    return x_rotated.to(dtype=x.dtype)
import torch
import torch.nn.functional as F

def temperature_sampling(logits: torch.Tensor, temperature: float = 1.0) -> torch.Tensor:
    """
    简化的温度采样函数，专注于核心的温度缩放逻辑
    假设输入为2D张量(batch_size, vocab_size)，温度参数>0
    
    参数:
        logits: 模型输出的原始logits，形状为 (batch_size, vocab_size) 或 (vocab_size)
        temperature: 温度参数（τ > 0）
                     - τ < 1: 降低随机性，增强确定性（分布更陡峭）
                     - τ = 1: 保持原始概率分布
                     - τ > 1: 增加随机性，分布更平缓
    
    返回:
        采样得到的token索引，形状为(batch_size,)
    """
    # 温度缩放logits，添加微小值防止数值不稳定
    scaled_logits = logits / (temperature + 1e-8)
    
    # 转换为概率分布
    probabilities = F.softmax(scaled_logits, dim=-1)
    
    # 基于概率分布采样
    sampled_indices = torch.multinomial(probabilities, num_samples=1).squeeze(1)
    
    return sampled_indices


# 使用示例
if __name__ == "__main__":
    # 模拟模型输出的logits (batch_size=2, vocab_size=10)
    torch.manual_seed(42)
    logits = torch.randn(2, 10)
    
    # 不同温度下的采样结果
    for temp in [0.1, 1.0, 5.0]:
        tokens = temperature_sampling(logits, temp)
        print(f"温度={temp}时采样结果: {tokens.tolist()}")
        
        # 展示概率分布变化
        probs = F.softmax(logits / (temp + 1e-8), dim=-1)
        print(f"  最大概率值: {probs.max(dim=1).values.tolist()}\n")

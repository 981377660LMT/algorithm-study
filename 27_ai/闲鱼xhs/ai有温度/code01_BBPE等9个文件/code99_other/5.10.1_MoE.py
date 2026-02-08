import torch
import torch.nn as nn
import torch.nn.functional as F

# 基础专家模块定义
class BasicExpert(nn.Module):
    def __init__(self, in_dim, out_dim):
        super().__init__()
        # 专家模块本质上是一个线性变换层
        self.fc = nn.Linear(in_dim, out_dim)
        
    def forward(self, x):
        # 对输入进行线性变换并返回结果
        return self.fc(x)

class MOERouter(nn.Module):
    def __init__(self, hidden_dim, expert_number, top_k):
        super().__init__()
        # 路由网络：一个线性层，用于计算每个专家的权重
        self.gate = nn.Linear(hidden_dim, expert_number)
        self.expert_number = expert_number  # 专家总数
        self.top_k = top_k  # 每个token选择的专家数量
    
    def forward(self, hidden_states):
        # 计算路由logits：每个token对每个专家的原始分数
        # hidden_states形状: (b * s, hidden_dim)，输出形状: (b * s, expert_number)
        router_logits = self.gate(hidden_states)
        
        # 计算专家经过softmax之后的概率分布
        routing_probs = F.softmax(router_logits, dim=-1, dtype=torch.float)
        
        # 选择概率最高的top_k个专家
        # router_weights: 选中专家的概率值，形状: (b * s, top_k)
        # selected_experts: 选中专家的索引，形状: (b * s, top_k)
        router_weights, selected_experts = torch.topk(
            routing_probs, self.top_k, dim=-1
        )
        
        # 专家权重归一化：确保选中的top_k个专家权重之和为1
        router_weights = router_weights / router_weights.sum(dim=-1, keepdim=True)
        # 转换权重数据类型以匹配隐藏状态
        router_weights = router_weights.to(hidden_states.dtype)
        
        # 生成专家掩码：标记每个token被分配给哪些专家
        # 先转为one-hot编码，形状: (b * s, top_k, expert_number)
        expert_mask = F.one_hot(
            selected_experts,
            num_classes=self.expert_number
        )
        # 调整维度顺序，方便后续处理: (expert_number, top_k, b * s)
        expert_mask = expert_mask.permute(2, 1, 0)
        
        return router_logits, router_weights, selected_experts, expert_mask


class SparseMOE(nn.Module):
    # 稀疏MOE模型：每个token仅经过top_k个专家，而不是所有专家
    def __init__(self, hidden_dim, expert_number, top_k, shared_experts_number=2):
        super().__init__()

        self.hidden_dim = hidden_dim  # 隐藏层维度
        self.expert_number = expert_number  # 专家总数
        self.top_k = top_k  # 每个token选择的专家数量
        self.shared_experts_number = shared_experts_number  # 共享专家数量（当前实现未使用）

        # 创建专家列表
        self.experts = nn.ModuleList(
            [
                BasicExpert(self.hidden_dim, self.hidden_dim) 
                for _ in range(self.expert_number)
            ]
        )

        # 创建路由模块
        self.router = MOERouter(self.hidden_dim, self.expert_number, self.top_k)
    
    def forward(self, x):
        # x输入形状: (batch_size, seq_len, hidden_dim)
        batch_size, seq_len, hidden_dim = x.size()

        # 合并batch和sequence维度，转为token级别处理
        # 新形状: (batch_size * seq_len, hidden_dim)
        hidden_states = x.view(-1, hidden_dim)

        # 通过路由模块获取路由信息
        # router_logits: 每个token对每个专家的原始分数
        # router_weights: 选中专家的归一化权重
        # selected_experts_indices: 选中专家的索引
        # expert_mask: 标记每个token被分配给哪些专家的掩码
        router_logits, router_weights, selected_experts_indices, expert_mask = self.router(hidden_states)
        
        # 初始化最终输出张量
        final_hidden_states = torch.zeros(
            (batch_size * seq_len, hidden_dim),
            dtype=hidden_states.dtype,
            device=hidden_states.device
        )

        # 逐个专家处理其负责的token
        for expert_idx in range(self.expert_number):
            # 获取当前专家
            expert_layer = self.experts[expert_idx]
            # expert_mask[expert_idx]形状: (top_k, batch_size * seq_len)
            
            # 找到当前专家需要处理的token
            # idx: 表示是该token的第几个选中专家(0或1，对应top1或top2)
            # top_x: 选中的token在(batch*seq_len)维度上的索引
            idx, top_x = torch.where(expert_mask[expert_idx])
            
            # 获取这些token的隐藏状态
            current_state = hidden_states.unsqueeze(0)[:, top_x, :].reshape(-1, hidden_dim)
            
            # 应用专家处理并乘以对应的权重
            # router_weights[top_x, idx]获取每个token对应于当前专家的权重
            current_hidden_states = expert_layer(current_state) * router_weights[top_x, idx].unsqueeze(-1)
            
            # 将当前专家的输出累加到最终结果中
            final_hidden_states.index_add_(0, top_x, current_hidden_states.to(hidden_states.dtype))

        # 将结果恢复为原始的(batch, seq_len, hidden_dim)形状
        final_hidden_states = final_hidden_states.reshape(batch_size, seq_len, hidden_dim)

        return final_hidden_states, router_logits  # 输出最终隐藏状态和路由logits


def test_moe():
    # 创建测试输入: batch_size=2, seq_len=4, hidden_dim=16
    x = torch.rand(2, 4, 16)
    # 初始化MOE模型: 16维隐藏层, 4个专家, 每个token选2个专家
    token_level_moe = SparseMOE(hidden_dim=16, expert_number=4, top_k=2)
    # 前向传播
    out = token_level_moe(x)
    # 打印输出形状，验证是否符合预期
    # 预期输出: torch.Size([2, 4, 16]) torch.Size([8, 4])
    print(out[0].shape, out[1].shape)


# 执行测试
test_moe()
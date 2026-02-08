class ShareExpertMOE(nn.Module):
    """
    共享专家混合模型
    结合了稀疏混合专家(SparseMOE)和共享专家(shared experts)的结构
    共享专家对所有输入进行处理，而稀疏MOE则通过路由选择部分专家处理
    """
    def __init__(self, hidden_dim, expert_number, top_k, shared_experts_number=2):
        super().__init__()
        
        # 初始化稀疏混合专家模型
        # hidden_dim: 输入特征维度
        # expert_number: 稀疏MOE中的专家数量
        # top_k: 每个输入选择的专家数量
        self.moe_model = SparseMOE(hidden_dim, expert_number, top_k)
        
        # 创建共享专家列表
        # 共享专家对所有输入进行处理，而不是通过路由选择
        self.shared_experts = nn.ModuleList(
            [
                BasicExpert(
                    hidden_dim, hidden_dim  # 输入输出维度均为hidden_dim
                ) for _ in range(shared_experts_number)  # 创建指定数量的共享专家
            ]
        )

    def forward(self, x):
        # x shape 是 (b, s, hidden_dim)，其中：
        # b: batch size（批次大小）
        # s: sequence length（序列长度）
        # hidden_dim: 特征维度
        
        # 首先通过稀疏MOE模型处理输入
        # 返回值：稀疏专家的输出和路由logits（用于选择专家的概率分布）
        sparse_moe_out, router_logits = self.moe_model(x)
        
        # 然后通过所有共享专家处理原始输入x
        # 每个专家的输出形状保持为 (b, s, hidden_dim)
        shared_experts_out = [
            expert(x) for expert in self.shared_experts
        ]
        
        # 将所有共享专家的输出在新维度上堆叠，然后求和
        # 堆叠后形状: (shared_experts_number, b, s, hidden_dim)
        # 求和后形状: (b, s, hidden_dim)，与稀疏MOE输出形状一致
        shared_experts_out = torch.stack(
            shared_experts_out, dim=0
        ).sum(dim=0)
        
        # 最终输出为稀疏MOE输出与共享专家输出的总和，以及路由logits
        return sparse_moe_out + shared_experts_out, router_logits


def test_share_expert_moe():
    """测试ShareExpertMOE模型的功能和输出形状"""
    # 创建随机输入张量: (batch_size=2, sequence_length=4, hidden_dim=16)
    x = torch.rand(2, 4, 16)
    
    # 初始化模型: 隐藏维度16，4个稀疏专家，每个输入选2个专家，2个共享专家
    share_expert_moe = ShareExpertMOE(hidden_dim=16, expert_number=4, top_k=2, shared_experts_number=2)
    
    # 前向传播获取输出
    out = share_expert_moe(x)
    
    # 打印输出形状，验证是否符合预期
    # 预期输出:
    # 模型输出张量形状: (2, 4, 16)，与输入形状一致
    # 路由logits形状: (2, 4, 4)，4对应稀疏专家数量
    print(out[0].shape, out[1].shape)


# 执行测试函数
test_share_expert_moe()
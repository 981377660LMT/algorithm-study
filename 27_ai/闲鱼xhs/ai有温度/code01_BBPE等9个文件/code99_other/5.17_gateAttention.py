import torch
import torch.nn as nn

class GatedAttentionHead(nn.Module):
    def __init__(self, d_model, n_heads, d_head, gate_hidden_dim=None):
        super().__init__()
        self.n_heads = n_heads
        self.d_head = d_head
        gate_hidden_dim = gate_hidden_dim or d_model
        
        # 用 pre-norm hidden state 生成 gate head-specific + element-wise
        self.gate_mlp = nn.Sequential(
            nn.Linear(d_model, gate_hidden_dim),
            nn.SiLU(),
            nn.Linear(gate_hidden_dim, n_heads * d_head),
        )
        self.ln = nn.LayerNorm(d_model)

    def forward(self, x, attn_out):
        """
        x: [B, L, d_model] pre-norm hidden
        attn_out: [B, n_heads, L, d_head] SDPA 输出
        """
        B, H, L, Dh = attn_out.shape
        assert H == self.n_heads and Dh == self.d_head
        
        gate_input = self.ln(x)                         # [B, L, d_model]
        gate = self.gate_mlp(gate_input)                # [B, L, H * Dh]
        gate = gate.view(B, L, H, Dh)                   # [B, L, H, Dh]
        gate = torch.sigmoid(gate)                      # [0,1] 的软 mask
        gate = gate.permute(0, 2, 1, 3)                 # [B, H, L, Dh], 对齐 attn_out
        gated_out = attn_out * gate                     # element-wise gating
        return gated_out

if __name__ == "__main__":
    # --- 超参数设置 ---
    batch_size = 2
    seq_len = 10
    d_model = 64
    n_heads = 4
    d_head = 16  # 注意：通常 n_heads * d_head == d_model，但这里为了测试独立性可以不等
    
    # --- 实例化模型 ---
    gated_attn = GatedAttentionHead(
        d_model=d_model, 
        n_heads=n_heads, 
        d_head=d_head
    )
    
    # --- 创建虚拟数据 ---
    # x: [Batch, Length, d_model]
    x = torch.randn(batch_size, seq_len, d_model)
    
    # attn_out: [Batch, n_heads, Length, d_head] (标准 SDPA 输出形状)
    attn_out = torch.randn(batch_size, n_heads, seq_len, d_head)
    
    print(f"输入 x 形状: {x.shape}")
    print(f"输入 attn_out 形状: {attn_out.shape}")
    
    # --- 前向传播 ---
    output = gated_attn(x, attn_out)
    
    print("-" * 30)
    print(f"输出 gated_out 形状: {output.shape}")
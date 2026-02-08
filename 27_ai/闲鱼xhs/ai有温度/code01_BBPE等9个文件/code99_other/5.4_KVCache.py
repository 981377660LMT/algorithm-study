import torch
import torch.nn as nn
import math
import torch.nn.functional as F

class Attention(nn.Module):
    """
    注意力机制模块，实现了带有KV缓存的自注意力机制
    用于在生成式模型中高效计算注意力，避免重复计算历史序列的键值对
    """
    def __init__(self, args):
        super().__init__()
        # 调用父类初始化方法
        super(Attention, self).__init__()
        
        # 此处省略了和 KV cache 无关的初始化代码
        # 例如可能包括：
        # - 定义查询、键、值的线性变换层 (wq, wk, wv)
        # - 定义输出线性变换层 (wo)
        # - 设置注意力头数、头维度等参数 (n_local_heads, head_dim)
        
        # 初始化KV缓存，用于存储历史序列的键和值
        # 缓存形状: (最大批次大小, 最大序列长度, 本地头数, 头维度)
        # .cuda() 表示将缓存存储在GPU上以加速计算
        self.cache_k = torch.zeros(
            (args.max_batch_size, args.max_seq_len, self.n_local_heads, self.head_dim)
        ).cuda()
        self.cache_v = torch.zeros(
            (args.max_batch_size, args.max_seq_len, self.n_local_heads, self.head_dim)
        ).cuda()

    def forward(self, x, start_pos, freqs_cis, mask):
        """
        前向传播函数，计算注意力输出
        
        参数:
            x: 输入张量，形状为 (批次大小, 序列长度, 隐藏层维度)
            start_pos: 当前序列在整体序列中的起始位置，用于KV缓存定位
            freqs_cis: 旋转位置编码的复数表示，用于对查询和键进行位置编码
            mask: 注意力掩码，用于防止关注未来的序列元素，可选
        """
        # 获取输入张量的形状信息
        bsz, seqlen, _ = x.shape  # bsz: 批次大小, seqlen: 当前输入序列长度
        
        # 通过线性变换得到查询(query)、键(key)、值(value)
        xq, xk, xv = self.wq(x), self.wk(x), self.wv(x)

        # 重塑张量形状，将注意力头维度分离出来
        # 新形状: (批次大小, 序列长度, 头数, 头维度)
        xq = xq.view(bsz, seqlen, self.n_local_heads, self.head_dim)
        xk = xk.view(bsz, seqlen, self.n_local_heads, self.head_dim)
        xv = xv.view(bsz, seqlen, self.n_local_heads, self.head_dim)

        # 应用旋转位置编码，将位置信息融入查询和键中
        # 旋转位置编码有助于模型理解序列的位置关系
        xq, xk = apply_rotary_emb(xq, xk, freqs_cis=freqs_cis)

        # 确保缓存与输入张量在同一设备上（如GPU）
        self.cache_k = self.cache_k.to(xq)
        self.cache_v = self.cache_v.to(xq)

        # 将当前序列的键和值存入缓存中
        # 缓存的位置由start_pos决定，确保历史序列不会被覆盖
        self.cache_k[:bsz, start_pos : start_pos + seqlen] = xk
        self.cache_v[:bsz, start_pos : start_pos + seqlen] = xv

        # 从缓存中获取所有相关的键和值（包括历史序列和当前序列）
        # 形状: (批次大小, 总序列长度, 头数, 头维度)
        keys = self.cache_k[:bsz, : start_pos + seqlen]
        values = self.cache_v[:bsz, : start_pos + seqlen]

        # 调整张量维度顺序，为注意力计算做准备
        # 将注意力头维度提前，新形状: (批次大小, 头数, 序列长度, 头维度)
        xq = xq.transpose(1, 2)
        keys = keys.transpose(1, 2)
        values = values.transpose(1, 2)
        
        # 计算注意力分数: 查询与键的点积，除以头维度的平方根进行缩放
        # 形状: (批次大小, 头数, 当前序列长度, 总序列长度)
        scores = torch.matmul(xq, keys.transpose(2, 3)) / math.sqrt(self.head_dim)
        
        # 应用掩码（如果提供），通常用于 decoder 中防止关注未来的 token
        if mask is not None:
            scores = scores + mask  # 掩码通过加上一个很大的负值来实现
        
        # 对注意力分数进行 softmax 归一化，得到注意力权重
        scores = F.softmax(scores.float(), dim=-1).type_as(xq)
        
        # 注意力权重与值相乘，得到注意力输出
        # 形状: (批次大小, 头数, 当前序列长度, 头维度)
        output = torch.matmul(scores, values)
        
        # 调整维度顺序并重塑，将注意力头合并回隐藏层维度
        # 最终形状: (批次大小, 序列长度, 隐藏层维度)
        output = output.transpose(1, 2).contiguous().view(bsz, seqlen, -1)

        # 通过输出线性变换层，返回最终结果
        return self.wo(output)
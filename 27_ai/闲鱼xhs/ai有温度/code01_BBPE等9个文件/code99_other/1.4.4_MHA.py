import torch
import torch.nn as nn
import torch.nn.functional as F

class MHA(torch.nn.Module):
    def __init__(self, embed_size, nums_head):
        # 调用父类的构造函数
        super(MHA, self).__init__()
        self.embed_size = embed_size
        self.nums_head = nums_head
        self.head_dim = embed_size // nums_head

        assert(self.head_dim * nums_head == embed_size), "嵌入维度必须是头数的整数倍"

        self.query = nn.Linear(embed_size, embed_size)
        self.key = nn.Linear(embed_size, embed_size)
        self.value = nn.Linear(embed_size, embed_size)
        self.attention_dropout = nn.Dropout(0.1)
        self.output_linear = nn.Linear(embed_size, embed_size)

    def forward(self, datas, mask):
        '''
        前向传播函数
        :param datas: 输入序列的值矩阵，大小为 (batch_size, seq_length, embed_size)
        :param mask: 遮罩矩阵，用于掩蔽无效的序列位置，大小为 (batch_size, 1, seq_length, seq_length)
        :return: 输出注意力加权的值矩阵，注意力分数
        '''
        batch_size, seq, _ = datas.size()

        # 初始化 qkv
        # （batch_size, seq_length, embed_size） -> （batch_size, seq_length, nums_head, head_dim）
        Q = self.query(datas).view(batch_size, seq, self.nums_head, self.head_dim)
        K = self.key(datas).view(batch_size, seq, self.nums_head, self.head_dim)
        V = self.value(datas).view(batch_size, seq, self.nums_head, self.head_dim)

        # （batch_size, seq_length, nums_head, head_dim） -> （batch_size, nums_head, seq_length, head_dim）
        Q = Q.transpose(1, 2)
        K = K.transpose(1, 2)
        V = V.transpose(1, 2)

        '''
        计算缩放点积注意力。
         Q: (batch_size,num_heads,seq_len, d_k)
         K: (batch_size,num_heads,seq_len, d_k)
         V: (batch_size, num_heads, seq_length, d_v)
         mask:(batch_size, 1, seq_length, seq_length)
        '''
        
        # 计算查询和键的点积，得到注意力分数
        P = torch.matmul(Q, K.transpose(-2, -1))
        attention_scores = P/torch.sqrt(torch.tensor(Q.size(-1)))

        if mask is not None:
            # 将是0的，被掩盖的，填写成负无穷，在计算softmax时不考虑权重
            attention_scores = attention_scores.masked_fill(mask == 0, float('-inf'))

        # 对注意力进行归一化
        attention_wight = F.softmax(attention_scores, dim = -1)
        attention_wight = self.attention_dropout(attention_wight)

        # 计算加权注意力向量
        # （batch_size, nums_head, seq_length, seq_length） * （batch_size, nums_head, seq_length, head_dim）
        output = torch.matmul(attention_wight, V)
        # （batch_size, nums_head, seq_length, head_dim） -> （batch_size, seq_length, nums_head, head_dim）
        output = output.transpose(1, 2)
        # （batch_size, seq_length, nums_head, head_dim） -> （batch_size, seq_length, embed_size）
        output = output.reshape(batch_size, seq_length, self.embed_size)
        output = self.output_linear(output)
        return output

# 生成简单的数据进行测试
batch_size = 2  # 批次大小
seq_length = 3  # 序列长度
embed_size = 128  # 嵌入维度
heads = 8  # 头数

datas = torch.randn(batch_size, seq_length, embed_size)
# 生成遮罩矩阵，这里我们假设没有需要遮罩的位置，所以全为1
mask = torch.ones(batch_size, 1, seq_length, seq_length)
MHA_process = MHA(embed_size, heads)
output = MHA_process(datas, mask)
print("输出的形状:", output.shape)
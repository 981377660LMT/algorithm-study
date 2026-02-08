import numpy as np

def swiglu(x, W, V, b=None, c=None):
    """
    整合版SwiGLU激活函数，包含所有内部计算步骤
    
    参数:
    x -- 输入张量，形状为(..., input_dim)
    W -- 第一个权重矩阵，形状为(input_dim, hidden_dim)
    V -- 第二个权重矩阵，形状为(input_dim, hidden_dim)
    b -- 第一个偏置项，形状为(hidden_dim, )，可选
    c -- 第二个偏置项，形状为(hidden_dim, )，可选
    
    返回:
    SwiGLU计算结果，形状为(..., hidden_dim)
    """
    # 计算xW + b
    xW = np.dot(x, W)
    if b is not None:
        xW += b
    
    # 计算xV + c
    xV = np.dot(x, V)
    if c is not None:
        xV += c
    
    # 内部计算Sigmoid和Swish1
    sigmoid_xW = 1 / (1 + np.exp(-xW))  # Sigmoid计算
    swish1_xW = xW * sigmoid_xW         # Swish1计算
    
    # 元素-wise乘法得到最终结果
    return swish1_xW * xV
    
## softmax
X = torch.tensor([-0.2, 0.2, 0.5, 0.9, 0.1, 0.8])
X_exp_sum = X.exp().sum()
X_softmax_hand = torch.exp(X) / X_exp_sum
print(X_softmax_hand)
# tensor([0.0864, 0.1289, 0.1739, 0.2595, 0.1166, 0.2348])



#-------------------------------------------------------------------
# safesoftmax
X = torch.tensor([-0.2, 0.2, 0.5, 0.9, 0.1, 0.8])
X_max = X.max()
X_exp_sum_sub_max = torch.exp(X-X_max).sum()
X_safe_softmax_hand = torch.exp(X - X_max) / X_exp_sum_sub_max
print(X_safe_softmax_hand)
# tensor([0.0864, 0.1289, 0.1739, 0.2595, 0.1166, 0.2348])


#-------------------------------------------------------------------
# online softmax
X = torch.tensor([-0.2, 0.2, 0.5, 0.9, 0.1, 0.8])
X_pre = X[:-1]
print(X)
print(X_pre)
print(X[-1])

# 1. we calculative t-1 time Online Softmax
X_max_pre = X_pre.max()
X_sum_pre = torch.exp(X_pre - X_max_pre).sum()

# 2. we calculative t time Online Softmax
X_max_cur = torch.max(X_max_pre, X[-1])  # X[-1] is new data
X_sum_cur = X_sum_pre * torch.exp(X_max_pre - X_max_cur) + torch.exp(X[-1] - X_max_cur)

# 3. final we calculative online softmax
X_online_softmax = torch.exp(X - X_max_cur) / X_sum_cur
print("online softmax result: ", X_online_softmax)

# tensor([-0.2000,  0.2000,  0.5000,  0.9000,  0.1000])
# tensor(0.8000)

# online softmax result:  
#tensor([0.0864, 0.1289, 0.1739, 0.2595, 0.1166, 0.2348])

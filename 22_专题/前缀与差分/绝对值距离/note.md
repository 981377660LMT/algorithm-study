**母题：有序数组中差绝对值之和**
nums 是 `非递减有序`整数数组
公式 (ni-n0) + (ni-n1) + (ni-n2) +...+(ni-ni) +(ni+1-ni) +(ni+2-ni)
第 i 项等于
`ni*(i+1)-preSum[i+1]`+`preSum[n]-preSum[i]-ni*(n-i)`
**前面有 i+1 项 后面有 n-i 项**

**矩阵快速幂加快行间 dp 转移递推**
适用于

- dp 转移表达式不依赖于 i 的情况
- 线性转移会超时 (n>1e7)

- 题目的暗示
  n<=50 k 在 1e9 左右
  O(n^3logk)复杂度
  转移矩阵为`n*n`的矩阵 做 k 次矩阵乘法
- 项数 k 不影响迁移结果
  状態遷移の式を見ると、`回数 n が遷移に影響しない`。こういうふうに同じ状態遷移を何度も繰り返すものは、行列累乗かも。

## 矩阵快速幂加速方法

- 使用预处理 mat 的每个幂次
- 降低快速幂的 log 底数(cacheLevel)
- 将二维矩阵用一维数组存储
- 加速矩阵乘法
  https://qiita.com/nocturnality/items/92b3cf7d7fc5dab4a3a7

  1. i, j, k ループ交換法 命中缓存

  ```java

  // 慢,2100ms
  for (i = 0; i < N; i++){
      for (j = 0; j < N; j++){
          for (k = 0; k < N; k++){
              c[i*N+j]+=a[i*N+k]*b[k*N+j];
          }
      }
  }

  // 快,1200ms
  for (i = 0; i < N; i++){
      for (k = 0; k < N; k++){
          for (j = 0; j < N; j++){
              c[i*N+j]+=a[i*N+k]*b[k*N+j];
          }
      }
  }
  ```

  2. ループアンローリング 分段减少循环,注意缓存局部变量

  ```JAVA
  for (i = 0; i < N; i+=4){
    i1 = i + 1;
    i2 = i + 2;
    i3 = i + 3;
    for (k = 0; k < N; k+=4){
        k1 = k + 1;
        k2 = k + 2;
        k3 = k + 3;

        for (j = 0; j < N; j++){
            c[i*N+j]+=a[i*N+k]*b[k*N+j];
            c[i*N+j]+=a[i*N+k1]*b[k1*N+j];
            c[i*N+j]+=a[i*N+k2]*b[k2*N+j];
            c[i*N+j]+=a[i*N+k3]*b[k3*N+j];

            c[i1*N+j]+=a[i1*N+k]*b[k*N+j];
            c[i1*N+j]+=a[i1*N+k1]*b[k1*N+j];
            c[i1*N+j]+=a[i1*N+k2]*b[k2*N+j];
            c[i1*N+j]+=a[i1*N+k3]*b[k3*N+j];

            c[i2*N+j]+=a[i2*N+k]*b[k*N+j];
            c[i2*N+j]+=a[i2*N+k1]*b[k1*N+j];
            c[i2*N+j]+=a[i2*N+k2]*b[k2*N+j];
            c[i2*N+j]+=a[i2*N+k3]*b[k3*N+j];

            c[i3*N+j]+=a[i3*N+k]*b[k*N+j];
            c[i3*N+j]+=a[i3*N+k1]*b[k1*N+j];
            c[i3*N+j]+=a[i3*N+k2]*b[k2*N+j];
            c[i3*N+j]+=a[i3*N+k3]*b[k3*N+j];
        }
    }
  }
  ```

  3. キャッシュブロッキング 合理利用 CPU 局部缓存块的大小

  ```java
  BLOCK = 10;

  for (i = 0; i <N; i+=BLOCK){
      for(k = 0; k < N; k+=BLOCK){
          for(j = 0; j < N; j+=BLOCK){
              for(ii = i; ii < (i + BLOCK); ii+=4){
                  ii1 = ii + 1;
                  ii2 = ii + 2;
                  ii3 = ii + 3;
                  for (kk = k; kk < (k + BLOCK); kk+=4){
                      kk1 = kk + 1;
                      kk2 = kk + 2;
                      kk3 = kk + 3;

                      for (jj = j; jj < (j + BLOCK); jj++){
                          c[ii*N+jj]+=a[ii*N+kk]*b[kk*N+jj];
                          c[ii*N+jj]+=a[ii*N+kk1]*b[kk1*N+jj];
                          c[ii*N+jj]+=a[ii*N+kk2]*b[kk2*N+jj];
                          c[ii*N+jj]+=a[ii*N+kk3]*b[kk3*N+jj];

                          c[ii1*N+jj]+=a[ii1*N+kk]*b[kk*N+jj];
                          c[ii1*N+jj]+=a[ii1*N+kk1]*b[kk1*N+jj];
                          c[ii1*N+jj]+=a[ii1*N+kk2]*b[kk2*N+jj];
                          c[ii1*N+jj]+=a[ii1*N+kk3]*b[kk3*N+jj];

                          c[ii2*N+jj]+=a[ii2*N+kk]*b[kk*N+jj];
                          c[ii2*N+jj]+=a[ii2*N+kk1]*b[kk1*N+jj];
                          c[ii2*N+jj]+=a[ii2*N+kk2]*b[kk2*N+jj];
                          c[ii2*N+jj]+=a[ii2*N+kk3]*b[kk3*N+jj];

                          c[ii3*N+jj]+=a[ii3*N+kk]*b[kk*N+jj];
                          c[ii3*N+jj]+=a[ii3*N+kk1]*b[kk1*N+jj];
                          c[ii3*N+jj]+=a[ii3*N+kk2]*b[kk2*N+jj];
                          c[ii3*N+jj]+=a[ii3*N+kk3]*b[kk3*N+jj];
                      }
                  }
              }
          }
      }
  }
  ```

- 针对 0 值和稀疏矩阵优化

// 有两个整形数组A和B
// A是有序递增的
// B是有序递减的
// 假设A外的空间
// 把B合并到A后需要去重并保证有序递增
// 需要去重

void merge(int A[], int m, int B[], int n) {
    int i = m - 1, j = n - 1, k = m + n - 1;
    while (i >= 0 && j >= 0) {
        if (A[i] > B[j]) {
            A[k--] = A[i--];
        } else {
            A[k--] = B[j--];
        }
    }
    while (j >= 0) {
        A[k--] = B[j--];
    }

    int last = A[0];
    int index = 1;
    for (int i = 1; i < m + n; i++) {
        if (A[i] != last) {
            A[index++] = A[i];
            last = A[i];
        }
    }

    for (int i = index; i < m + n; i++) {
        A[i] = 0;
    }
}



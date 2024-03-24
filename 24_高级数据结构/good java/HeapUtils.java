package template.utils;

import java.util.Comparator;

public class HeapUtils {
    public static <T> void moveUp(T[] arr, int n, int index, Comparator<T> cmp) {
        while (index > 1) {
            int p = index >> 1;
            if (cmp.compare(arr[p], arr[index]) > 0) {
                SequenceUtils.swap(arr, p, index);
                index = p;
            } else {
                break;
            }
        }
    }

    public static <T> void moveDown(T[] arr, int n, int index, Comparator<T> cmp) {
        while (index <= (n >> 1)) {
            int best = left(index);
            if (right(index) <= n) {
                best = cmp.compare(arr[right(index)], arr[best]) < 0 ? right(index) : best;
            }
            if (cmp.compare(arr[best], arr[index]) < 0) {
                SequenceUtils.swap(arr, index, best);
                index = best;
            } else {
                break;
            }
        }
    }

    public static <T> void heapify(T[] arr, int n, Comparator<T> cmp) {
        for (int i = n / 2; i >= 1; i--) {
            moveDown(arr, n, i, cmp);
        }
    }

    private static int left(int i) {
        return i << 1;
    }

    private static int right(int i) {
        return (i << 1) | 1;
    }
}

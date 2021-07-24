import java.util.Arrays;

public class InsertionSort {

    private InsertionSort(){}

    public static <E extends Comparable<E>> void sort(E[] arr){

        for(int i = 0; i < arr.length; i ++){

            // 将 arr[i] 插入到合适的位置
            E t = arr[i];
            int j;
            for(j = i; j - 1 >= 0 && t.compareTo(arr[j - 1]) < 0; j --){
                arr[j] = arr[j - 1];
            }
            arr[j] = t;
        }
    }

    // 换个方法实现插入排序法，我们叫 sort2
    public static <E extends Comparable<E>> void sort2(E[] arr){

        for(int i = arr.length - 1; i >= 0; i --){

            // 将 arr[i] 插入到合适的位置
            E t = arr[i];
            int j;
            for(j = i; j + 1 < arr.length && t.compareTo(arr[j + 1]) > 0; j ++){
                arr[j] = arr[j + 1];
            }
            arr[j] = t;
        }
    }

    private static <E extends Comparable<E>> boolean isSorted(E[] arr){

        for(int i = 1; i < arr.length; i ++)
            if(arr[i - 1].compareTo(arr[i]) > 0)
                return false;
        return true;
    }

    public static void main(String[] args){

        int[] dataSize = {10000, 100000};
        for(int n: dataSize){

            System.out.println("Random Array : ");

            Integer[] arr = ArrayGenerator.generateRandomArray(n, n);
            Integer[] arr2 = Arrays.copyOf(arr, arr.length);
            SortingHelper.sortTest("InsertionSort", arr);
            SortingHelper.sortTest("InsertionSort2", arr2);

            System.out.println();

            System.out.println("Ordered Array : ");

            arr = ArrayGenerator.generateOrderedArray(n);
            arr2 = Arrays.copyOf(arr, arr.length);
            SortingHelper.sortTest("InsertionSort", arr);
            SortingHelper.sortTest("InsertionSort2", arr2);

            System.out.println();
        }
    }
}

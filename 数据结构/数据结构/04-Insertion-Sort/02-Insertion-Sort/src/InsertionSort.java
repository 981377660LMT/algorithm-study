public class InsertionSort {

    private InsertionSort(){}

    public static <E extends Comparable<E>> void sort(E[] arr){

        for(int i = 0; i < arr.length; i ++){

            // 将 arr[i] 插入到合适的位置
//            for(int j = i; j - 1 >= 0; j --){
//                if(arr[j].compareTo(arr[j - 1]) < 0)
//                    swap(arr, j - 1, j);
//                else break;
//            }

            for(int j = i; j - 1 >= 0 && arr[j].compareTo(arr[j - 1]) < 0; j --)
                swap(arr, j - 1, j);
        }
    }

    private static <E> void swap(E[] arr, int i, int j){

        E t = arr[i];
        arr[i] = arr[j];
        arr[j] = t;
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
            Integer[] arr = ArrayGenerator.generateRandomArray(n, n);
            SortingHelper.sortTest("InsertionSort", arr);
        }
    }
}

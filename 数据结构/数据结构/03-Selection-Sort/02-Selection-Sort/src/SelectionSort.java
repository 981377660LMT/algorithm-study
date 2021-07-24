public class SelectionSort {

    private SelectionSort(){}

    public static void sort(int[] arr){

        // arr[0...i) 是有序的；arr[i...n) 是无序的
        for(int i = 0; i < arr.length; i ++){

            // 选择 arr[i...n) 中的最小值的索引
            int minIndex = i;
            for(int j = i; j < arr.length; j ++){
                if(arr[j] < arr[minIndex])
                    minIndex = j;
            }

            swap(arr, i, minIndex);
        }
    }

    private static void swap(int[] arr, int i, int j){

        int t = arr[i];
        arr[i] = arr[j];
        arr[j] = t;
    }

    public static void main(String[] args){

        int[] arr = {1, 4, 2, 3, 6, 5};
        SelectionSort.sort(arr);
        for(int e: arr)
            System.out.print(e + " ");
        System.out.println();
    }
}

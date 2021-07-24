public class Main {

    public static void main(String[] args) {

        Array<Integer> array = new Array<>();
        LinkedList<Integer> list = new LinkedList<>();

        int n = 10000000;
        System.out.println("n = " + n);

        long startTime = System.nanoTime();
        for(int i = 0; i < n; i ++)
            array.addLast(i);
        long endTime = System.nanoTime();
        double time = (endTime - startTime) / 1000000000.0;
        System.out.println("Array : " + time + " s");

        startTime = System.nanoTime();
        for(int i = 0; i < n; i ++)
            list.addFirst(i);
        endTime = System.nanoTime();
        time = (endTime - startTime) / 1000000000.0;
        System.out.println("LinkedList : " + time + " s");
    }
}

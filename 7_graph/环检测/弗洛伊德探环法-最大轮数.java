package template.algo;

public class FloydSearchCircle {
    /**
     * Initially each participants located at start vertex v[1]
     */
    public static interface Interactor {
        void next(int i);

        boolean equal(int i, int j);
    }

    /**
     * <pre>
     * a circle comprised by start outcircle vertex and period incircle vertex
     * v[1] -> v[2] -> ... -> v[start] -> v[start + 1] -> v[start + 2] -> ... -> v[start + period - 1] -> v[start + 1] -> ...
     * use three participants in interactor
     * if (start + period) > maxRound, null will be returned
     * time complexity: O(2t+2c)
     * space complexity: O(1)
     * </pre>
     */
    public static int[] search(Interactor interactor, int maxRound) {
        interactor.next(0);
        interactor.next(0);
        interactor.next(1);
        maxRound--;
        while (!interactor.equal(0, 1) && maxRound > 0) {
            interactor.next(0);
            interactor.next(0);
            interactor.next(1);
            maxRound--;
        }
        if (!interactor.equal(0, 1)) {
            return null;
        }
        int start = 0;
        while (!interactor.equal(0, 2)) {
            interactor.next(0);
            interactor.next(2);
            start++;
        }
        int period = 1;
        interactor.next(0);
        while (!interactor.equal(0, 2)) {
            interactor.next(0);
            period++;
        }
        return new int[]{start, period};
    }

}
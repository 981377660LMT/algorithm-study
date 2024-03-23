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
     * a circle comprised by t outcircle vertex and c incircle vertex
     * v[1] -> v[2] -> ... -> v[t] -> v[t + 1] -> v[t + 2] -> ... -> v[t + c - 1] -> v[t + 1] -> ...
     * use three participants in interactor
     * if (t + c) > maxRound, null will be returned
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
        int t = 0;
        while (!interactor.equal(0, 2)) {
            interactor.next(0);
            interactor.next(2);
            t++;
        }
        int c = 1;
        interactor.next(0);
        while (!interactor.equal(0, 2)) {
            interactor.next(0);
            c++;
        }
        return new int[]{t, c};
    }

}
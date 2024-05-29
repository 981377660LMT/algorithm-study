package template.problem;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntegerArrayList;
import template.primitve.generated.datastructure.LongObjectHashMap;

/**
 * There is c points on a circle, points are enumerated as 0, 1, ... , c - 1.
 * <br>
 * points i and i + 1 are adjacent (0 and c - 1 are adjacent too).
 * <br>
 * There are some people and some houses(equal number). People can move between adjacent points, cost 1 for moving one step.
 * <br>
 * Now you should build a matching between people and houses. You need to make the total cost minimum
 */
public class OnCircleMinCostMatchProblem {
    int[] matching;
    long minimumCost;
    LongObjectHashMap<IntegerArrayList> peopleMap;
    LongObjectHashMap<IntegerArrayList> houseMap;
    CandyAssignProblem problem;

    public long getMinimumCost(){
        return minimumCost;
    }

    /**
     * Find which house will be lived by person i
     */
    public int getMateOf(int i){
        return matching[i];
    }

    private IntegerArrayList getIntegerList(LongObjectHashMap<IntegerArrayList> map, long key) {
        IntegerArrayList list = map.get(key);
        if (list == null) {
            list = new IntegerArrayList(1);
            map.put(key, list);
        }
        return list;
    }

    public OnCircleMinCostMatchProblem(int c, int[] people, int[] houses) {
        if (c <= 0 || people.length != houses.length) {
            throw new IllegalArgumentException();
        }

        int m = people.length;
        problem = new CandyAssignProblem(c, m * 2);
        peopleMap = new LongObjectHashMap<>(m, false);
        houseMap = new LongObjectHashMap<>(m, false);
        matching = new int[m];
        pending = new IntegerArrayList(m);

        for (int i = 0; i < m; i++) {
            getIntegerList(peopleMap, people[i]).add(i);
            getIntegerList(houseMap, houses[i]).add(i);
        }

        for (int i = 0; i < m; i++) {
            problem.assignCandyOn(people[i], 1);
        }

        for (int i = 0; i < m; i++) {
            problem.requireCandyOn(houses[i], 1);
        }

        problem.solve();
        minimumCost = problem.minimumCost;

        for (int i = 0; i < problem.candieCnt; i++) {
            int last = DigitUtils.mod(i - 1, problem.candieCnt);
            if (problem.candies[i].a >= 0 && problem.candies[last].a <= 0) {
                buildMatching(i);
            }
        }
    }

    IntegerArrayList pending;

    private void buildMatching(int i) {
        IntegerArrayList people = peopleMap.get(problem.candies[i].location);
        IntegerArrayList houses = houseMap.get(problem.candies[i].location);

        if (people != null && houses != null) {
            while (!people.isEmpty() && !houses.isEmpty()) {
                matching[people.pop()] = houses.pop();
            }
        }

        if (houses != null) {
            while (!houses.isEmpty() && !pending.isEmpty()) {
                matching[pending.pop()] = houses.pop();
            }
        }

        if (problem.candies[i].a > 0) {
            while (pending.size() < problem.candies[i].a) {
                pending.add(people.pop());
            }
            buildMatching((i + 1) % problem.candieCnt);
        }
        
        int last = DigitUtils.mod(i - 1, problem.candieCnt);
        if (problem.candies[last].a < 0) {
            while (pending.size() < -problem.candies[last].a) {
                pending.add(people.pop());
            }
            buildMatching(last);
        }
    }
}

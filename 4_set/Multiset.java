import java.util.TreeMap;

class MultiSet {
  TreeMap<Integer, Integer> map = new TreeMap<>();

  private int size = 0;

  public MultiSet() {
  }

  public void add(int val) {
    map.put(val, map.getOrDefault(val, 0) + 1);
    size++;
  }

  public void remove(int val) {
    map.put(val, map.get(val) - 1);
    if (map.get(val) == 0) {
      map.remove(val);
    }
    size--;
  }

  public int size() {
    return size;
  }

  public Integer higher(int val) {
    return map.higherKey(val);
  }

  public Integer lower(int val) {
    return map.lowerKey(val);
  }

  public Integer ceiling(int val) {
    return map.ceilingKey(val);
  }

  public Integer floor(int val) {
    return map.floorKey(val);
  }

  public Integer first() {
    if (map.isEmpty())
      return null;
    return map.firstKey();
  }

  public Integer last() {
    if (map.isEmpty())
      return null;
    return map.lastKey();
  }

  public boolean isEmpty() {
    return map.isEmpty();
  }

  @Override
  public String toString() {
    return map.toString();
  }
}
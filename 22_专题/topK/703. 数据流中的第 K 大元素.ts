import { MinHeap } from '../../8_heap/MinHeap'

class KthLargest {
  private readonly pq: MinHeap

  constructor(private readonly k: number, private readonly nums: number[]) {
    this.pq = new MinHeap((a, b) => a - b, k)
    nums.forEach(num => this.pq.heappush(num))
  }

  add(val: number): number {
    this.pq.heappush(val)
    return this.pq.heappop()!
  }

  static main() {
    const k = new KthLargest(3, [4, 5, 8, 2])
    console.log(k)
  }
}

KthLargest.main()

import { MinHeap } from '../../2_queue/minheap'

class KthLargest {
  private pq: MinHeap
  constructor(private k: number, private nums: number[]) {
    this.pq = new MinHeap((a, b) => a - b, k)
    nums.forEach(num => this.pq.push(num))
  }

  add(val: number): number {
    this.pq.push(val)
    return this.pq.shift()!
  }

  static main() {
    const k = new KthLargest(3, [4, 5, 8, 2])
    console.log(k)
  }
}
KthLargest.main()

"use strict";
// 1.插入 insert
// 2.删除堆顶 pop
// 3.获取堆顶 peek
// 4.获取堆大小 size
Object.defineProperty(exports, "__esModule", { value: true });
exports.MinHeap = void 0;
class MinHeap {
    constructor(heap, volumn) {
        this.heap = heap;
        this.volumn = volumn;
        this.heap = heap;
        this.volumn = volumn;
    }
    /**
     *
     * @param val 插入的值
     * @description 将值插入数组(堆)的尾部，然后上移直至父节点不超过它
     * @description 时间复杂度为`O(log(k))`
     */
    insert(val) {
        if (this.volumn !== undefined && this.heap.length >= this.volumn) {
            this.pop();
        }
        this.heap.push(val);
        this.shiftUp(this.heap.length - 1);
        return this;
    }
    /**
     * @description 用数组尾部元素替换堆顶(直接删除会破坏堆结构),然后下移动直至子节点都大于新堆顶
     * @description 时间复杂度为`O(log(k))`
     */
    pop() {
        this.heap[0] = this.heap.pop();
        this.shiftDown(0);
        return this;
    }
    peek() {
        return this.heap[0];
    }
    size() {
        return this.heap.length;
    }
    // 上移
    shiftUp(index) {
        if (index <= 0)
            return;
        const parentIndex = this.getParentIndex(index);
        while (this.heap[parentIndex] > this.heap[index]) {
            this.swap(parentIndex, index);
            this.shiftUp(parentIndex);
        }
    }
    // 下移
    shiftDown(index) {
        const leftChildIndex = this.getLeftChildIndex(index);
        const rightChildIndex = this.getRightChildIndex(index);
        if (this.heap[leftChildIndex] < this.heap[index]) {
            this.swap(leftChildIndex, index);
            this.shiftDown(leftChildIndex);
        }
        if (this.heap[rightChildIndex] < this.heap[index]) {
            this.swap(rightChildIndex, index);
            this.shiftDown(rightChildIndex);
        }
    }
    getParentIndex(index) {
        // 二进制数向右移动一位，相当于Math.floor((index-1)/2)
        return (index - 1) >> 1;
    }
    getLeftChildIndex(index) {
        return index * 2 + 1;
    }
    getRightChildIndex(index) {
        return index * 2 + 2;
    }
    swap(parentIndex, index) {
        ;
        [this.heap[parentIndex], this.heap[index]] = [this.heap[index], this.heap[parentIndex]];
    }
}
exports.MinHeap = MinHeap;
if (require.main === module) {
    const minHeap = new MinHeap([10, 20, 30]);
    console.log(minHeap);
    minHeap.insert(4);
    console.log(minHeap);
    minHeap.pop();
    console.log(minHeap);
}

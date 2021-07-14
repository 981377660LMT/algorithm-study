"use strict";
// 第k个最大元素=>直接使用最小堆
Object.defineProperty(exports, "__esModule", { value: true });
const _1_js___1 = require("./1_js\u5B9E\u73B0\u5806");
// 将数组的数值插入堆中，如果堆的容量超过K则删除堆顶
// 时间复杂度O(n*log(k))
const findKthLargest = (nums, k) => {
    const minHeap = new _1_js___1.MinHeap([], k);
    nums.forEach(num => minHeap.insert(num));
    return minHeap.peek();
};
console.log(findKthLargest([1, 2, 3, 4, 5, 7], 2));

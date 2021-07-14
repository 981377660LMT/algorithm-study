"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 每次抽出第i个元素 然后向前比较插入到合适的位置
const insertSort = (arr) => {
    for (let i = 1; i < arr.length; i++) {
        for (let j = i; j > 0; j--) {
            if (arr[j - 1] > arr[j]) {
                ;
                [arr[j - 1], arr[j]] = [arr[j], arr[j - 1]];
            }
        }
    }
};
const arr = [4, 1, 2, 5, 3, 6, 7];
insertSort(arr);
console.log(arr);

"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const biSearch = (arr, target) => {
    if (arr.length === 0)
        return -1;
    let leftPoint = 0;
    let rightPoint = arr.length - 1;
    while (leftPoint <= rightPoint) {
        const mid = Math.floor((leftPoint + rightPoint) / 2);
        const midElement = arr[mid];
        if (midElement === target) {
            return mid;
        }
        else if (midElement < target) {
            leftPoint = mid + 1;
        }
        else {
            rightPoint = mid - 1;
        }
    }
    return -1;
};
const arr = [1, 2, 3, 4, 5, 6, 7];
console.log(biSearch(arr, 3));

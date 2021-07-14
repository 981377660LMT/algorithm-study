"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const sequencialSearch = (arr, target) => {
    for (let i = 0; i < arr.length; i++) {
        if (arr[i] === target)
            return i;
    }
    return -1;
};
const arr = [4, 1, 2, 5, 3, 6, 7];
console.log(sequencialSearch(arr, 3));

"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const T = (matrix) => {
    for (const m of matrix) {
        m.reverse();
    }
    matrix.reverse();
    return matrix;
};
console.log(T([
    [1, 2, 3, 4],
    [5, 6, 7, 8],
    [9, 10, 11, 12],
]));

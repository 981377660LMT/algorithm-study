"use strict";
// 合并所有重叠的区间，并返回一个不重叠的区间数组
Object.defineProperty(exports, "__esModule", { value: true });
const merge = (intervals) => {
    if (intervals.length <= 1)
        return intervals;
    intervals.sort((a, b) => a[0] - b[0]);
    const res = [intervals[0]];
    for (let index = 1; index < intervals.length; index++) {
        const interval = intervals[index];
        const preLeft = res[res.length - 1][0];
        const preRight = res[res.length - 1][1];
        const curLeft = interval[0];
        const curRight = interval[1];
        // 三种关系:包含，相交，相离
        if (curRight <= preRight) {
            continue;
        }
        else if (curLeft <= preRight && curRight >= preRight) {
            res.pop();
            res.push([preLeft, curRight]);
        }
        else {
            res.push(interval);
        }
    }
    return res;
};
console.log(merge([
    [1, 3],
    [2, 6],
    [8, 10],
    [15, 18],
]));

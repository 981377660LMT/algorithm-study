"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const json = [
    {
        id: 1,
        children: [
            {
                id: 3,
                children: [
                    {
                        id: 5,
                        children: [],
                    },
                ],
            },
            {
                id: 4,
                children: [],
            },
        ],
    },
    {
        id: 2,
        children: [],
    },
];
// 为每个node增加key属性
const bfs = (n) => {
    console.log(n);
    n.key = n.id;
    n.children.map(node => {
        return bfs(node);
    });
};
json.map(bfs);
console.dir(json, { depth: null });
exports.default = {};

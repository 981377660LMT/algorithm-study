// https://www.developers.pub/article/1171875

// "mappings": "CAAA,WACE,IAAK,IAAIA,
// EAAI,EAAGA,EAAI,EAAGA,IACrBC,QAAQC,IAAI,KAGhBC",

/**
 * 为了尽可能减少存储空间但同时要达到记录原始位置和目标位置映射关系的目的，mappings字段按照了一些特殊的规则来生成。
 * 对应 sourcemap 的mappings，有三种形式
 * 1. [number]：表示行数相同，列数相同
 */
export type SourceMapSegment = [number] | [number, number, number, number] | [number, number, number, number, number]

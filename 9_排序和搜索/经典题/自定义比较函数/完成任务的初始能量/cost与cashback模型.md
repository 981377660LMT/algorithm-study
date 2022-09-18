transactions[i] = [costi, cashbacki]
为了完成交易 i ,money >= costi 这个条件必须为真。
执行交易后，你的钱数 money 变成 money - costi + cashbacki 。
请你返回 任意一种 交易顺序下，你都能完成所有交易的最少钱数 money 是多少

交易最差的情况:
贪心排序:

1. 亏钱的排前面
2. 亏钱的中, cashback 小的排前面 (不给发展机会，!前面 cashback 越小就越难)
3. 赚钱的中, cost 大的排前面 (拦路虎，!前面 cost 越大就越难)

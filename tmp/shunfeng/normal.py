# 	顺丰里面有一个处罚管理系统，为了管理处罚系统，请你帮忙实现如下函数:


# 	List<HistoryResult> query(int pageNo, int pageSize) :查询处罚记录，记录按照id升序排序


# 	List<HistoryResult> getByUserId(int userId) ：查询某个用户的处罚记录，记录按照id升序排序


# 	int punish(String operatorUserName, int userId, int punishStatus) : 处罚操作, 如果用户已经有被处罚了，新的处罚不能低于当前处罚等级才能生效，operatorUserName为执行处罚的人，userId为被处罚用户的编号


# 	int relieve(String operatorUserName, int userId) : 解除当前处罚，operatorUserName为解除处罚的人，userId为被解除处罚用户的编号。若当前用户无处罚记录则该操作无效。解除处罚的status设置为0


# 	只有有效的处罚或解除处罚操作，才会被记录id。这两种操作的返回值为该操作的id，从1开始自增。若操作为无效操作，返回-1


#  输入描述 第一行输入一个N(2<N<=200000),代表操作个数接下来的每2行代表1个操作，每2行的第1行输入一个操作的英文名字，代表对应的操作每2行的第2行输入操作对应的参数 输出描述 对应每个查询操作，记录一个结果，最后统一输出所有处罚操作和解除处罚操作的返回值，你只需要完成函数即可，输出无需理会

from collections import defaultdict
from typing import List


# id从1开始自增
class HistoryResult:
    def __init__(self, _id: int, user_id: int, operatorUserName: str, status: int):
        self.id = _id
        self.user_id = user_id
        self.operatorUserName = operatorUserName
        self.status = status

    def get_id(self):
        """处罚记录id"""
        return self.id

    def get_user_id(self):
        """被处罚的用户id"""
        return self.user_id

    def get_operator_user_name(self):
        """操作人的用户名"""
        return self.operatorUserName

    def get_status(self):
        """处罚等级.解除处罚为0."""
        return self.status


INF = int(1e18)


class Main:
    def __init__(self):
        self.id = 1
        self.history = []  # HistoryResult
        self.idToPu = dict()  # id => punish
        self.userPids = defaultdict(list)  # user_id => [punish]
        self.userLevel = defaultdict(list)

    # 查询处罚记录，记录按照id升序排序
    def query(self, page_no, page_size) -> List["HistoryResult"]:
        start = (page_no - 1) * page_size
        end = page_no * page_size
        res = self.history[start:end][:]
        res.sort(key=lambda x: x.id)
        return res

    # 获得查询某个用户的处罚记录，记录按照id升序排序
    def get_by_user_id(self, user_id) -> List["HistoryResult"]:
        ids = self.userPids[user_id]
        res = [self.idToPu[id] for id in ids][:]
        res.sort(key=lambda x: x.id)
        return res

    # 处罚操作, 如果用户已经有被处罚了，新的处罚不能低于当前处罚等级才能生效
    # 返回操作的id,从1开始自增.无效操作返回-1
    def punish(self, operator_user_name, user_id, punish_status) -> int:
        preLevel = self.userLevel[user_id]
        if not preLevel or punish_status >= preLevel[-1]:
            newPu = HistoryResult(self.id, user_id, operator_user_name, punish_status)
            self.history.append(newPu)
            self.idToPu[self.id] = newPu
            self.userPids[user_id].append(self.id)
            self.userLevel[user_id].append(punish_status)
            self.id += 1
            return self.id - 1
        return -1

    # 解除当前处罚, 如果当前用户正在被处罚中，解除当前处罚，返回处罚记录id，如果用户没有被处罚，返回-1表示解除处罚非法
    # 返回操作的id,从1开始自增.无效操作返回-1
    def relieve(self, operator_user_name, user_id) -> int:
        level = self.userLevel[user_id]
        if not level:
            return -1

        pu = HistoryResult(self.id, user_id, operator_user_name, 0)
        self.history.append(pu)
        self.idToPu[self.id] = pu
        self.userPids[user_id].append(self.id)
        self.userLevel[user_id].pop()
        self.id += 1
        return self.id - 1


if __name__ == "__main__":
    main = Main()
    op_num = int(input())
    print_history = []
    print_operators = []
    for i in range(op_num):
        op = input()
        data = input()
        if op == "punish":
            operator_user_name, user_id, punish_status = map(str, data.split())
            print_operators.append(
                main.punish(operator_user_name, int(user_id), int(punish_status))
            )
        elif op == "relieve":
            operator_user_name, user_id = map(str, data.split())
            print_operators.append(main.relieve(operator_user_name, int(user_id)))
        elif op == "getByUserId":
            user_id = int(data)
            results = main.get_by_user_id(user_id)
            if results:
                print_history.extend(results)
        elif op == "query":
            page_no, page_size = map(int, data.split())
            results = main.query(page_no, page_size)
            if results:
                print_history.extend(results)
        else:
            print("错误的输入")
    for result in print_history:
        print(
            str(result.get_id())
            + "_"
            + str(result.get_user_id())
            + "_"
            + result.get_operator_user_name()
            + "_"
            + str(result.get_status())
        )

    for operator_result in print_operators:
        print(operator_result)

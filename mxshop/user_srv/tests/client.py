import grpc
from user_srv.proto import user_pb2_grpc, user_pb2

class UserTest:
    def __init__(self):
        # 连接grpc 服务器
        channel = grpc.insecure_channel("127.0.0.1:50051")
        self.stub = user_pb2_grpc.UserStub(channel)

    def user_list(self):
        rsp: user_pb2.UserListResponse = self.stub.GetUserList(user_pb2.PageInfo(pn=1, pSize=5))
        # print(rsp.total)
        for user in rsp.data:
            print(user.mobile, user.birthDay)

    def get_user_by_mobile(self):
        rsp: user_pb2.UserInfoResponse = self.stub.GetUserByMobile(user_pb2.MobileReques(mobile="18782222222"))
        print(rsp)

    def get_user_by_id(self):
        rsp: user_pb2.UserInfoResponse = self.stub.GetUserById(user_pb2.IdRequests(id=4))
        print(rsp)

    def create_user(self):
        rsp: user_pb2.CreateUserInfo = self.stub.CreateUser(
            user_pb2.CreateUserInfo(
                mobile="18212341232",
                password="a1234",
                nick_name="bobby6"
            )
        )

    def update_user(self):
        rsp: user_pb2.UpdateUserInfo = self.stub.UpdateUser(
            user_pb2.UpdateUserInfo(
                id=3,
                nickName="jack",
                gender="男",
            )
        )

if __name__ == '__main__':
    user = UserTest()
    user.user_list()
    # user.get_user_by_mobile()
    # user.get_user_by_id()
    # user.create_user()
    # user.update_user()
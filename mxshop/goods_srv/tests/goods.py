import grpc
from goods_srv.proto import goods_pb2_grpc, goods_pb2

class GoodsTest:
    def __init__(self):
        # 连接grpc 服务器
        channel = grpc.insecure_channel("127.0.0.1:49457")
        self.goods_stub = goods_pb2_grpc.GoodsStub(channel)

    def goods_list(self):
        rsp: goods_pb2.GoodsListResponse = self.goods_stub.GoodsList(
            goods_pb2.GoodsFilterRequest(priceMin = 50)
        )
        # print(rsp.total)
        for item in rsp.data:
            print(item.name, item.shopPrice)



if __name__ == '__main__':
    user = GoodsTest()
    user.goods_list()
    # user.get_user_by_mobile()
    # user.get_user_by_id()
    # user.create_user()
    # user.update_user()
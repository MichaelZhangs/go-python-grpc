from user_srv.model.models import User
from user_srv.proto import user_pb2, user_pb2_grpc
import time
from loguru import logger
import grpc
from passlib.hash import pbkdf2_sha256
from datetime import date
from google.protobuf import empty_pb2


class UserServicer(user_pb2_grpc.UserServicer):

    def convert_user_to_rsp(self, user):
        user_info_rsp = user_pb2.UserInfoResponse()
        user_info_rsp.id = user.id
        user_info_rsp.passWord = user.password
        user_info_rsp.mobile = user.mobile
        user_info_rsp.role = user.role

        if user.nick_name:
            user_info_rsp.nickName = user.nick_name
        if user.gender:
            user_info_rsp.gender = user.gender
        if user.birthday:
            # print(user.birthday, type(user.birthday))
            # print(user.birthday.timetuple())
            user_info_rsp.birthDay = int(time.mktime(user.birthday.timetuple()))
        return user_info_rsp

    @logger.catch
    def GetUserList(self, request: user_pb2.PageInfo, context):
        #  获取用户的列表
        rsp = user_pb2.UserListResponse()
        users = User.select()
        rsp.total = users.count()
        start = 0
        per_page_nums = 10
        if request.pSize:
            per_page_nums = request.pSize
        if request.pn:
            start = per_page_nums * (request.pn - 1)

        users = users.limit(per_page_nums).offset(start)

        for user in users:
            user_info_rsp = self.convert_user_to_rsp(user)
            rsp.data.append(user_info_rsp)
        return rsp

    @logger.catch
    def GetUserByMobile(self, request: user_pb2.MobileReques, context):

            try:
                user = User.get(User.mobile == request.mobile)
            except Exception:
                context.set_code(grpc.StatusCode.NOT_FOUND)
                context.set_details("用户不存在")
                return user_pb2.UserInfoResponse()

            rsp = self.convert_user_to_rsp(user)
            return rsp

    @logger.catch
    def GetUserById(self, request: user_pb2.IdRequests, context):
        try:
            user = User.get(User.id == request.id)
        except Exception:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在")
            return user_pb2.UserInfoResponse()

        rsp = self.convert_user_to_rsp(user)
        return rsp

    @logger.catch
    def CreateUser(self, request: user_pb2.CreateUserInfo, context):
        # 新建用户
        try:
            user = User.get(User.Mobile == request.mobile)
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details("用户已经存在")
            return user_pb2.UserInfoResponse()
        except Exception:
            pass

        user = User()
        user.nick_name = request.nick_name
        user.mobile = request.mobile
        user.password = pbkdf2_sha256.hash(request.password)
        user.save()
        return self.convert_user_to_rsp(user)

    @logger.catch
    def UpdateUser(self, request: user_pb2.UpdateUserInfo, context):
        #  更新用户
        try:
            user = User.get(User.id == request.id)
        except Exception:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details("用户不存在")
            return user_pb2.UserInfoResponse()

        user.nick_name = request.nickName
        user.gender = request.gender
        user.birthday = date.fromtimestamp(request.birthDay)
        user.save()
        return empty_pb2.Empty()

    @logger.catch
    def CheckPassWord(self, request: user_pb2.PasswordCheckInfo, context):
        return user_pb2.CheckResponse(success=pbkdf2_sha256.verify(request.password, request.encryptPassword))



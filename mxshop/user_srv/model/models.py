from peewee import *
from user_srv.setting import settings

class BaseMode(Model):
    class Meta:
        database = settings.DB



class User(BaseMode):
    # 用户模型
    GENDER_CHOICES = (
        ("female", "女"),
        ("male", "男")
    )
    ROLE_CHOICES = (
        (1, "普通用户"),
        (2, "管理员")
    )

    mobile = CharField(max_length=11, index=True, unique=True, verbose_name="手机号码")
    password = CharField(max_length=100, verbose_name="密码")
    nick_name = CharField(max_length=20, null=True, verbose_name="昵称")
    head_url = CharField(max_length=200, null=True, verbose_name="头像")
    birthday = DateField(null=True, verbose_name="生日")
    address = CharField(max_length=200, null=True, verbose_name="地址")
    desc = TextField(null=True, verbose_name="个人简介")
    gender = CharField(max_length=6, choices= GENDER_CHOICES, null=True, verbose_name="性别")
    role = IntegerField(default=1, choices=ROLE_CHOICES, verbose_name="")
    # salt = CharField(max_length=200, null=True, verbose_name="地址")

if __name__ == '__main__':
    # 生成表结构
    # settings.DB.create_tables([User])
    # import hashlib
    # m = hashlib.md5()
    # print("sssse13".encode())
    # m.update(b"ssss")
    # print(m.hexdigest())
    # 9f6e6800cfae7749eb6c486619254b9c
    # 8f60c8102d29fcd525162d02eed4566b
    from passlib.hash import pbkdf2_sha256
    # h = pbkdf2_sha256.hash("123456")
    # print(h)
    # # 验证
    # print(pbkdf2_sha256.verify("123456", h))
    # print(pbkdf2_sha256.verify("123465", h))
    # for i in range(10):
    #     user = User()
    #     user.nick_name = f"bobby{i}"
    #     user.mobile = f"1878222222{i}"
    #     user.password = pbkdf2_sha256.hash("admin123")
    #     user.save()
    import time
    from datetime import date
    for user in User.select():
        print(user.mobile , pbkdf2_sha256.verify("admin123", user.password))
        # if user.birthday:
            # print(user.birthday)
            # utime = int(time.mktime(user.birthday.timetuple()))
            # print(utime)
            # print(date.fromtimestamp(utime))
    users = User.get(User.mobile=="18782222222")
    # print(users.mobile)


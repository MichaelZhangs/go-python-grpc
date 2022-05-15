from playhouse.pool import PooledMySQLDatabase
from playhouse.shortcuts import ReconnectMixin

# 不用reconnectMirxin 会出现数据库断开连接


class ReconnectMysqlDatabase(ReconnectMixin, PooledMySQLDatabase):
    pass

MYSQL_DB = "goods_srv"
MYSQL_HOST = "127.0.0.1"
MYSQL_PORT = 3306
MYSQL_USER = "root"
MYSQL_PASSWORD = "123456"

DB = ReconnectMysqlDatabase(MYSQL_DB, host=MYSQL_HOST, port=MYSQL_PORT, user=MYSQL_USER, password=MYSQL_PASSWORD)

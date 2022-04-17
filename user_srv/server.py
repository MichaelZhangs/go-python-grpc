import sys
import os
BASE_DIR = os.path.dirname(os.path.abspath(os.path.dirname(__file__)))
# BASE_DIR = os.path.dirname(__file__)
sys.path.insert(0, BASE_DIR)

import grpc
from user_srv.proto import  user_pb2_grpc
from concurrent import futures
from user_srv.handler.user import UserServicer
from loguru import logger
import signal

# from user_srv.logs.my_log import mylogger
BASE_DIR = os.path.dirname(os.path.abspath(os.path.dirname(__file__)))
# BASE_DIR = os.path.dirname(__file__)
sys.path.insert(0, BASE_DIR)
print(BASE_DIR)
logger.add("user_srv_{time:YYYY-MM-DD}.log", rotation="1 MB")
import argparse





def on_exist(signo, frame):
    logger.info("终止进程.... ")
    sys.exit(0)

def server():

    parser = argparse.ArgumentParser()
    parser.add_argument("--ip",
                        nargs="?",
                        type=str,
                        default="127.0.0.1",
                        help="binding ip")
    parser.add_argument("--port",
                        nargs="?",
                        type=int,
                        default=50051,
                        help="port")
    args = parser.parse_args()
    print(args.ip, args.port)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=3))
    user_pb2_grpc.add_UserServicer_to_server(UserServicer(), server)
    server.add_insecure_port(f"{args.ip}:{args.port}")
    # 主进程退出信号监听
    """
    windows 下支持的信号是有限的
        SIGINT ctrl+C 中断
        SIGTREM: kill
    """
    signal.signal(signal.SIGINT, on_exist)
    signal.signal(signal.SIGTERM, on_exist)
    logger.info(f"启动服务器 {args.ip}:{args.port}")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    # logging.basicConfig()
    server()
    # my_fun(0)


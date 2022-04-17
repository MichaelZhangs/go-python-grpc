
from loguru import logger

class Loggings:
    __instance = None

    def __new__(cls, *args, **kwargs):
        if not cls.__instance:
            cls.__instance = super(Loggings, cls).__new__(cls, *args, **kwargs)

        return cls.__instance

    def info(self, msg):
        return logger.info(msg)

    def debug(self, msg):
        return logger.debug(msg)

    def warning(self, msg):
        return logger.warning(msg)

    def error(self, msg):
        return logger.error(msg)


mylogger = Loggings()
if __name__ == '__main__':
    loggings = Loggings()
    loggings.info("中文test")



# if __name__ == '__main__':
#     mylog = MyLog()
#     mylog.info("adsdd")

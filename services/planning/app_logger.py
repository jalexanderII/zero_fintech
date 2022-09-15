import logging

_log_format = f"%(asctime)s [%(levelname)s] %(message)s"


class CustomFilter(logging.Filter):
    COLOR = {
        "DEBUG": "GREEN",
        "INFO": "GREEN",
        "WARNING": "YELLOW",
        "ERROR": "RED",
        "CRITICAL": "RED",
    }

    def filter(self, record):
        record.color = CustomFilter.COLOR[record.levelname]
        return True


def get_file_handler():
    file_handler = logging.FileHandler("x.log")
    file_handler.setLevel(logging.INFO)
    file_handler.setFormatter(logging.Formatter(_log_format))
    return file_handler


def get_stream_handler():
    stream_handler = logging.StreamHandler()
    stream_handler.setLevel(logging.DEBUG)
    stream_handler.setFormatter(logging.Formatter(_log_format))
    return stream_handler


def get_logger(name):
    logger = logging.getLogger(name)
    logger.setLevel(logging.DEBUG)
    # logger.addHandler(get_file_handler())
    logger.addHandler(get_stream_handler())
    logger.addFilter(CustomFilter())
    return logger

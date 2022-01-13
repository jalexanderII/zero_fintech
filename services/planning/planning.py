from concurrent import futures
import grpc

from server.server import PlanningServicer
from gen.planning import planning_pb2_grpc

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    planning_pb2_grpc.add_PlanningServicer_to_server(
        PlanningServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    # logging.basicConfig()
    serve()
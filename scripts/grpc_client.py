import asyncio

from grpclib.client import Channel

from gempellm.logistic_parcel_api.v1.logistic_parcel_api_grpc import logisticparcelApiServiceStub
from gempellm.logistic_parcel_api.v1.logistic_parcel_api_pb2 import DescribeParcelV1Request

async def main():
    async with Channel('127.0.0.1', 8082) as channel:
        client = logisticparcelApiServiceStub(channel)

        req = DescribeParcelV1Request(parcel_id=1)
        reply = await client.DescribeParcelV1(req)
        print(reply.message)


if __name__ == '__main__':
    asyncio.run(main())

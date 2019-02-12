# Basic gRPC in Python

Contains a minimal working example for gRPC in Python.
Based on https://github.com/ramananbalakrishnan/basic-grpc-python

## Quickstart

```shell
pip install -r requirements.txt
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. grpcDemo.proto
python backend.py
python frontend.py
```
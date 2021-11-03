#!/bin/bash
GRPCURL_BIN=/home/maxima/soft/grpcurl/grpcurl

$GRPCURL_BIN -plaintext -d '{"foo": "1"}'  127.0.0.1:8082 ozonmp.bss_workplace_api.v1.BssWorkplaceApiService/CreateWorkplaceV1
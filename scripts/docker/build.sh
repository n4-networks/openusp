#!/bin/sh
# Copyright 2023 N4-Networks.com
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# Controller
docker buildx build -t n4networks/openusp-controller:latest -f build/controller/Dockerfile --push --platform=linux/amd64,linux/arm64 .

# ApiServer
docker buildx build -t n4networks/openusp-apiserver:latest -f build/apiserver/Dockerfile --push --platform=linux/amd64,linux/arm64 .

# Cli
docker buildx build -t n4networks/openusp-cli:latest -f build/cli/Dockerfile --push --platform=linux/amd64,linux/arm64 .

# OBUSPA
docker buildx build -t n4networks/openusp-agent:latest -f build/obuspa/Dockerfile --push --platform=linux/amd64,linux/arm64 .

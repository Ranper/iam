#!/usr/bin/env bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}/../..")

# iam 配置
readonly IAM_CONFIG_DIR=${IAM_CONFIG_DIR:-/etc/iam} # iam 配置文件存放目录
readonly IAM_LOG_DIR=${IAM_LOG_DIR:-/var/log/iam} # iam 日志文件存放目录

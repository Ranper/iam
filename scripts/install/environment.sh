#!/usr/bin/env bash

IAM_ROOT=$(dirname "${BASH_SOURCE[0]}/../..")

# iam 配置
readonly IAM_DATA_DIR=${IAM_DATA_DIR:-/data/iam} # iam 各组件数据目录
readonly IAM_INSTALL_DIR=${IAM_INSTALL_DIR:-/opt/iam} # iam 安装文件存放目录
readonly IAM_CONFIG_DIR=${IAM_CONFIG_DIR:-/etc/iam} # iam 配置文件存放目录
readonly IAM_LOG_DIR=${IAM_LOG_DIR:-/var/log/iam} # iam 日志文件存放目录
readonly CA_FILE=${CA_FILE:-${IAM_CONFIG_DIR}/cert/ca.pem} # CA

# iam-apiserver 配置
readonly IAM_APISERVER_SECURE_BIND_ADDRESS=${IAM_APISERVER_SECURE_BIND_ADDRESS:-0.0.0.0}
readonly IAM_APISERVER_SECURE_BIND_PORT=${IAM_APISERVER_SECURE_BIND_PORT:-8443}
readonly IAM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE=${IAM_APISERVER_SECURE_TLS_CERT_KEY_CERT_FILE:-${IAM_CONFIG_DIR}/cert/iam-apiserver.pem}
readonly IAM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE=${IAM_APISERVER_SECURE_TLS_CERT_KEY_PRIVATE_KEY_FILE:-${IAM_CONFIG_DIR}/cert/iam-apiserver-key.pem}
#!/usr/bin/env bash

env_file="$1"
template_file="$2"


IAM_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${IAM_ROOT}/scripts/lib/init.sh"

if [ $# -ne 2 ]; then
    ima::log::error "Usage: genconfig.sh scripts/environment.sh configs/iam-apiserver.yaml"
    exit 1
fi

source "${env_file}"

# 声明关联数组（用于存储环境变量）
declare -A envs

# 临时禁用未定义变量报错
set +u

# 遍历模板文件中所有需要替换的变量
# ^[^#]：匹配不以#开头的行（排除注释行）
# \\1：替换为第一个捕获组（即环境变量名）
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ]; then
        iam::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" = "true" ]; then
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
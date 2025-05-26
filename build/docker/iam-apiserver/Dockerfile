FROM BASE_IMAGE
# LABEL maintainer="<colin404@foxmail.com>"

# 设置容器内部的工作目录, 后续的RUN、CMD、COPY等指令默认会在/opt/iam目录下执行。
WORKDIR /opt/iam
# docker build 的构建上下文
# 构建上下文是主机上的目录（如docker build -t myimage .中的.）。
# 它决定了COPY或ADD指令可以访问哪些主机文件。
# 关键区别：
# WORKDIR是容器内部的路径，影响容器运行时的行为。
# 构建上下文是主机路径，仅用于构建时提供文件。

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/iam

COPY iam-apiserver /opt/iam/bin/
ENTRYPOINT ["/opt/iam/bin/iam-apiserver"]
CMD ["-c", "/etc/iam/iam-apiserver.yaml"]

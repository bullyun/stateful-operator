FROM registry.cn-hangzhou.aliyuncs.com/bullyun/ubi-minimal:1.0

ENV OPERATOR=/usr/local/bin/stateful-operator \
    USER_UID=1001 \
    USER_NAME=stateful-operator

# install operator binary
COPY build/_output/bin/stateful-operator ${OPERATOR}

COPY build/bin /usr/local/bin

RUN chmod +x /usr/local/bin/*
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}

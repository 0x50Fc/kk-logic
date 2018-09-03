FROM alpine:latest

COPY ./etc/timezone /etc/timezone

COPY ./etc/localtime /etc/localtime

COPY ./main /bin/kk-logic

RUN chmod +x /bin/kk-logic

ENV KK_ENV_PORT 80

ENV KK_ENV_DIR /home

ENV KK_ENV_PREFIX /

ENV KK_ENV_SESSION_KEY kk

ENV KK_ENV_SESSION_MAX_AGE 1800

ENV KK_ENV_MAX_MEMORY 4096000

VOLUME /home

EXPOSE 80

CMD kk-logic -p $KK_ENV_PORT -r $KK_ENV_DIR --prefix $KK_ENV_PREFIX --sessionKey $KK_ENV_SESSION_KEY --sessionMaxAge $KK_ENV_SESSION_MAX_AGE --maxMemory $KK_ENV_MAX_MEMORY


FROM public.ecr.aws/lambda/provided:al2 as build
LABEL Terry Kim <terry960302@gmail.com>

ARG DIR=flower-batch-upload-scheduler

RUN yum install -y golang
RUN go env -w GOPROXY=direct

ADD ${DIR}/go.mod ${DIR}/go.sum ./
RUN go mod download

ADD ${DIR}/ .
RUN env GOOS=linux GOARCH=amd64 go build -o /main

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main

ENV PROFILE=prod
ENV DB_USER=marketbill
ENV DB_PW=marketbill1234!
ENV DB_NET=tcp
ENV DB_HOST=marketbill-db.ciegftzvpg1l.ap-northeast-2.rds.amazonaws.com
ENV DB_PORT=5432
ENV DB_NAME=prod-db

ENV PORT=8080
ENV API_KEY=4DC6A10B4F5D43D5977F364FC0DFE81C

ENTRYPOINT [ "/main" ]
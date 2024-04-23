FROM golang:1.22.2-alpine3.19

WORKDIR /HOME

COPY ./pkg /HOME/

RUN cd /HOME && go build -o kubelib

CMD [ "/HOME/kubelib" ]

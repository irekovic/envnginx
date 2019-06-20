FROM golang:alpine

ADD ./ /envnginx/
RUN cd /envnginx && go build

FROM nginx:alpine
COPY start.sh .
COPY --from=0 /envnginx/envnginx .
WORKDIR /
CMD ["./start.sh"]
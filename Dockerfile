FROM ubuntu:18.04
WORKDIR gogsc
COPY gsc .
RUN mkdir -p /home/runner/work/go-gsc/go-gsc/vendor/github.com/yanyiwu
COPY vendor/github.com/yanyiwu/ /home/runner/work/go-gsc/go-gsc/vendor/github.com/yanyiwu/
CMD ["./gsc"]

FROM golang:1.17.1 as builder

RUN mkdir /action
COPY . /action
WORKDIR /action
RUN GO111MODULE=on go list /action/...  && go build -o /action/pr_title_lint_action /action/main.go

FROM ubuntu:18.04
RUN apt-get update -y \
        && apt-get install ca-certificates -y --no-install-recommends
RUN mkdir /action
COPY --from=builder /action/pr_title_lint_action /action/pr_title_lint_action
ENTRYPOINT ["/action/pr_title_lint_action"]
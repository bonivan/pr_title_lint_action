FROM golang:1.17.1 as builder
RUN mkdir /action
COPY . /action
WORKDIR /action
RUN GO111MODULE=on go list /action/...  && go build -o /action/pr_title_lint_action /action/main.go

FROM public.ecr.aws/lts/ubuntu:latest
RUN mkdir /action
COPY --from=builder /action/pr_title_lint_action /action/pr_title_lint_action
WORKDIR /action
ENTRYPOINT ["./pr_title_lint_action"]
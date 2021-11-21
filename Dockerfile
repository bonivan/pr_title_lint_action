FROM golang:1.16.3 as builder
RUN mkdir /action
COPY . /action
RUN go build -o /action/pr_title_lint_action /action/main.go

FROM docker.internal.sysdig.com/sysdig-mini-ubi:1.1.12
RUN mkdir /action
COPY --from=builder /action/pr_title_lint_action /action/pr_title_lint_action
WORKDIR /action
ENTRYPOINT ["/pr_title_lint_action"]

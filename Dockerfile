FROM golang:1.16.3 as builder

COPY . .

GOPRIVATE=github.com/bonivan/pr_lint_action GO111MODULE=on go list ./...
go test ./...
go build -o /pr_lint_action ./cmd/main.go

FROM docker.internal.sysdig.com/sysdig-mini-ubi:1.1.12

COPY --from=builder /pr_lint_action /pr_lint_action

ENTRYPOINT ["/pr_lint_action"]

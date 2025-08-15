# 멀티스테이지 빌드를 사용한 Go 파일 서버 Dockerfile

# 빌드 스테이지
FROM golang:1.22-alpine AS builder

# 필요한 패키지 설치
RUN apk add --no-cache git ca-certificates tzdata

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사 및 의존성 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사
COPY . .

# 바이너리 빌드 (정적 링킹, 최적화)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o server ./cmd/server

# 최종 실행 스테이지
FROM alpine:3.19

# 필요한 패키지 설치 (CA 인증서, 타임존)
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 타임존 설정
ENV TZ=Asia/Seoul

# 작업 디렉토리 생성
WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/server /app/server

# 스토리지 디렉토리 생성 및 권한 설정
RUN mkdir -p /app/storage && \
    chown -R appuser:appgroup /app

# 비권한 사용자로 전환
USER appuser

# 헬스체크 추가
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 포트 노출
EXPOSE 8080

# 볼륨 마운트 포인트
VOLUME ["/app/storage"]

# 바이너리 실행
ENTRYPOINT ["/app/server"]

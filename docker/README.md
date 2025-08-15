# Docker 사용법

이 문서는 Go File Server를 Docker로 실행하는 방법을 설명합니다.

## 목차

- [빠른 시작](#빠른-시작)
- [Docker 빌드](#docker-빌드)
- [Docker Compose](#docker-compose)
- [환경 변수](#환경-변수)
- [볼륨 관리](#볼륨-관리)
- [프로덕션 배포](#프로덕션-배포)
- [문제 해결](#문제-해결)

## 빠른 시작

### 1. Docker Compose를 사용한 실행 (권장)

```bash
# 개발 환경용
cd docker
docker-compose -f docker-compose.dev.yml up -d

# 또는 프로덕션 환경용
docker-compose up -d
```

### 2. Docker만 사용한 실행

```bash
# 이미지 빌드
docker build -f docker/Dockerfile -t go-file-server .

# 컨테이너 실행
docker run -d \
  --name go-file-server \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  go-file-server
```

## Docker 빌드

### 멀티스테이지 빌드 특징

- **빌드 스테이지**: Go 1.22 Alpine에서 바이너리 컴파일
- **실행 스테이지**: 경량 Alpine Linux 기반 (약 15MB)
- **보안**: 비권한 사용자(appuser)로 실행
- **최적화**: 정적 링킹, 바이너리 스트리핑

### 수동 빌드

```bash
# 기본 빌드
docker build -f docker/Dockerfile -t go-file-server .

# 태그와 함께 빌드
docker build -f docker/Dockerfile -t go-file-server:v1.0.0 .

# 빌드 인자 사용 (필요시)
docker build -f docker/Dockerfile \
  --build-arg GO_VERSION=1.22 \
  -t go-file-server .
```

## Docker Compose

### 사용 가능한 구성

1. **개발 환경** (`docker-compose.dev.yml`)

   - 로컬 storage 디렉토리 마운트
   - 빠른 헬스체크
   - 간단한 네트워크 설정

2. **프로덕션 환경** (`docker-compose.yml`)
   - Named volume 사용
   - Nginx 리버스 프록시 (선택사항)
   - 보안 강화 설정

### 명령어

```bash
# 개발 환경 시작
docker-compose -f docker-compose.dev.yml up -d

# 프로덕션 환경 시작
docker-compose up -d

# Nginx와 함께 시작 (프로덕션)
docker-compose --profile production up -d

# 로그 확인
docker-compose logs -f go-file-server

# 서비스 중지
docker-compose down

# 볼륨까지 제거
docker-compose down -v
```

## 환경 변수

| 변수명     | 기본값         | 설명                    |
| ---------- | -------------- | ----------------------- |
| `PORT`     | `8080`         | 서버 포트               |
| `ROOT_DIR` | `/app/storage` | 파일 저장 루트 디렉토리 |
| `TZ`       | `Asia/Seoul`   | 타임존 설정             |

### 사용 예시

```bash
# Docker run에서
docker run -d \
  -e PORT=9090 \
  -e ROOT_DIR=/data \
  -p 9090:9090 \
  go-file-server

# Docker Compose에서
environment:
  - PORT=9090
  - ROOT_DIR=/data
  - TZ=UTC
```

## 볼륨 관리

### 1. Named Volume (권장)

```yaml
volumes:
  - file_storage:/app/storage
```

### 2. 바인드 마운트

```yaml
volumes:
  - ./storage:/app/storage
```

### 3. 임시 볼륨

```yaml
volumes:
  - /app/storage # 컨테이너 삭제 시 데이터 손실
```

### 볼륨 백업

```bash
# 백업 생성
docker run --rm \
  -v go-file-server_file_storage:/source:ro \
  -v $(pwd):/backup \
  alpine tar czf /backup/backup-$(date +%Y%m%d-%H%M%S).tar.gz -C /source .

# 복원
docker run --rm \
  -v go-file-server_file_storage:/target \
  -v $(pwd):/backup \
  alpine tar xzf /backup/backup-20240101-120000.tar.gz -C /target
```

## 프로덕션 배포

### 1. 기본 프로덕션 설정

```bash
# 프로덕션 환경 시작
docker-compose up -d

# 상태 확인
docker-compose ps
docker-compose logs -f
```

### 2. Nginx 리버스 프록시 사용

```bash
# Nginx와 함께 시작
docker-compose --profile production up -d

# SSL 인증서 설정 (선택사항)
mkdir -p ssl
# SSL 인증서를 ssl/ 디렉토리에 배치하고 nginx.conf 수정
```

### 3. 모니터링

```bash
# 헬스체크 상태 확인
docker inspect go-file-server | grep -A5 "Health"

# 리소스 사용량 모니터링
docker stats go-file-server

# 컨테이너 로그
docker logs -f go-file-server
```

### 4. 자동 재시작 설정

```yaml
services:
  go-file-server:
    restart: unless-stopped # 또는 'always'
```

## 문제 해결

### 일반적인 문제

1. **포트 충돌**

   ```bash
   # 다른 포트 사용
   docker run -p 8081:8080 go-file-server
   ```

2. **권한 문제**

   ```bash
   # storage 디렉토리 권한 확인
   sudo chown -R 1001:1001 ./storage
   ```

3. **메모리 부족**
   ```bash
   # 메모리 제한 설정
   docker run --memory=512m go-file-server
   ```

### 디버깅

```bash
# 컨테이너 내부 접근
docker exec -it go-file-server sh

# 헬스체크 수동 테스트
curl http://localhost:8080/health

# 상세 로그 확인
docker logs --details go-file-server
```

### 성능 최적화

1. **리소스 제한 설정**

   ```yaml
   deploy:
     resources:
       limits:
         memory: 256M
         cpus: "0.5"
       reservations:
         memory: 128M
         cpus: "0.25"
   ```

2. **로그 로테이션**
   ```yaml
   logging:
     driver: "json-file"
     options:
       max-size: "10m"
       max-file: "3"
   ```

## API 사용 예시

서버가 실행된 후 다음과 같이 API를 사용할 수 있습니다:

```bash
# 헬스체크
curl http://localhost:8080/health

# 파일 목록 조회
curl http://localhost:8080/api/files

# 파일 업로드
curl -X POST -F "file=@example.txt" http://localhost:8080/api/files

# 파일 다운로드
curl -O http://localhost:8080/api/files/example.txt

# 파일 삭제
curl -X DELETE http://localhost:8080/api/files/example.txt
```

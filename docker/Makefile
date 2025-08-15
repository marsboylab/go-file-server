# Go File Server Docker Makefile

# 변수 설정
IMAGE_NAME := go-file-server
IMAGE_TAG := latest
CONTAINER_NAME := go-file-server
PORT := 8080

# 기본 타겟
.PHONY: help
help: ## 사용 가능한 명령어를 보여줍니다
	@echo "Go File Server Docker Commands:"
	@echo "==============================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 빌드 관련
.PHONY: build
build: ## Docker 이미지를 빌드합니다
	docker build -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) ..

.PHONY: build-no-cache
build-no-cache: ## 캐시 없이 Docker 이미지를 빌드합니다
	docker build --no-cache -f Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) ..

# 실행 관련
.PHONY: run
run: ## Docker 컨테이너를 실행합니다 (개발용)
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		-v $$(pwd)/../storage:/app/storage \
		-e PORT=8080 \
		-e ROOT_DIR=/app/storage \
		$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: run-prod
run-prod: ## Docker 컨테이너를 실행합니다 (프로덕션용, named volume)
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		-v go-file-server-storage:/app/storage \
		-e PORT=8080 \
		-e ROOT_DIR=/app/storage \
		--restart unless-stopped \
		$(IMAGE_NAME):$(IMAGE_TAG)

# Docker Compose 관련
.PHONY: up
up: ## docker-compose로 서비스를 시작합니다 (개발용)
	docker-compose -f docker-compose.dev.yml up -d

.PHONY: up-prod
up-prod: ## docker-compose로 서비스를 시작합니다 (프로덕션용)
	docker-compose up -d

.PHONY: up-with-nginx
up-with-nginx: ## Nginx와 함께 서비스를 시작합니다
	docker-compose --profile production up -d

.PHONY: down
down: ## docker-compose 서비스를 중지하고 제거합니다
	docker-compose -f docker-compose.dev.yml down
	docker-compose down

.PHONY: down-volumes
down-volumes: ## docker-compose 서비스와 볼륨을 모두 제거합니다
	docker-compose -f docker-compose.dev.yml down -v
	docker-compose down -v

# 관리 관련
.PHONY: logs
logs: ## 컨테이너 로그를 확인합니다
	docker logs -f $(CONTAINER_NAME)

.PHONY: logs-compose
logs-compose: ## docker-compose 로그를 확인합니다
	docker-compose logs -f

.PHONY: shell
shell: ## 컨테이너 내부 쉘에 접속합니다
	docker exec -it $(CONTAINER_NAME) sh

.PHONY: stop
stop: ## 컨테이너를 중지합니다
	docker stop $(CONTAINER_NAME)

.PHONY: start
start: ## 중지된 컨테이너를 시작합니다
	docker start $(CONTAINER_NAME)

.PHONY: restart
restart: ## 컨테이너를 재시작합니다
	docker restart $(CONTAINER_NAME)

.PHONY: rm
rm: ## 컨테이너를 제거합니다
	docker rm -f $(CONTAINER_NAME)

# 정리 관련
.PHONY: clean
clean: ## 컨테이너와 이미지를 정리합니다
	docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) 2>/dev/null || true

.PHONY: clean-all
clean-all: ## 모든 Docker 리소스를 정리합니다 (주의!)
	docker system prune -af
	docker volume prune -f

# 상태 확인
.PHONY: status
status: ## 컨테이너 상태를 확인합니다
	@echo "=== Container Status ==="
	docker ps -a --filter name=$(CONTAINER_NAME)
	@echo ""
	@echo "=== Image Info ==="
	docker images $(IMAGE_NAME)
	@echo ""
	@echo "=== Volume Info ==="
	docker volume ls | grep file-server || echo "No volumes found"

.PHONY: health
health: ## 서버 헬스체크를 수행합니다
	@echo "Checking server health..."
	@curl -f http://localhost:$(PORT)/health && echo " ✓ Server is healthy" || echo " ✗ Server is not responding"

# 개발 워크플로우
.PHONY: dev
dev: build up logs-compose ## 개발 환경을 빌드하고 시작합니다

.PHONY: prod
prod: build up-prod ## 프로덕션 환경을 빌드하고 시작합니다

.PHONY: rebuild
rebuild: down clean build up ## 전체 재빌드 및 재시작

# 백업/복원
.PHONY: backup
backup: ## 파일 저장소를 백업합니다
	@echo "Creating backup..."
	docker run --rm \
		-v go-file-server_file_storage:/source:ro \
		-v $$(pwd):/backup \
		alpine tar czf /backup/backup-$$(date +%Y%m%d-%H%M%S).tar.gz -C /source .
	@echo "Backup completed!"

.PHONY: list-backups
list-backups: ## 백업 파일 목록을 보여줍니다
	@ls -la backup-*.tar.gz 2>/dev/null || echo "No backup files found"

# 도움말을 기본 타겟으로 설정
.DEFAULT_GOAL := help

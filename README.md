# go-file-server
go file server using gPRC
# Go File Server

간단한 파일 업로드/다운로드 API 서버입니다. `ROOT_DIR` 하위에서 파일을 나열, 업로드, 다운로드, 삭제할 수 있습니다.

## 실행

```bash
go mod tidy
go run ./cmd/server
```

환경변수:

- `PORT`: 기본 `8080`
- `ROOT_DIR`: 기본 `./storage`

## API

- `GET /health` 헬스체크
- `GET /api/files?path=.` 파일/디렉토리 목록
- `POST /api/files` 멀티파트 업로드 (`file` 필드, 선택적 `path`)
- `GET /api/files/{path/to/file}` 파일 다운로드
- `DELETE /api/files/{path/to/file}` 파일 삭제

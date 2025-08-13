package server

import (
	"net/http"
)

// OpenAPI 3.0 spec with examples for requests and responses
func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	const specJSON = `{
		"openapi": "3.0.3",
		"info": {"title": "Go File Server API", "version": "1.0.0", "description": "간단한 파일 서버 API. 디렉터리 목록 조회, 파일 업로드/다운로드/삭제를 제공합니다."},
		"servers": [{"url": "http://localhost:8080"}],
		"tags": [{"name": "files", "description": "파일 작업"}],
		"paths": {
			"/health": {"get": {"summary": "헬스체크", "operationId": "getHealth", "responses": {"200": {"description": "ok", "content": {"text/plain": {"examples": {"ok": {"value": "ok"}}}}}}}},
			"/api/files": {
				"get": {"tags": ["files"], "summary": "디렉터리 내 파일/폴더 목록", "operationId": "listFiles", "parameters": [{"name": "path", "in": "query", "required": false, "description": "조회할 상대 경로 (기본 '.')", "schema": {"type": "string", "default": "."}, "example": "."}], "responses": {"200": {"description": "목록", "content": {"application/json": {"schema": {"type": "array", "items": {"$ref": "#/components/schemas/FileInfo"}}, "examples": {"root": {"summary": "루트 목록", "value": [{"name": "subdir", "size": 0, "modTime": "2024-06-01T12:00:00Z", "relativePath": "subdir", "isDir": true}, {"name": "readme.txt", "size": 128, "modTime": "2024-06-01T12:10:00Z", "relativePath": "readme.txt", "isDir": false}]}}}}}, "400": {"description": "잘못된 경로", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "examples": {"badPath": {"value": {"error": "path escapes root"}}}}}}}},
				"post": {"tags": ["files"], "summary": "파일 업로드", "operationId": "uploadFiles", "requestBody": {"required": true, "content": {"multipart/form-data": {"schema": {"type": "object", "required": ["file"], "properties": {"file": {"type": "array", "items": {"type": "string", "format": "binary"}}, "path": {"type": "string", "example": "subdir"}}}, "encoding": {"file": {"style": "form", "explode": true}}, "examples": {"single": {"summary": "단일 파일 업로드", "value": {"path": "subdir"}}}}}}, "responses": {"201": {"description": "업로드 성공", "content": {"application/json": {"schema": {"type": "array", "items": {"$ref": "#/components/schemas/FileInfo"}}, "examples": {"saved": {"value": [{"name": "a.jpg", "size": 10240, "modTime": "2024-06-01T13:00:00Z", "relativePath": "subdir/a.jpg", "isDir": false}]}}}}}, "400": {"description": "요청 오류", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "examples": {"missing": {"value": {"error": "file field is required"}}}}}}}}
			},
			"/api/files/{filePath}": {
				"parameters": [{"name": "filePath", "in": "path", "required": true, "description": "다운로드/삭제할 파일의 상대 경로", "schema": {"type": "string"}, "example": "subdir/a.jpg"}],
				"get": {"tags": ["files"], "summary": "파일 다운로드", "operationId": "downloadFile", "responses": {"200": {"description": "바이너리 파일", "content": {"application/octet-stream": {}}}, "404": {"description": "파일 없음", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "examples": {"notFound": {"value": {"error": "file not found"}}}}}}}},
				"delete": {"tags": ["files"], "summary": "파일 삭제", "operationId": "deleteFile", "responses": {"204": {"description": "삭제됨"}, "404": {"description": "파일 없음", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "examples": {"notFound": {"value": {"error": "file not found"}}}}}}}}
			}
		},
		"components": {"schemas": {"FileInfo": {"type": "object", "required": ["name", "size", "modTime", "relativePath", "isDir"], "properties": {"name": {"type": "string", "example": "a.jpg"}, "size": {"type": "integer", "format": "int64", "example": 10240}, "modTime": {"type": "string", "format": "date-time", "example": "2024-06-01T13:00:00Z"}, "relativePath": {"type": "string", "example": "subdir/a.jpg"}, "isDir": {"type": "boolean", "example": false}}}, "ErrorResponse": {"type": "object", "required": ["error"], "properties": {"error": {"type": "string", "example": "file not found"}}}}}
	}`

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(specJSON))
}



# NeoNeuro Backend (Gin)

与项目规划一致：提供对话、情感、TTS、人格等 API；前后端解耦，前端（Live2D/Vue）通过 HTTP/WS 调用。

## Endpoints (v1)
- `GET    /api/v1/healthz`
- `POST   /api/v1/chat`        入参：`{ "text": "你好", "persona_id": "default" }`
- `POST   /api/v1/tts`         入参：`{ "text": "早上好", "voice_id": "meiko" }`
- `POST   /api/v1/emotion`     入参：`{ "text": "我有点难过" }`
- `GET    /api/v1/personas`    列表
- `POST   /api/v1/personas`    新建（最小字段）

> 当前均为 **可运行的占位实现**，已分层（handler/service/repo），方便替换为真实 LLM/TTS/情感模型。

## Run
```bash
cd backend
go run ./cmd/server
# 环境变量：HTTP_ADDR (默认 :8080), APP_ENV (dev/prod)
```

## Test
```bash
go test ./...
```

# WhatsApp Provider - cURL Commands for Postman

## 1. Provedores

### Listar Provedores
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/providers \
  -H "Content-Type: application/json"
```

## 2. Instâncias

### Criar Instância
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/instances \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Minha Instância Z-API",
    "provider": "z-api",
    "instance_id": "3E68B22C8B61603262EC967D54735262",
    "token": "657054F1246E3A6A2049CD9E",
    "config": {}
  }'
```

### Listar Instâncias
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/instances \
  -H "Content-Type: application/json"
```

### Obter Instância por ID
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/instances/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json"
```

### Status da Instância
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/status/SEU_TOKEN_Z_API \
  -H "Content-Type: application/json"
```

### Deletar Instância
```bash
curl -X DELETE \
  http://localhost:8080/api/v1/whatsapp/instances/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json"
```

## 3. Mensagens

### Enviar Mensagem de Texto
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/messages \
  -H "Content-Type: application/json" \
  -d '{
    "instance_id": "SEU_TOKEN_Z_API",
    "phone": "5511999999999",
    "type": "text",
    "content": "Olá! Mensagem de teste."
  }'
```

### Enviar Imagem
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/messages \
  -H "Content-Type: application/json" \
  -d '{
    "instance_id": "SEU_TOKEN_Z_API",
    "phone": "5511999999999",
    "type": "image",
    "content": "Legenda da imagem",
    "media_url": "https://picsum.photos/800/600"
  }'
```

### Enviar Vídeo
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/messages \
  -H "Content-Type: application/json" \
  -d '{
    "instance_id": "SEU_TOKEN_Z_API",
    "phone": "5511999999999",
    "type": "video",
    "content": "Legenda do vídeo",
    "media_url": "https://sample-videos.com/zip/10/mp4/SampleVideo_1280x720_1mb.mp4"
  }'
```

### Enviar Áudio
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/messages \
  -H "Content-Type: application/json" \
  -d '{
    "instance_id": "SEU_TOKEN_Z_API",
    "phone": "5511999999999",
    "type": "audio",
    "content": "",
    "media_url": "https://www.soundjay.com/misc/sounds/bell-ringing-05.wav"
  }'
```

### Enviar Documento
```bash
curl -X POST \
  http://localhost:8080/api/v1/whatsapp/messages \
  -H "Content-Type: application/json" \
  -d '{
    "instance_id": "SEU_TOKEN_Z_API",
    "phone": "5511999999999",
    "type": "document",
    "content": "Documento em anexo",
    "media_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
  }'
```

### Obter Mensagem por ID
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/messages/456e7890-e89b-12d3-a456-426614174001 \
  -H "Content-Type: application/json"
```

### Histórico de Mensagens da Instância
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/messages/instance/SEU_TOKEN_Z_API \
  -H "Content-Type: application/json"
```

### Histórico com Paginação
```bash
curl -X GET \
  http://localhost:8080/api/v1/whatsapp/messages/instance/SEU_TOKEN_Z_API?limit=10&offset=0 \
  -H "Content-Type: application/json"
```

## 4. Monitoramento

### Health Check
```bash
curl -X GET \
  http://localhost:8080/health \
  -H "Content-Type: application/json"
```

### Welcome
```bash
curl -X GET \
  http://localhost:8080/api/v1/ \
  -H "Content-Type: application/json"
```

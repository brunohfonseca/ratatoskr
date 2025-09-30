┌──────────────────────────────────────────────────────────────┐
│ 1. USUÁRIO CRIA ENDPOINT                                     │
└──────────────┬───────────────────────────────────────────────┘
│
↓
┌──────────────────────────────────────────────────────────────┐
│ 2. SERVICE                                                    │
│    ├─ Salva no Postgres: endpoints table                     │
│    └─ Envia pro Redis: stream "endpoints"                    │
└──────────────┬───────────────────────────────────────────────┘
│
↓
┌──────────────────────────────────────────────────────────────┐
│ REDIS STREAM "endpoints"                                      │
│ [msg1] [msg2] [msg3] ← Fila de endpoints para processar     │
└──────────────┬───────────────────────────────────────────────┘
│
↓
┌──────────────────────────────────────────────────────────────┐
│ 3. WORKER (loop infinito)                                    │
│    XReadGroup() → Lê mensagem da fila                        │
└──────────────┬───────────────────────────────────────────────┘
│
↓
┌──────────────────────────────────────────────────────────────┐
│ 4. PROCESSA MENSAGEM                                          │
│    ├─ Extrai: uuid, domain, path, check_ssl                  │
│    ├─ Faz health check: GET https://domain/path              │
│    ├─ Resultado: online/offline + tempo                      │
│    ├─ XAck() → Marca como processado                         │
│    └─ XDel() → Remove da fila                                │
└───────────────────────────────────────────────────────────────┘
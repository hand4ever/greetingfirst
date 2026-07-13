# API Contract: SHA256 Demo

**Endpoint**: `GET /demo/sha256`

## Request

```http
GET /demo/sha256?text={input_text}
```

| Parameter | Type | Required | Location | Description |
|-----------|------|----------|----------|-------------|
| text | string | Yes | Query string | Input text to compute SHA256 hash for |

**Example**:
```http
GET /demo/sha256?text=hello
```

## Response

### Success (200)

```json
{
  "code": 0,
  "message": "",
  "data": {
    "input": "hello",
    "hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
  },
  "trace_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "cost": "xxms"
}
```

| Field | Type | Description |
|-------|------|-------------|
| data.input | string | Original input text echoed back |
| data.hash | string | SHA256 hash in lowercase hex (64 characters) |
| code | integer | Always 0 on success |

### Error — Missing Parameter (400)

```json
{
  "code": 1,
  "message": "text parameter is required",
  "data": null,
  "trace_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "cost": "xxms"
}
```

### Error — Wrong Method (405)

Returned when using POST, PUT, DELETE etc.

```json
{
  "code": 1,
  "message": "method not allowed",
  "data": null,
  "trace_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "cost": "xxms"
}
```

## Behavior Notes

- Accepts any UTF-8 encoded text including emoji, special characters, Chinese characters
- Empty string (`text=`) is valid and returns hash of empty string
- Hash output is always 64 lowercase hexadecimal characters
- No authentication required (demo endpoint)
- Maximum input size limited only by Go runtime memory

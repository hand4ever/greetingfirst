# Quickstart: SHA256 Demo API

**Feature**: SHA256 Demo API | **Date**: 2026-07-13

## Prerequisites

- Go ≥ 1.22 installed
- Project dependencies resolved (`go mod tidy`)

## Setup

```bash
# Start the server
go run main.go
```

## Validation Scenarios

### 1. Basic Hash Computation

**Command**:
```bash
curl -s "http://localhost:1323/demo/sha256?text=hello" | jq
```

**Expected**:
```json
{
  "code": 0,
  "data": {
    "input": "hello",
    "hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
  }
}
```

**Cross-verify**:
```bash
echo -n "hello" | sha256sum
# Should match the hash field above
```

### 2. Chinese Text

**Command**:
```bash
curl -s "http://localhost:1323/demo/sha256?text=%E4%BD%A0%E5%A5%BD" | jq
```

**Expected**: `data.hash` is 64 lowercase hex characters, `data.input` is "你好".

### 3. Empty String

**Command**:
```bash
curl -s "http://localhost:1323/demo/sha256?text=" | jq
```

**Expected**:
```json
{
  "code": 0,
  "data": {
    "input": "",
    "hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  }
}
```

### 4. Missing Parameter (Error)

**Command**:
```bash
curl -s -w "\nHTTP %{http_code}" "http://localhost:1323/demo/sha256"
```

**Expected**: HTTP 400 with error message indicating missing text parameter.

### 5. Wrong HTTP Method (Error)

**Command**:
```bash
curl -s -w "\nHTTP %{http_code}" -X POST "http://localhost:1323/demo/sha256?text=test"
```

**Expected**: HTTP 405 Method Not Allowed.

### 6. Special Characters

**Command**:
```bash
curl -s "http://localhost:1323/demo/sha256?text=%F0%9F%98%80" | jq
```

**Expected**: Returns correct hash for emoji character 😀.

## Run Unit Tests

```bash
go test -v ./handler/ -run TestSha256 -count=1
```

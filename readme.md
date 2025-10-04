# Kirana AI Backend

A lightweight Go backend service built using **Gin**, **Zap**, and **Viper**, following clean and modular design principles based on the **Hexagonal Architecture**. This backend powers **kirana AI**, a Malayalam voice-based assistant for merchants in **Kerala**.

---

## 🌐 Supported Languages

- Malayalam

> *More to be added soon*

---

## 🚀 Overview

The backend currently supports uploading audio files from a client application. The uploaded audio can be stored locally or in **AWS S3**, based on the configuration set in `config.yaml`. Future stages will include Malayalam speech transcription, guardrail filtering, and AI inference pipelines.

---

## 🧱 Project Structure

```
cmd/server/main.go        # Entry point
config/config.yaml        # Centralized configuration file (YAML)
internal/adapter/http     # HTTP handler (Gin)
internal/adapter/storage  # Local and S3 storage adapters
internal/service          # Core business logic
internal/models           # API response models
internal/port             # Storage interface definitions
internal/pkg/logger       # Zap logging setup
internal/pkg/config       # Viper configuration loader
```

---

## ⚙️ Endpoint

### **POST /v1/audio**

Uploads an audio file to the server. Depending on configuration, it will either be stored locally under `/uploads` or uploaded directly to an **S3 bucket**.

#### Request

```bash
curl -X POST http://localhost:8080/v1/audio \
  -F 'file=@"/path/to/audio/file.m4a"'
```

#### Successful Response

```json
{
  "status": "STATUS_CREATED",
  "message": "file uploaded successfully",
  "status_code": 201,
  "data": {
    "file": "20251004_115155_Elavakattumoola Road.m4a"
  }
}
```

#### Error Response

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "file is required",
    "status_code": 400
  }
}
```

---

## ⚙️ Configuration

The application configuration is managed using **Viper** and stored in `config/config.yaml`.

Example:

```yaml
app:
  port: "8080"

aws:
  use_s3: true
  bucket: "kirana-ai-audio"
  region: "ap-south-1"

logging:
  level: "info"
```

* Change `use_s3` to `false` to use local storage.
* The AWS SDK automatically picks credentials from your AWS CLI or environment.

---

## 🧩 Tech Stack

* **Language:** Go 1.22+
* **Framework:** Gin
* **Logging:** Uber Zap
* **Configuration:** Viper (YAML-based)
* **Cloud Storage:** AWS S3

---

## 🧠 Author

**Ashwin Nambiar** — Founder & Engineer, kirana AI

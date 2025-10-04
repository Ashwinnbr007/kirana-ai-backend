# Kinara AI Backend

A lightweight Go backend service built using **Gin** and **Zap**, following clean and modular design principles. This backend powers **Kinara AI**, a Malayalam voice-based assistant for merchants in **Kerala**.

---

## Supported languages

 - Malayalam

> *More to be added soon*

## 🚀 Overview

The backend currently supports uploading audio files from a client application. The uploaded audio will later flow through transcription, guardrail filtering, and AI inference pipelines.

---

## 🧱 Project Structure

```
cmd/server/main.go        # Entry point
internal/adapter/http     # HTTP handler (Gin)
internal/adapter/storage  # Local file storage
internal/service          # Core business logic
internal/models           # API response models
internal/pkg/logger       # Zap logging setup
```

---

## ⚙️ Endpoint

### **POST /v1/audio**

Uploads an audio file to the server. Current inplementation saves it directly on the server, the next plan is to use an s3 bucket.

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

## 🧩 Tech Stack

* **Language:** Go 1.22+
* **Framework:** Gin
* **Logging:** Uber Zap

---

## 🧠 Author

**Ashwin Nambiar** — Founder & Engineer, Kinara AI
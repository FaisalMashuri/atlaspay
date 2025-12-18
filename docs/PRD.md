# 📕 PRD

## **Global Wallet, Ledger, Fraud & AI Investigator Platform**

---

## 1. Product Overview

### 1.1 Product Name

**AtlasPay (working title)**

> _Global Transaction, Risk & AI Investigation Platform_

---

### 1.2 Problem Statement

Sistem finansial modern menghadapi masalah:

- Transaksi global dengan **konsistensi tinggi**

- Fraud yang semakin kompleks & cepat

- Investigasi manual lambat dan mahal

- Log & data tersebar, sulit dianalisis

- Tim kecil tapi beban operasional besar


---

### 1.3 Product Vision

Menyediakan **platform transaksi finansial global** yang:

- Aman

- Konsisten

- Observable

- Audit-friendly

- Dibantu AI untuk investigasi & operasional


---

### 1.4 Target Users

|User|Kebutuhan|
|---|---|
|End User|Transaksi cepat & aman|
|Merchant|Settlement & monitoring|
|Fraud Analyst|Deteksi & investigasi|
|Auditor|Audit trail|
|Admin|Konfigurasi & kontrol|

---

## 2. Goals & Non-Goals

### 2.1 Goals

- Menjamin **business invariant** (saldo, ledger)

- Deteksi fraud near real-time

- Audit trail 100% traceable

- AI membantu investigasi, bukan menggantikan manusia

- Siap di-scale secara horizontal


### 2.2 Non-Goals

- KYC real

- Integrasi bank production

- Real money settlement


---

## 3. User Journeys

---

### 3.1 Wallet Transaction Flow

1. User request transfer

2. API Gateway validate auth + rate limit

3. Wallet service:

    - lock logical balance

    - validate invariant

4. Ledger entry dibuat (double entry)

5. Event dipublish ke Kafka

6. Fraud engine consume event

7. Read model di-update via CQRS

8. User melihat balance update


---

### 3.2 Fraud Investigation Flow

1. Fraud engine flag transaksi

2. Alert dikirim ke Telegram / WhatsApp

3. Fraud analyst membuka dashboard

4. Analyst bertanya ke AI:

   > “Kenapa transaksi ini ditandai?”

5. AI:

    - Query ClickHouse

    - Search Elastic

    - Vector similarity

6. AI memberikan explanation + evidence


---

## 4. Functional Requirements (Product View)

---

### 4.1 Authentication & Access

- Login via OAuth2

- Token JWT RS256

- Refresh token rotation

- API key untuk system-to-system

- Role-based access (RBAC)


---

### 4.2 Wallet

- Multi-currency wallet

- Real-time balance

- Strong consistency

- Cached read (Redis)


---

### 4.3 Transactions

- Transfer internal

- Idempotent request

- Atomic execution

- Retry-safe


---

### 4.4 Ledger

- Double-entry accounting

- Immutable record

- Time-ordered (logical time)


---

### 4.5 Fraud Detection

- Rule-based engine (v1)

- Pattern similarity (vector)

- Velocity checks

- Threshold-based scoring


---

### 4.6 AI Investigator

- Natural language query

- Explain fraud reasoning

- Reference historical cases

- Provide links to raw data


---

### 4.7 Observability

- End-to-end tracing

- Error rate & latency

- Business metrics (volume, fraud rate)


---

### 4.8 Notifications

- Telegram bot

- WhatsApp API

- Severity-based alerting


---

## 5. Non-Functional Requirements

---

### 5.1 Performance

|Metric|Target|
|---|---|
|Read latency p95|< 200ms|
|Write latency p99|< 500ms|
|Fraud scoring|< 1s|

---

### 5.2 Scalability

- Stateless services

- Horizontal scaling

- Kafka partitioned by key


---

### 5.3 Reliability

- At-least-once messaging

- Idempotent consumers

- Circuit breaker & backoff


---

### 5.4 Security

- Rate limiting

- Secret rotation

- Audit logging

- Zero trust internal comms (mTLS gRPC)


---

## 6. Data & Storage Strategy

|Data|Tech|
|---|---|
|Transaction OLTP|PostgreSQL|
|Analytics|ClickHouse|
|Cache|Redis|
|Logs|Elastic|
|Vector|pgvector / Milvus|
|CDC|Debezium|
|Local analytics|DuckDB|

---

## 7. API & Integration

- API Gateway (REST)

- gRPC internal

- BFF per channel

- Webhook-ready

- Event-driven core


---

## 8. Success Metrics (KPIs)

### 8.1 Product Metrics

- Zero negative balance

- Fraud detection accuracy

- Investigation time reduction


### 8.2 Engineering Metrics

- Error rate < 0.1%

- Event lag < 2s

- Service SLO met


---

## 9. Risks & Tradeoffs

|Risk|Tradeoff|
|---|---|
|Strong consistency|Higher latency|
|CQRS complexity|More infra|
|AI hallucination|Human validation|
|Eventual read model|Slight staleness|

---

## 10. Release Phases

### Phase 1 – Core Ledger

- Wallet

- Ledger

- Transaction

- Kafka + Outbox


### Phase 2 – Fraud & Analytics

- Fraud engine

- ClickHouse

- Elastic

- CDC


### Phase 3 – AI & Channels

- LLM integration

- Telegram / WhatsApp bot

- Investigator UI


---

## 11. Open Questions

- Multi-region write?

- Event schema governance?

- Fraud rules vs ML?
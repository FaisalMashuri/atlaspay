git # 🔐 1️⃣ AUTH SERVICE — Identity Database (PostgreSQL)

## ERD

`USERS - id (UUID, PK)                ← technical ID - user_ref (ULID, UNIQUE)      ← exposed ke luar - email (UNIQUE) - password_hash - status (PENDING | ACTIVE | SUSPENDED) - created_at - updated_at  USER_OTPS - id (UUID, PK) - user_id (UUID)               ← internal join - otp_code - expires_at - used_at - created_at  REFRESH_TOKENS - id (UUID, PK) - user_id (UUID) - token_hash - expires_at - revoked_at - created_at`

## Keterangan

- `user_ref` dipakai di:

    - JWT (`sub`)

    - Event auth

- **Tidak ada relasi ke wallet / ledger**

- OTP & refresh token **append / revoke only**


---

# 🌐 2️⃣ BACKEND GATEWAY

❌ **TIDAK PUNYA DATABASE**

> Gateway adalah stateless edge.

---

# 🧠 3️⃣ TRANSACTION ORCHESTRATOR — Workflow DB (PostgreSQL)

## ERD

`TRANSACTIONS - id (UUID, PK) - transaction_ref (ULID, UNIQUE)   ← BUSINESS ID - user_ref (ULID)                  ← dari Auth - type (TRANSFER | TOPUP) - status (INIT | WALLET_OK | LEDGER_OK | COMPLETED | FAILED) - idempotency_key (UNIQUE) - created_at - updated_at  TRANSACTION_STEPS - id (UUID, PK) - transaction_id (UUID)            ← join internal - step (WALLET | LEDGER) - status (SUCCESS | FAILED) - error_reason - created_at`

## Keterangan

- `transaction_ref`:

    - Dipakai di **ledger**

    - Dipakai di **event**

    - Dipakai di **audit**

- Orchestrator **BUKAN source of truth uang**


---

# 👛 4️⃣ WALLET COMMAND SERVICE — Operational State (PostgreSQL)

## ERD

`WALLETS - id (UUID, PK) - wallet_ref (ULID, UNIQUE)     ← exposed - user_ref (ULID) - currency (IDR, USD, etc) - balance (BIGINT) - status (ACTIVE | FROZEN) - created_at - updated_at  WALLET_MUTATIONS - id (UUID, PK) - wallet_id (UUID)              ← join internal - transaction_ref (ULID)        ← business reference - mutation_type (DEBIT | CREDIT) - amount - balance_after - created_at`

## Keterangan

- `balance` = **state operasional**

- `wallet_ref`:

    - Dipakai API

    - Dipakai event

- **Tidak ada join ke ledger**


---

# 📚 5️⃣ LEDGER COMMAND SERVICE — Accounting Truth (PostgreSQL)

## ERD (Double-Entry, Immutable)

`LEDGER_ACCOUNTS - id (UUID, PK) - account_ref (ULID, UNIQUE) - owner_type (WALLET | SYSTEM) - owner_ref (ULID) - currency - created_at  LEDGER_JOURNALS - id (UUID, PK) - journal_ref (ULID, UNIQUE) - transaction_ref (ULID)        ← BUSINESS FACT - created_at  LEDGER_ENTRIES - id (UUID, PK) - journal_id (UUID)             ← join internal - account_id (UUID) - entry_type (DEBIT | CREDIT) - amount - created_at`

## Invariant (WAJIB)

- Σ(DEBIT) = Σ(CREDIT) per `journal_ref`

- **NO UPDATE**

- **NO DELETE**


## Keterangan penting

- Ledger **tidak join ke wallet table**

- Ledger hanya tahu:

    - `transaction_ref`

    - `account_ref`


---

# 📖 6️⃣ WALLET QUERY SERVICE — UX Read Model (Redis + ClickHouse)

## ClickHouse ERD

`wallet_balances_view - wallet_ref - user_ref - currency - balance - last_updated  wallet_transactions_view - transaction_ref - wallet_ref - mutation_type - amount - created_at`

## Redis

`Key: wallet:{wallet_ref}:balance`

## Keterangan

- Dibangun dari **event**

- Boleh di-drop & rebuild

- Tidak pernah write ke core DB


---

# 📖 7️⃣ LEDGER QUERY SERVICE — Audit Read Model (ClickHouse)

`ledger_entries_view - journal_ref - transaction_ref - account_ref - entry_type - amount - created_at`

## Keterangan

- Admin / auditor only

- Immutable view


---

# 📊 8️⃣ GLOBAL READ MODEL — Analytics (ClickHouse / DuckDB)

`transaction_facts - transaction_ref - user_ref - amount - currency - type - status - created_at  daily_wallet_snapshot - wallet_ref - date - opening_balance - closing_balance - total_debit - total_credit`

## Keterangan

- Cross-domain

- Dipakai BI, Fraud, AI

- **Bukan source of truth**


---

# 🚨 9️⃣ RISK / FRAUD SERVICE — Feature Store

`risk_scores - transaction_ref - risk_score - reason - created_at`

## Keterangan

- Advisory only

- Tidak block transaksi


---

# 🧠 1️⃣0️⃣ VECTOR SERVICE — Vector DB

`transaction_embeddings - transaction_ref - embedding (vector) - created_at`

---

# 📣 1️⃣1️⃣ NOTIFICATION SERVICE

`notifications - id (UUID, PK) - event_type - reference_ref (ULID)      ← transaction_ref / wallet_ref - channel (WA | Telegram) - status (SENT | FAILED) - created_at`

---

# 🧩 ATURAN EMAS ERD (INI WAJIB DIPEGANG)

### 1️⃣ JOIN

`JOIN INTERNAL  → id (UUID)`

### 2️⃣ EXPOSE

`API / EVENT / LOG → *_ref (ULID)`

### 3️⃣ LEDGER

`Ledger tidak pernah join ke DB lain`

### 4️⃣ EVENT

`Event tidak pernah kirim UUID internal`

CREATE TABLE outbox_events (
id UUID PRIMARY KEY,
aggregate_type TEXT NOT NULL,
aggregate_ref TEXT NOT NULL,
event_type TEXT NOT NULL,
payload JSONB NOT NULL,
created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
published BOOLEAN DEFAULT false
);
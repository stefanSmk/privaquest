# PrivaQuest

Self-hosted GDPR / DSGVO / RGPD **data subject request** manager for small EU teams.

When someone emails "delete my data" or "send me everything you have on me", you have **one month** to respond (Art. 12 GDPR). Most small companies track that in email threads or spreadsheets. This tool gives you a proper queue, deadlines, and an audit trail — on your own server.

Built for freelancers, agencies, and Mittelstand teams who host on Hetzner/OVH and don't want another US SaaS with DPA paperwork.

## Why this exists

I kept seeing the same pattern: EU startups know GDPR exists, but when an actual request lands in the inbox, nobody knows who owns it or whether the 30-day clock is ticking.

Legal uncertainty is the #1 blocker for digital tools in Germany (Bitkom, 2025). France has similar pressure — CNIL enforcement, cookie banners, Schrems II anxiety about US processors.

PrivaQuest doesn't replace a lawyer. It helps you **operationalize** the boring part: intake, tracking, status changes, proof you responded on time.

## Features

- Public form + API for intake (access, delete, rectify, object)
- **30-day deadline** calculated automatically
- Admin API: list, update status, dashboard (open / overdue / due this week)
- **Immutable audit log** per request
- **EN / DE / FR** out of the box
- SQLite, one binary, Docker
- No telemetry, no third-party processors

## Quick start

```bash
git clone https://github.com/stefanSmk/privaquest.git
cd privaquest
go run ./cmd/server
```

Open `http://localhost:8080` for the public form.

### Docker

```bash
docker compose up --build
```

### Submit via API

```bash
curl -X POST http://localhost:8080/api/requests \
  -H "Content-Type: application/json" \
  -d '{
    "type": "delete",
    "email": "user@example.com",
    "full_name": "Anna Müller",
    "description": "Please delete my account data",
    "locale": "de"
  }'
```

### Admin (change `ADMIN_API_KEY` in production)

```bash
curl http://localhost:8080/api/admin/dashboard \
  -H "Authorization: Bearer change-me-admin-key"
```

## Config

| Variable | Default | Notes |
|----------|---------|-------|
| `PORT` | `8080` | HTTP port |
| `ADMIN_API_KEY` | `change-me-admin-key` | Admin API auth |
| `PUBLIC_TOKEN` | _(empty)_ | Optional token for POST /api/requests |
| `DATABASE_URL` | `file:privaquest.db?...` | SQLite path |

## Who is this for?

| Market | Why it matters |
|--------|----------------|
| **Germany (DE)** | DSGVO + BDSG; Mittelstand prefers self-hosted on DE infra |
| **France (FR)** | RGPD + CNIL; French UI matters for ops teams |
| **EN** | GDPR applies to EU residents worldwide |

Details: [docs/market-research.md](./docs/market-research.md)

## Not legal advice

Organizes requests. You still need a privacy policy and real deletion/export processes.

## Other languages

- [Deutsch](./README.de.md)
- [Français](./README.fr.md)
- [Русский](./README.ru.md)

## License

MIT

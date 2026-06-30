# Market research — DE / FR / EN

Why PrivaQuest targets these markets and how to position it.

## Summary

| Signal | Germany | France | English (EU + global) |
|--------|---------|--------|------------------------|
| Regulation pressure | DSGVO + BDSG, BfDI scrutiny | RGPD + CNIL (Matomo recommended over GA) | GDPR applies to EU residents everywhere |
| SME readiness | 94% feel prepared (Qonto) but 53% cite legal uncertainty for AI/digital (Bitkom) | 48% feel **unprepared** — highest in surveyed EU | UK/US firms serving EU still need DSAR process |
| Self-hosting trend | Hetzner, Mittelstand culture | OVH, data sovereignty | r/selfhosted, privacy-first stack |
| Language | German UI for ops teams | French UI important | English = default docs + HN/Reddit |

## Germany (DE)

**Demand drivers:**
- 109k+ unfilled IT jobs (BITKOM) — backend Go/Python in demand
- Mittelstand (~99% of firms) prefers tools that integrate without replacing SAP/legacy
- **Legal uncertainty (53%)** blocks SaaS adoption — self-hosted compliance tools win trust
- BDSG requires documented processes; audit trails matter

**Where to promote:**
- Heise / Golem (if you write a technical article)
- German self-hosted communities
- LinkedIn DACH — target "Datenschutzbeauftragter", "IT-Leiter KMU"
- Hetzner community

**Messaging (DE):**
> DSGVO-Anfragen ohne Excel: Self-hosted Queue mit Frist und Audit-Log. Läuft auf Ihrem Hetzner-Server.

## France (FR)

**Demand drivers:**
- CNIL enforcement on analytics/cookies — privacy tooling is mainstream conversation
- French SMEs least prepared digitally in Qonto survey
- Schrems II → distrust of US-only processors
- Symfony/PHP strong locally, but Go backend still valued for infra tools

**Where to promote:**
- CNIL-adjacent communities (careful: not legal advice)
- French dev Twitter/X, r/france tech threads
- OVH ecosystem

**Messaging (FR):**
> Demandes RGPD : file d'attente self-hosted, délai 30 jours, journal d'audit. Hébergement EU.

## English (EN)

**Demand drivers:**
- Largest reach for GitHub stars/clones
- Show HN, r/selfhosted, r/gdpr
- UK GDPR post-Brexit still similar DSAR rules
- Remote EU companies operate in English internally

**Messaging (EN):**
> Self-hosted GDPR request queue for small teams. 30-day deadlines, audit log, no US SaaS.

## Competitive landscape

| Tool | Gap PrivaQuest fills |
|------|----------------------|
| Spreadsheets / email | No deadline tracking, no audit trail |
| OneTrust, TrustArc | Enterprise $$$, overkill for 5-person team |
| Generic ticketing (Zendesk) | Not GDPR-specific, no 30-day SLA baked in |
| URL shorteners / CLI toys | Wrong problem |

Not saturated like "another URL shortener". Niche is smaller but **buyer intent is real** (compliance pain).

## GitHub About (copy-paste)

**EN:**
```
Self-hosted GDPR/DSGVO/RGPD data subject request manager. 30-day deadlines, audit log, EN/DE/FR. For EU SMEs on their own infra.
```

**DE:**
```
Self-hosted DSGVO-Anfragenverwaltung. Fristen, Audit-Log, DE/EN/FR. Für KMU auf eigenem Server (Hetzner/OVH).
```

**FR:**
```
Gestion self-hosted des demandes RGPD. Délais 30 jours, audit, FR/EN/DE. Pour PME sur infra EU.
```

## Topics

```
gdpr, dsgvo, rgpd, privacy, compliance, self-hosted, go, sqlite, docker, audit-log, data-subject-request, dsar, cnil, mittelstand
```

## Sources (2025–2026)

- BITKOM IT labour market / AI adoption surveys
- Qonto European SME digital readiness report
- CNIL guidance on analytics alternatives
- EU GDPR Art. 12–17 (response timelines and request types)

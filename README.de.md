# PrivaQuest

Self-hosted **DSGVO-/RGPD-Anfragenverwaltung** für kleine Teams in der EU.

Wenn jemand schreibt „Löscht meine Daten" oder „Schickt mir alles, was ihr über mich habt", habt ihr **einen Monat** Zeit (Art. 12 DSGVO). Viele KMU tracken das in E-Mail-Threads. PrivaQuest gibt euch eine Queue, Fristen und ein Audit-Log — auf eurem Server.

Für Freelancer, Agenturen und Mittelstand, die auf Hetzner/OVH hosten und kein US-SaaS mit DPA-Pingpong wollen.

## Warum

Rechtliche Unsicherheit ist laut Bitkom 2025 der größte Blocker für Digitalisierung in deutschen Firmen (~53%). Die DSGVO kennen alle — aber wenn die erste Anfrage im Postfach landet, fehlt oft der Prozess.

PrivaQuest ersetzt keinen Anwalt. Es hilft beim Operativen: Eingang, Status, Nachweis, dass ihr fristgerecht geantwortet habt.

## Funktionen

- Öffentliches Formular + API (Auskunft, Löschung, Berichtigung, Widerspruch)
- **30-Tage-Frist** automatisch
- Admin-API + Dashboard (offen / überfällig / diese Woche fällig)
- **Audit-Log** pro Anfrage
- **DE / EN / FR**
- SQLite, ein Binary, Docker

## Start

```bash
go run ./cmd/server
```

Formular: `http://localhost:8080`

```bash
curl -X POST http://localhost:8080/api/requests \
  -H "Content-Type: application/json" \
  -d '{"type":"delete","email":"nutzer@beispiel.de","locale":"de"}'
```

## Kein Rechtsrat

Software zur Organisation — Datenschutzerklärung und echte Löschprozesse braucht ihr trotzdem.

## Weitere Sprachen

- [English](./README.md)
- [Français](./README.fr.md)

## Lizenz

MIT

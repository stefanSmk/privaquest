# PrivaQuest

Self-hosted менеджер **GDPR/DSGVO/RGPD-запросов** для маленьких EU-команд.

Когда клиент пишет «удалите мои данные» — у вас **месяц** на ответ (Art. 12 GDPR). Многие KMU ведут это в почте. PrivaQuest — очередь, дедлайны, audit log на **вашем** сервере.

## Зачем

В Германии ~53% компаний называют правовую неопределённость главным барьером цифровизации (Bitkom). Во Франции ~48% SME чувствуют себя не готовы (Qonto). Self-hosted + EU hosting = меньше DPA/Schrems II головной боли.

Не заменяет юриста. Помогает не пропустить срок.

## Языки

EN / DE / FR из коробки.

## Запуск

```bash
go run ./cmd/server
```

Форма: `http://localhost:8080`

## Другие языки

- [English](./README.md)
- [Deutsch](./README.de.md)
- [Français](./README.fr.md)

## Лицензия

MIT

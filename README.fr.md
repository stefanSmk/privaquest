# PrivaQuest

Gestion **self-hosted des demandes RGPD** pour petites équipes en Europe.

Quand quelqu'un écrit « supprimez mes données » ou « envoyez-moi tout ce que vous avez », vous avez **un mois** pour répondre (Art. 12 RGPD). Beaucoup de PME gèrent ça par e-mail. PrivaQuest ajoute une file d'attente, des échéances et une piste d'audit — sur votre serveur.

Pour freelances et PME qui hébergent chez OVH/Hetzner et veulent éviter un SaaS US + DPA.

## Pourquoi

En France, ~48 % des PME se disent mal préparées au numérique (Qonto, 2025). La CNIL pousse à la conformité ; les outils US posent des questions Schrems II.

PrivaQuest ne remplace pas un DPO avocat. Il structure l'opérationnel : réception, statut, preuve de réponse dans les délais.

## Fonctions

- Formulaire public + API (accès, suppression, rectification, opposition)
- **Délai 30 jours** automatique
- API admin + tableau de bord
- **Journal d'audit** immuable
- **FR / EN / DE**
- SQLite, Docker

## Démarrage

```bash
go run ./cmd/server
```

Formulaire : `http://localhost:8080`

## Pas un conseil juridique

Organisation seulement — politique de confidentialité et vrais processus de suppression restent nécessaires.

## Projets associés

- [RopaDesk](https://github.com/stefanSmk/ropadesk) — registre des activités (Art. 30)
- [CookieAudit](https://github.com/stefanSmk/cookieaudit) — scanner cookies et trackers

## Autres langues

- [English](./README.md)
- [Deutsch](./README.de.md)

## Licence

MIT

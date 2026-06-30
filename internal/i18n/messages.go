package i18n

import "strings"

var messages = map[string]map[string]string{
	"en": {
		"request_received":  "Your privacy request has been received.",
		"reference_label":   "Reference number",
		"due_label":         "Response due by",
		"type_access":       "Access my data (Art. 15 GDPR)",
		"type_delete":       "Delete my data (Art. 17 GDPR)",
		"type_rectify":      "Correct my data (Art. 16 GDPR)",
		"type_object":       "Object to processing (Art. 21 GDPR)",
		"status_open":       "Open",
		"status_in_progress": "In progress",
		"status_resolved":   "Resolved",
		"status_rejected":   "Rejected",
	},
	"de": {
		"request_received":  "Ihre Datenschutzanfrage wurde empfangen.",
		"reference_label":   "Referenznummer",
		"due_label":         "Antwort fällig bis",
		"type_access":       "Auskunft über meine Daten (Art. 15 DSGVO)",
		"type_delete":       "Löschung meiner Daten (Art. 17 DSGVO)",
		"type_rectify":      "Berichtigung meiner Daten (Art. 16 DSGVO)",
		"type_object":       "Widerspruch gegen Verarbeitung (Art. 21 DSGVO)",
		"status_open":       "Offen",
		"status_in_progress": "In Bearbeitung",
		"status_resolved":   "Erledigt",
		"status_rejected":   "Abgelehnt",
	},
	"fr": {
		"request_received":  "Votre demande relative à la protection des données a été reçue.",
		"reference_label":   "Numéro de référence",
		"due_label":         "Réponse attendue avant le",
		"type_access":       "Accès à mes données (Art. 15 RGPD)",
		"type_delete":       "Suppression de mes données (Art. 17 RGPD)",
		"type_rectify":      "Rectification de mes données (Art. 16 RGPD)",
		"type_object":       "Opposition au traitement (Art. 21 RGPD)",
		"status_open":       "Ouverte",
		"status_in_progress": "En cours",
		"status_resolved":   "Traitée",
		"status_rejected":   "Rejetée",
	},
}

func T(locale, key string) string {
	locale = Normalize(locale)
	if m, ok := messages[locale]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return messages["en"][key]
}

func Normalize(locale string) string {
	locale = strings.ToLower(strings.TrimSpace(locale))
	if strings.HasPrefix(locale, "de") {
		return "de"
	}
	if strings.HasPrefix(locale, "fr") {
		return "fr"
	}
	return "en"
}

func SupportedLocales() []string {
	return []string{"en", "de", "fr"}
}

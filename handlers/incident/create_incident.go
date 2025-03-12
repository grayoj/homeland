package incident

import (
	"context"
	"encoding/json"
	"homeland/models"
	"homeland/utils"
	"net/http"

	"github.com/uptrace/bun"
)

type CreateIncidentRequest struct {
	AgentID           string                  `json:"agent_id"`
	Department        models.DepartmentEnum   `json:"department"`
	IncidentType      models.IncidentTypeEnum `json:"incident_type"`
	Severity          models.SeverityEnum     `json:"severity"`
	CallerFullName    string                  `json:"caller_full_name"`
	CallerPhoneNumber string                  `json:"caller_phone_number"`
	CallerLocation    string                  `json:"caller_location"`
	PeopleInvolved    int                     `json:"people_involved"`
	IncidentReport    string                  `json:"incident_report"`
}

func CreateIncidentHandler(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateIncidentRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var staff models.Staff
		err := db.NewSelect().
			Model(&staff).
			Where("agent_id = ?", req.AgentID).
			Where("department = ?", models.DeptHomelandSecurity).
			Scan(context.Background())

		if err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "Unauthorized: Only Homeland Security staff can report incidents")
			return
		}

		incident := models.Incident{
			AgentID:           req.AgentID,
			Department:        req.Department,
			IncidentType:      req.IncidentType,
			Severity:          req.Severity,
			CallerFullName:    req.CallerFullName,
			CallerPhoneNumber: req.CallerPhoneNumber,
			CallerLocation:    req.CallerLocation,
			PeopleInvolved:    req.PeopleInvolved,
			IncidentReport:    req.IncidentReport,
			StaffID:           staff.ID,
		}

		_, err = db.NewInsert().Model(&incident).Exec(context.Background())
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create incident")
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Incident created successfully"})
	}
}

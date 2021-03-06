package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) Test_ServiceAreaEffectiveDateValidation() {
	now := time.Now()

	validServiceArea := Tariff400ngServiceArea{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(1, 0, 0),
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		EffectiveDateLower: now,
		EffectiveDateUpper: now.AddDate(-1, 0, 0),
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors = map[string][]string{
		"effective_date_upper": []string{"EffectiveDateUpper must be after EffectiveDateLower."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaServiceChargeValidation() {
	validServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: 100,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validServiceArea, expErrors)

	invalidServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: -1,
		SIT185ARateCents:   unit.Cents(50),
		SIT185BRateCents:   unit.Cents(50),
		SITPDSchedule:      1,
	}

	expErrors = map[string][]string{
		"service_charge_cents": []string{"-1 is not greater than -1."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

func (suite *ModelSuite) Test_ServiceAreaSITRatesValidation() {
	invalidServiceArea := Tariff400ngServiceArea{
		ServiceChargeCents: 1,
	}

	expErrors := map[string][]string{
		"s_i_t185_b_rate_cents": []string{"SIT185BRateCents can not be blank."},
		"s_i_t_p_d_schedule":    []string{"SITPDSchedule can not be blank."},
		"s_i_t185_a_rate_cents": []string{"SIT185ARateCents can not be blank."},
	}
	suite.verifyValidationErrors(&invalidServiceArea, expErrors)
}

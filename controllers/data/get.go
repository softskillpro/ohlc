// Package data defines the data controllers for our web application.
package data

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ohcl/controllers/inputForms"
	"ohcl/controllers/outputForms"
	"ohcl/database"
	"ohcl/models/ohcl"
)

// Get returns all data.
// @Description get
// @Summary returns all data.
// @Tags data
// @Router /data [GET]
// @Accept application/json
// @Param page query int true "page" default(1)
// @Param per_page query int true "per page" default(5)
// @Param symbol query string false "symbol"
// @Param open query string false "open"
// @Param unix query string false "unix"
// @Param high query string false "high"
// @Param low query string false "low"
// @Param close query string false "close"
// @Param created_at query string false "created at"
// @Success 200 {object} outputForms.State
// @Failure 500 {object} outputForms.State
func Get(c *gin.Context) {

	// Parse inputForm.
	// Create an instance of the struct GetData and bind the input parameters to it.
	inputForm := new(inputForms.GetData)
	if err := c.BindQuery(inputForm); err != nil {
		// Return a JSON response with an error message if the input validation fails.
		c.JSON(http.StatusBadRequest, outputForms.NewState().SetCode(http.StatusBadRequest).
			SetStatus(false).SetMessage(invalidInputQuery).SetDetail(err.Error()))
		return
	}

	// handle inputForm, fetch perPage and page variables.
	// Call the Handle method of the inputForm to set default values if not present and validate and sanitize the input values.
	inputForm.Handle()

	// Declare variables for holding the queried records.
	query := database.QueryOption().SetPagination(inputForm.PerPage, inputForm.Page).SetTableName(ohcl.TableName)

	// check all inputform fields, if fields are exist, then it add into query.
	if inputForm.Symbol != "" {
		format := "%" + inputForm.Symbol + "%"
		query.SetCondition("symbol LIKE ", format)
	}

	if inputForm.Unix != "" {
		format := "%" + inputForm.Unix + "%"
		query.SetCondition("CAST (unix AS VARCHAR) LIKE ", format)
	}

	if inputForm.Open != "" {
		format := "%" + inputForm.Open + "%"
		query.SetCondition("CAST (open AS VARCHAR) LIKE ", format)
	}

	if inputForm.Close != "" {
		format := "%" + inputForm.Close + "%"
		query.SetCondition("CAST (close AS VARCHAR) LIKE ", format)
	}

	if inputForm.High != "" {
		format := "%" + inputForm.High + "%"
		query.SetCondition("CAST (high AS VARCHAR) LIKE ", format)
	}

	if inputForm.Low != "" {
		format := "%" + inputForm.Low + "%"
		query.SetCondition("CAST (low AS VARCHAR) LIKE ", format)
	}

	// Declare variables for holding the queried records.
	records := make([]ohcl.Model, 0)

	// Execute the query to retrieve records based on the input parameters.
	// Limit the records by the perPage count, and offset them by the page count.
	if err := database.Db().GetAll(query, &records).Error(); err != nil {
		// If any error occurs, return a JSON response with an error message.
		c.JSON(http.StatusInternalServerError, outputForms.NewState().
			SetStatus(false).SetCode(http.StatusInternalServerError).SetMessage(internalServerError).
			SetDetail(err.Error()))
		return
	}

	// Return a JSON response with the queried records and the count information.
	c.JSON(http.StatusOK, outputForms.NewState().SetStatus(true).
		SetCode(http.StatusOK).SetData(records).SetMessage(dataFoundedSuccessfully))
}

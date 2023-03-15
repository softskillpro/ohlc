// Package data defines the controllers for our web application.
package data

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"ohcl/controllers/outputForms"
	"ohcl/database"
	"ohcl/models/ohcl"
	"ohcl/tools"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Post function handles the POST request sent to the /data endpoint.
// It reads and validates the files provided in the formData and constructs an array of Model data records accordingly.
// The records are then saved to the database in chunks with a Go routine per chunk to speed up the insertion.
// The function returns a JSON response with the number of records saved and the time duration of the insertion process.
// @Description post data
// @Summary post gives a files from formData and insert records quickly.
// @Tags data
// @Router /data [POST]
// @Param files formData file true "resource files"
// @Accept application/json
// @Success 200 {object} outputForms.State
// @Failure 400 {object} outputForms.State
// @Failure 500 {object} outputForms.State
func Post(c *gin.Context) {

	// Start timer to calculate duration of function
	start := time.Now()

	// Create wait group to synchronously handle the records being written to the database
	var wg sync.WaitGroup
	var goRoutinesCount int

	// get fileForms from context.
	form, err := c.MultipartForm()
	if err != nil {
		// If there is an error, return an internal server error with a JSON response
		c.JSON(http.StatusInternalServerError, outputForms.NewState().
			SetStatus(false).SetCode(http.StatusInternalServerError).SetMessage(internalServerError).
			SetDetail(err.Error()))
		return
	}

	files := form.File["files"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, outputForms.NewState().SetCode(http.StatusBadRequest).SetMessage(noFilesUploaded).
			SetStatus(false))
		return
	}

	// Set initial lineCount and records slice
	lineCount := 1
	records := make([]ohcl.Model, 0)

	// Iterate through each file and read its contents.
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, outputForms.NewState().
				SetStatus(false).SetCode(http.StatusInternalServerError).SetMessage(internalServerError).
				SetDetail(err.Error()))
			return
		}

		if !strings.Contains(file.Header.Get("Content-Type"), "text/csv") {
			c.JSON(http.StatusBadRequest, outputForms.NewState().SetCode(http.StatusBadRequest).
				SetStatus(false).SetMessage(invalidCSVFile).SetDetail(err.Error()))
			return
		}

		// Create a CSV reader for the file.
		reader := csv.NewReader(f)

		// Validate the CSV file for proper formatting.
		if ok, err := tools.ValidateCSVHeader(reader); !ok || err != nil {
			c.JSON(http.StatusBadRequest, outputForms.NewState().SetCode(http.StatusBadRequest).
				SetStatus(false).SetMessage(invalidCSVFile).SetDetail(err.Error()))
			return
		}

		// Read the CSV file line by line.
		for {
			lineCount++

			// Read the next record from the CSV reader
			record, err := reader.Read()

			if err != nil {
				// If the error is an EOF, break out of the loop
				if err == io.EOF {
					break
				}

				// If there is another error, return an internal server error with a JSON response
				c.JSON(http.StatusInternalServerError, outputForms.NewState().
					SetStatus(false).SetCode(http.StatusInternalServerError).SetMessage(internalServerError).
					SetDetail(err.Error()))
				return
			}

			// validate row.
			if err := tools.ValidateOneRow(record); err != nil {
				c.JSON(http.StatusBadRequest, outputForms.NewState().SetStatus(false).SetCode(http.StatusBadRequest).
					SetMessage(err.Error()).SetDetail(err.Error()))
				return
			}

			// fetch records into model.
			// Create a new instance of the Model struct and populate it with the data from the CSV file.
			d := new(ohcl.Model)
			d.Symbol = record[1]

			// Parse the string values from the CSV file to the necessary data types.
			u, _ := strconv.Atoi(record[0])
			d.Unix = int64(u)
			d.Open, _ = strconv.ParseFloat(record[2], 64)
			d.High, _ = strconv.ParseFloat(record[3], 64)
			d.Low, _ = strconv.ParseFloat(record[4], 64)
			d.Close, _ = strconv.ParseFloat(record[5], 64)

			// Append the new Model record to the slice of records
			records = append(records, *d)

			//  If the size of the records slice is greater than or equal to 3000
			// (which is close to the maximum number of records that can be written in a single bulk query to the database),
			// use a new Goroutine to asynchronously save the records to the database
			if len(records) == 3000 {
				wg.Add(1)
				goRoutinesCount++

				go func() {
					defer wg.Done()

					recordGoRoutine := records

					if err := database.Db().InsertMany(ohcl.TableName, recordGoRoutine).Error(); err != nil {
						return
					}

				}()

				// Zero out the records array to make it ready for the next chunk
				records = []ohcl.Model{}
			}

			if goRoutinesCount == 10 {
				wg.Wait()
				goRoutinesCount = 0
			}

		}

		if len(records) != 0 {

			recordGoRoutine := records

			if err := database.Db().InsertMany(ohcl.TableName, recordGoRoutine).Error(); err != nil {
				c.JSON(http.StatusInternalServerError, outputForms.NewState().
					SetStatus(false).SetCode(http.StatusInternalServerError).SetMessage(internalServerError).
					SetDetail(err.Error()))
				return
			}
		}
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	// End timer and set the duration time
	end := time.Now()
	duration := end.Sub(start)

	// Return the total number of records saved, the total number of files read, and the time duration of the insertion process in a JSON response
	c.JSON(http.StatusOK, outputForms.NewState().SetCode(http.StatusOK).SetStatus(false).
		SetMessage(fmt.Sprintf("%v records stored and %v files read successfully, seconds spend %v",
			lineCount, len(files), duration.Seconds())))
}

package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/stepkareserva/obsermon/internal/models"
	"github.com/stepkareserva/obsermon/internal/server/http/constants"
	"github.com/stepkareserva/obsermon/internal/server/http/errors"
	"go.uber.org/zap"
)

type ValuesHandler struct {
	service Service
	errors.ErrorsWriter
}

func NewValuesHandler(s Service, log *zap.Logger) (*ValuesHandler, error) {
	if s == nil {
		return nil, fmt.Errorf("service not exists")
	}
	return &ValuesHandler{
		service:      s,
		ErrorsWriter: errors.NewErrorsWriter(log),
	}, nil
}

func (h *ValuesHandler) MetricValuesHandler() http.HandlerFunc {
	var tmpl = template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>
			Obsermon metrics
		</title>
	    <style>
			body {
				font-family: monospace;
			}
			table {
				border-collapse: collapse;
				width: auto;
			}
			td {
				padding: 0;
    			padding-right: 1ch;
				border: none;
				line-height: 1.5; 
			}
		</style>
	</head>
	<body>
		<h1>Gauges:</h1>
		<table>
		{{range .Gauges}}
		<tr>
			<td>{{.Name}}</td>
			<td>{{.Value.PrettyString}}</td>
		</tr>
		{{end}}
		</table>

		<h1>Counters:</h1>
		<table>
		{{range .Counters}}
		<tr>
			<td>{{.Name}}</td>
			<td>{{.Value}}</td>
		</tr>
		{{end}}
		</table>
	</body>
	</html>`))

	return func(w http.ResponseWriter, r *http.Request) {
		gauges, err := h.service.ListGauges(r.Context())
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		counters, err := h.service.ListCounters(r.Context())
		if err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}

		templateData := struct {
			Gauges   []models.Gauge
			Counters []models.Counter
		}{
			Gauges:   gauges,
			Counters: counters,
		}

		w.Header().Set(constants.ContentType, constants.ContentTypeHTML)
		if err := tmpl.Execute(w, templateData); err != nil {
			h.WriteError(w, errors.ErrInternalServerError, err.Error())
			return
		}
	}
}

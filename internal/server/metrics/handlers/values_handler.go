package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/internal/models"
	"go.uber.org/zap"

	hc "github.com/stepkareserva/obsermon/internal/server/httpconst"
)

func ValuesHandler(s Service, log *zap.Logger) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics service is nil")
	}

	r := chi.NewRouter()
	r.Get("/", metricValuesHandler(s, log))

	return r, nil
}

func metricValuesHandler(s Service, log *zap.Logger) http.HandlerFunc {
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
		gauges, err := s.ListGauges()
		if err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		counters, err := s.ListCounters()
		if err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}

		templateData := struct {
			Gauges   []models.Gauge
			Counters []models.Counter
		}{
			Gauges:   gauges,
			Counters: counters,
		}

		w.Header().Set(hc.ContentType, hc.ContentTypeHTML)
		if err := tmpl.Execute(w, templateData); err != nil {
			WriteError(w, ErrInternalServerError, log)
			return
		}
	}
}

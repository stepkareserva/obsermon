package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/stepkareserva/obsermon/cmd/server/internal/metrics/server"
	"github.com/stepkareserva/obsermon/internal/models"
)

func ValuesHandler(s *server.Server) (http.Handler, error) {
	if s == nil {
		return nil, fmt.Errorf("metrics server is nil")
	}

	r := chi.NewRouter()
	r.Get("/", metricValuesHandler(s))

	return r, nil
}

func metricValuesHandler(s *server.Server) http.HandlerFunc {
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
			WriteError(w, ErrInternalServerError)
			return
		}

		counters, err := s.ListCounters()
		if err != nil {
			WriteError(w, ErrInternalServerError)
			return
		}

		templateData := struct {
			Gauges   []models.Gauge
			Counters []models.Counter
		}{
			Gauges:   gauges,
			Counters: counters,
		}

		// Q: maybe some package with such constants exists?
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, templateData); err != nil {
			return
		}
	}
}

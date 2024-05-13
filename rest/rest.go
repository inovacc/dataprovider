package rest

//
//import (
//	"encoding/json"
//	"github.com/dyammarcano/dataprovider/querier"
//	"net/http"
//)
//
//type Rest interface {
//	GetResource(w http.ResponseWriter, r *http.Request)
//	SetResource(w http.ResponseWriter, r *http.Request)
//	UpdateResource(w http.ResponseWriter, r *http.Request)
//	DeleteResource(w http.ResponseWriter, r *http.Request)
//}
//
//type ResourceHandler struct {
//	querier querier.Querier
//}
//
//func NewResourceHandler(q querier.Querier) *ResourceHandler {
//	return &ResourceHandler{
//		querier: q,
//	}
//}
//
//func (h *ResourceHandler) GetResource(w http.ResponseWriter, r *http.Request) {
//	// Use the querier to build and execute a SELECT query
//	query := h.querier.Select("column1", "column2").From("table").Build()
//
//	// Execute the query and get the result
//	// This is a placeholder, replace with your actual database query execution
//	result, err := executeQuery(query)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	// Convert the result to JSON and write it to the response
//	json.NewEncoder(w).Encode(result)
//}
//
//func main() {
//	q := querier.NewQuerier()
//	handler := rest.NewResourceHandler(q)
//
//	http.HandleFunc("/resource", func(w http.ResponseWriter, r *http.Request) {
//		switch r.Method {
//		case http.MethodGet:
//			handler.GetResource(w, r)
//		// TODO: Add cases for other HTTP methods (POST, PUT, DELETE)
//		default:
//			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
//		}
//	})
//
//	http.ListenAndServe(":8080", nil)
//}

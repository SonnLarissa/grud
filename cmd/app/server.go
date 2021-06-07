package app

import (
	"encoding/json"
	"errors"
	"github.com/SonnLarissa/grud/pkg/customers"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Server представляет собой логический сервер нашего приложения
type Server struct {
	mux         *http.ServeMux
	customerSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customerSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

//INIt инициализирует сервер (регистрирует все Handlerы)
func (s *Server) Init() {
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
	s.mux.HandleFunc("/customers.blockById", s.handleBlockByID)
	s.mux.HandleFunc("/customers.unblockById", s.handleUnBlockByID)
	s.mux.HandleFunc("/customers.removeById", s.handleDelete)
	s.mux.HandleFunc("/customers.save", s.handleSave)
}


func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	item, err := s.customerSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}


func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customerSvc.All(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, items)
}


func (s *Server) handleGetAllActiveCustomers(writer http.ResponseWriter, request *http.Request) {

	items, err := s.customerSvc.AllActive(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	parceJSON(writer, items)
}


func (s *Server) handleBlockByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.ChangeActive(request.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		parceJSON(writer, item)
	}
}



func (s *Server) handleUnBlockByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.ChangeActive(request.Context(), id, true)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, item)
}



func (s *Server) handleDelete(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.Delete(request.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, item)
}


func (s *Server) handleSave(writer http.ResponseWriter, request *http.Request) {
	idParam := request.FormValue("id")
	name := request.FormValue("name")
	phone := request.FormValue("phone")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if name == "" && phone == "" {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item := &customers.Customer{
		ID:      id,
		Name:    name,
		Phone:   phone,
		Active:  true,
		Created: time.Now(),
	}
	customer, err := s.customerSvc.Save(request.Context(), item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	parceJSON(writer, customer)
}




func parceJSON(writer http.ResponseWriter, iData interface{}) {

	data, err := json.Marshal(iData)

	if err != nil {
		log.Println(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

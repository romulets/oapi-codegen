// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns all pets
	// (GET /pets)
	FindPets(w http.ResponseWriter, r *http.Request, params FindPetsParams)
	// Creates a new pet
	// (POST /pets)
	AddPet(w http.ResponseWriter, r *http.Request)
	// Deletes a pet by ID
	// (DELETE /pets/{id})
	DeletePet(w http.ResponseWriter, r *http.Request, id int64)
	// Returns a pet by ID
	// (GET /pets/{id})
	FindPetByID(w http.ResponseWriter, r *http.Request, id int64)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// FindPets operation middleware
func (siw *ServerInterfaceWrapper) FindPets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FindPetsParams

	// ------------- Optional query parameter "tags" -------------

	err = runtime.BindQueryParameter("form", true, false, "tags", r.URL.Query(), &params.Tags)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "tags", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	headers := r.Header

	// ------------- Required header parameter "Accept-Language" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Accept-Language")]; found {
		var AcceptLanguage []string

		for _, value := range valueList {
			var temp []string
			err = runtime.BindStyledParameterWithLocation("simple", false, "Accept-Language", runtime.ParamLocationHeader, value, &temp)
			if err != nil {
				siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "Accept-Language", Err: err})
				return
			}
			AcceptLanguage = append(AcceptLanguage, temp...)
		}

		params.AcceptLanguage = AcceptLanguage

	} else {
		err := fmt.Errorf("Header parameter Accept-Language is required, but not found")
		siw.ErrorHandlerFunc(w, r, &RequiredHeaderError{ParamName: "Accept-Language", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindPets(w, r, params)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddPet operation middleware
func (siw *ServerInterfaceWrapper) AddPet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddPet(w, r)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeletePet operation middleware
func (siw *ServerInterfaceWrapper) DeletePet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeletePet(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// FindPetByID operation middleware
func (siw *ServerInterfaceWrapper) FindPetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, chi.URLParam(r, "id"), &id)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindPetByID(w, r, id)
	})

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshallingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshallingParamError) Error() string {
	return fmt.Sprintf("Error unmarshalling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshallingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pets", wrapper.FindPets)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/pets", wrapper.AddPet)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/pets/{id}", wrapper.DeletePet)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pets/{id}", wrapper.FindPetByID)
	})

	return r
}

type FindPetsRequestObject struct {
	Params FindPetsParams
}

type FindPetsResponseObject interface {
	VisitFindPetsResponse(w http.ResponseWriter) error
}

type FindPets200JSONResponse []Pet

func (response FindPets200JSONResponse) VisitFindPetsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type FindPetsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response FindPetsdefaultJSONResponse) VisitFindPetsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type AddPetRequestObject struct {
	Body *AddPetJSONRequestBody
}

type AddPetResponseObject interface {
	VisitAddPetResponse(w http.ResponseWriter) error
}

type AddPet200JSONResponse Pet

func (response AddPet200JSONResponse) VisitAddPetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AddPetdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AddPetdefaultJSONResponse) VisitAddPetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type DeletePetRequestObject struct {
	Id int64 `json:"id"`
}

type DeletePetResponseObject interface {
	VisitDeletePetResponse(w http.ResponseWriter) error
}

type DeletePet204Response struct {
}

func (response DeletePet204Response) VisitDeletePetResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type DeletePetdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response DeletePetdefaultJSONResponse) VisitDeletePetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type FindPetByIDRequestObject struct {
	Id int64 `json:"id"`
}

type FindPetByIDResponseObject interface {
	VisitFindPetByIDResponse(w http.ResponseWriter) error
}

type FindPetByID200JSONResponse Pet

func (response FindPetByID200JSONResponse) VisitFindPetByIDResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type FindPetByIDdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response FindPetByIDdefaultJSONResponse) VisitFindPetByIDResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Returns all pets
	// (GET /pets)
	FindPets(ctx context.Context, request FindPetsRequestObject) (FindPetsResponseObject, error)
	// Creates a new pet
	// (POST /pets)
	AddPet(ctx context.Context, request AddPetRequestObject) (AddPetResponseObject, error)
	// Deletes a pet by ID
	// (DELETE /pets/{id})
	DeletePet(ctx context.Context, request DeletePetRequestObject) (DeletePetResponseObject, error)
	// Returns a pet by ID
	// (GET /pets/{id})
	FindPetByID(ctx context.Context, request FindPetByIDRequestObject) (FindPetByIDResponseObject, error)
}

type StrictHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// FindPets operation middleware
func (sh *strictHandler) FindPets(w http.ResponseWriter, r *http.Request, params FindPetsParams) {
	var request FindPetsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.FindPets(ctx, request.(FindPetsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "FindPets")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(FindPetsResponseObject); ok {
		if err := validResponse.VisitFindPetsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// AddPet operation middleware
func (sh *strictHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	var request AddPetRequestObject

	var body AddPetJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.AddPet(ctx, request.(AddPetRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddPet")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(AddPetResponseObject); ok {
		if err := validResponse.VisitAddPetResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// DeletePet operation middleware
func (sh *strictHandler) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {
	var request DeletePetRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.DeletePet(ctx, request.(DeletePetRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeletePet")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(DeletePetResponseObject); ok {
		if err := validResponse.VisitDeletePetResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// FindPetByID operation middleware
func (sh *strictHandler) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) {
	var request FindPetByIDRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.FindPetByID(ctx, request.(FindPetByIDRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "FindPetByID")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(FindPetByIDResponseObject); ok {
		if err := validResponse.VisitFindPetByIDResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("Unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RYS28jyQ3+K0Qlx96WM7vIQad4x7OAgN0ZJ7Oby44PdDUlcVEvV7HkEQz994DVrZfl",
	"8eSFIEEultRdLH7k95HF8pOx0acYKEgx8ydT7Jo8tq/vco5Zv6QcE2Vhao9tHEg/Byo2cxKOwczHxdDe",
	"dWYZs0cxc8NBvn1jOiPbRONPWlE2u854KgVXX9xo//pgWiRzWJndrjOZHipnGsz8VzM53C+/23XmPT3e",
	"klziDuhfcPcePUFcgqwJEsmlw84Iri7tft6m1+2eAW3eFd6EDZ37sDTzX5/M7zMtzdz8bnYkYjaxMJti",
	"2XXPg+HhEtIvgR8qAQ/nuE7J+ON3L5DxDCkP5m53t9PHHJZxpDwI2oabPLIzc4OJhdD/qTziakW552i6",
	"KcXm4/gMrm8X8DOhN52pWY3WIqnMZ7MTo133LIprKOiTo2YtaxSohQqgRlMkZgIsgAHo87hMIgzkYyiS",
	"UQiWhFIzFeDQcvAhUdCdvu2voCSyvGSLzVVnHFsKhY7iMNcJ7ZrgTX91gfnx8bHH9rqPeTWbbMvsx8Xb",
	"d+8/vvvmTX/Vr8W7phjKvnxYfqS8YUsvBj5ra2ZKB4s7zdrtFKfpzIZyGbPyh/6qv9KtY6KAic3cfNse",
	"dSahrJsmZpoh/bIaJXae17+Q1BwKoHMtlbDM0bcUlW0R8mOu9XctlGGtWbaWSgGJn8J79FBoABvDwJ6C",
	"VA9UpIefkCwFLCDkU8xQcMUiXKBgYgodBLKQ1zHYWqCQP1nAAuhJerimQBgABVYZNzwgYF1V6gAtMNrq",
	"uJn28LZmvGepGeLAEVzM5DuIOWAmoBUJkKMJXSDbga251KIl4chKLT3cVC7gGaTmxKWDVN2GA2b1RTlq",
	"0B0IB8tDDQIbzFwL/FaLxB4WAdZoYa0gsBSC5FAIYWAr1Ws6FmNRaSw4cOJiOawAg2g0x9gdr6rDQ+Rp",
	"jZkk4z6Juh58dFSECdgnygNrpv7KG/RjQOj4oaKHgVEzk7HAg8a2IccCIQaQmCVmTQkvKQwH7z3cZqRC",
	"QRQmBfZHADUHhE10VRIKbChQQAU8Jlf/eKxZ91iE485LylPWl2jZcTlz0jzon+7Ir4USB3SkxA6d5tFS",
	"RtHA9LOHj7UkCgNrlh2qeIboYu5UgYWsqJpblE0qGnUHG1qzrQ5BW1seqgfH95RjDz/FfM9AlYuPwykN",
	"+roJ26HlwNh/Cp/CRxoaE7XAklR8Lt7H3AwoHhWTq+Tqe9Da8Ng2nJLPxXVA9axaRsrBVdWhqrOH2zUW",
	"cm4sjER5Mm9pbvSSwBKr5fs6Jhz3fnTdqf2G3EQdbyhn7M5da50AD92hEAPfr3v4RSCRcxSEip4cKZZK",
	"Wkn7IupBU4H7KtCi2+dyv9M+rJbJrgE5yCLUYEEyF2kH04YFqYcfarEEJK0bDJUPVaCdolhylLnBGfW7",
	"N/CqlopNPLb6ggE8rjRkchNbPfy5jqY+OuVtZI/qqJ0jlO7QfACr1SIZV07yHMOexDE1mUM1qliUYODQ",
	"HaFMhRu48B5wUQyWpQ6sUEtBqLLX2UTk6Oksac1fD7enxLTMTRhTJuHqTzrXKJranehbW2//Sc84HRra",
	"ebcYzNz8wGHQ86UdG1kTQLm0KeT8sBBcad+HJTuhDPdbo8OAmZuHSnl7POl1nemmobHNJUK+nUGXU9T4",
	"AHPGrf4usm3Hno4nbcA5R+DxM3tt49XfU9aJJlOpThqs3M6yL2By7FnOQH11HFX3bas14UD5uNe1tZTk",
	"mx8xrOo4jh4nJcmV/pnQd3e6S0nax9r6N1dX+yGLwjgcpuSmMWX2W9F8PL3k6LXJcRwbn7u+mLYSCezB",
	"jLPYEquTfwjPazDGO8QLjmugz0n7uDb8cU1nSvUe8/aFaUWxpVhemGveZkJpA2KgR127n/zaEKUH/ohd",
	"l+jw6Fx8pOGiMq4HLYyJYCryfRy2/7Ys7Mf4yzTckqigcRj04wD7Qmi7f1EzX5XK/440Lghv79vwO3vi",
	"YTdKxJG8cNsbn6tt4bBy7YoE96g9PY6qWdxAqRrTCxq5adajTF5tn4sbbVhp5HbCMjUrndaP/YWHV1vK",
	"169ul73ku8uoFciIYvhvIvLmQEZjYQuLG4X3+u3lnLEDj4ubL51132/bu7+fryWJXf/H6Pq/LeNnjI7s",
	"tyWUN3uazq/g+/8A9CfXaL0L7+52fwsAAP//ysTaOscSAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

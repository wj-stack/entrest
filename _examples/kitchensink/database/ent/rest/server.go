// Code generated by ent, DO NOT EDIT.

package rest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/category"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/friendship"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/pet"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/privacy"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/settings"
	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent/user"
)

//go:embed openapi.json
var OpenAPI []byte // OpenAPI contains the JSON schema of the API.

// Operation represents the CRUD operation(s).
type Operation string

const (
	// OperationCreate represents the create operation (method: POST).
	OperationCreate Operation = "create"
	// OperationRead represents the read operation (method: GET).
	OperationRead Operation = "read"
	// OperationUpdate represents the update operation (method: PATCH).
	OperationUpdate Operation = "update"
	// OperationDelete represents the delete operation (method: DELETE).
	OperationDelete Operation = "delete"
	// OperationList represents the list operation (method: GET).
	OperationList Operation = "list"
)

// ErrorResponse is the response structure for errors.
type ErrorResponse struct {
	Error     string `json:"error"`                // The underlying error, which may be masked when debugging is disabled.
	Type      string `json:"type"`                 // A summary of the error code based off the HTTP status code or application error code.
	Code      int    `json:"code"`                 // The HTTP status code or other internal application error code.
	RequestID string `json:"request_id,omitempty"` // The unique request ID for this error.
	Timestamp string `json:"timestamp,omitempty"`  // The timestamp of the error, in RFC3339 format.
}

type ErrBadRequest struct {
	Err error
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", e.Err)
}

func (e ErrBadRequest) Unwrap() error {
	return e.Err
}

// IsBadRequest returns true if the unwrapped/underlying error is of type ErrBadRequest.
func IsBadRequest(err error) bool {
	var target *ErrBadRequest
	return errors.As(err, &target)
}

var ErrEndpointNotFound = errors.New("endpoint not found")

// IsEndpointNotFound returns true if the unwrapped/underlying error is of type ErrEndpointNotFound.
func IsEndpointNotFound(err error) bool {
	return errors.Is(err, ErrEndpointNotFound)
}

var ErrMethodNotAllowed = errors.New("method not allowed")

// IsMethodNotAllowed returns true if the unwrapped/underlying error is of type ErrMethodNotAllowed.
func IsMethodNotAllowed(err error) bool {
	return errors.Is(err, ErrMethodNotAllowed)
}

// JSON marshals 'v' to JSON, and setting the Content-Type as application/json.
// Note that this does NOT auto-escape HTML. If 'v' cannot be marshalled to JSON,
// this will panic.
//
// JSON also supports prettification when the origin request has a query parameter
// of "pretty" set to true.
func JSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)

	if pretty, _ := strconv.ParseBool(r.FormValue("pretty")); pretty {
		enc.SetIndent("", "    ")
	}

	if err := enc.Encode(v); err != nil && err != io.EOF {
		panic(fmt.Sprintf("failed to marshal response: %v", err))
	}
}

var (
	// DefaultDecoder is the default decoder used by Bind. You can either override
	// this, or provide your own. Make sure it is set before Bind is called.
	DefaultDecoder = form.NewDecoder()

	// DefaultDecodeMaxMemory is the maximum amount of memory in bytes that will be
	// used for decoding multipart/form-data requests.
	DefaultDecodeMaxMemory int64 = 8 << 20
)

// Bind decodes the request body to the given struct. At this time the only supported
// content-types are application/json, application/x-www-form-urlencoded, as well as
// GET parameters.
func Bind(r *http.Request, v any) error {
	err := r.ParseForm()
	if err != nil {
		return &ErrBadRequest{Err: fmt.Errorf("parsing form parameters: %w", err)}
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		err = DefaultDecoder.Decode(v, r.Form)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		switch {
		case strings.HasPrefix(r.Header.Get("Content-Type"), "application/json"):
			dec := json.NewDecoder(r.Body)
			defer r.Body.Close()
			err = dec.Decode(v)
		case strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data"):
			err = r.ParseMultipartForm(DefaultDecodeMaxMemory)
			if err == nil {
				err = DefaultDecoder.Decode(v, r.MultipartForm.Value)
			}
		default:
			err = DefaultDecoder.Decode(v, r.PostForm)
		}
	default:
		return &ErrBadRequest{Err: fmt.Errorf("unsupported method %s", r.Method)}
	}

	if err != nil {
		return &ErrBadRequest{Err: fmt.Errorf("error decoding %s request into required format (%T): %w", r.Method, v, err)}
	}
	return nil
}

// Req simplifies making an HTTP handler that returns a single result, and an error.
// The result, if not nil, must be JSON-marshalable. If result is nil, [http.StatusNoContent]
// will be returned.
func Req[Resp any](s *Server, op Operation, fn func(*http.Request) (*Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := fn(r)
		handleResponse(s, w, r, op, results, err)
	}
}

// ReqID is similar to Req, but also processes an "id" path parameter and provides it to the
// handler function.
func ReqID[Resp any](s *Server, op Operation, fn func(*http.Request, int) (*Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			handleResponse[Resp](s, w, r, op, nil, err)
			return
		}
		results, err := fn(r, id)
		handleResponse(s, w, r, op, results, err)
	}
}

// ReqParam is similar to Req, but also processes a request body/query params and provides it
// to the handler function.
func ReqParam[Params, Resp any](s *Server, op Operation, fn func(*http.Request, *Params) (*Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := new(Params)
		if err := Bind(r, params); err != nil {
			handleResponse[Resp](s, w, r, op, nil, err)
			return
		}
		results, err := fn(r, params)
		handleResponse(s, w, r, op, results, err)
	}
}

// ReqIDParam is similar to ReqParam, but also processes an "id" path parameter and request
// body/query params, and provides it to the handler function.
func ReqIDParam[Params, Resp any](s *Server, op Operation, fn func(*http.Request, int, *Params) (*Resp, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			handleResponse[Resp](s, w, r, op, nil, err)
			return
		}
		params := new(Params)
		err = Bind(r, params)
		if err != nil {
			handleResponse[Resp](s, w, r, op, nil, err)
			return
		}
		results, err := fn(r, id, params)
		handleResponse(s, w, r, op, results, err)
	}
}

// Links represents a set of linkable-relationsips that can be represented through
// the "Link" header. Note that all urls must be url-encoded already.
type Links map[string]string

func (l Links) String() string {
	var links []string
	var keys []string
	for k := range l {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		links = append(links, fmt.Sprintf(`<%s>; rel=%q`, l[k], k))
	}
	return strings.Join(links, ", ")
}

type linkablePagedResource interface {
	GetPage() int
	GetIsLastPage() bool
}

// Spec returns the OpenAPI spec for the server implementation.
func (s *Server) Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(OpenAPI)
}

var scalarTemplate = template.Must(template.New("docs").Parse(`<!DOCTYPE html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <link rel="icon" type="image/svg+xml"
      href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='32' height='32' viewBox='0 0 1024 1024'%3E%3Cpath fill='currentColor' d='m917.7 148.8l-42.4-42.4c-1.6-1.6-3.6-2.3-5.7-2.3s-4.1.8-5.7 2.3l-76.1 76.1a199.27 199.27 0 0 0-112.1-34.3c-51.2 0-102.4 19.5-141.5 58.6L432.3 308.7a8.03 8.03 0 0 0 0 11.3L704 591.7c1.6 1.6 3.6 2.3 5.7 2.3c2 0 4.1-.8 5.7-2.3l101.9-101.9c68.9-69 77-175.7 24.3-253.5l76.1-76.1c3.1-3.2 3.1-8.3 0-11.4M769.1 441.7l-59.4 59.4l-186.8-186.8l59.4-59.4c24.9-24.9 58.1-38.7 93.4-38.7s68.4 13.7 93.4 38.7c24.9 24.9 38.7 58.1 38.7 93.4s-13.8 68.4-38.7 93.4m-190.2 105a8.03 8.03 0 0 0-11.3 0L501 613.3L410.7 523l66.7-66.7c3.1-3.1 3.1-8.2 0-11.3L441 408.6a8.03 8.03 0 0 0-11.3 0L363 475.3l-43-43a7.85 7.85 0 0 0-5.7-2.3c-2 0-4.1.8-5.7 2.3L206.8 534.2c-68.9 69-77 175.7-24.3 253.5l-76.1 76.1a8.03 8.03 0 0 0 0 11.3l42.4 42.4c1.6 1.6 3.6 2.3 5.7 2.3s4.1-.8 5.7-2.3l76.1-76.1c33.7 22.9 72.9 34.3 112.1 34.3c51.2 0 102.4-19.5 141.5-58.6l101.9-101.9c3.1-3.1 3.1-8.2 0-11.3l-43-43l66.7-66.7c3.1-3.1 3.1-8.2 0-11.3zM441.7 769.1a131.32 131.32 0 0 1-93.4 38.7c-35.3 0-68.4-13.7-93.4-38.7a131.32 131.32 0 0 1-38.7-93.4c0-35.3 13.7-68.4 38.7-93.4l59.4-59.4l186.8 186.8z'/%3E%3C/svg%3E" />
  </head>
  <body>
    <script id="api-reference"></script>
    <script>
      document.getElementById("api-reference").dataset.configuration = JSON.stringify({
        spec: {
          url: "{{ . }}",
        },
        theme: "kepler",
        isEditable: false,
        hideDownloadButton: true,
        customCss: ".darklight-reference-promo { visibility: hidden !important; height: 0 !important; }",
      });
    </script>
    <script
      src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@1.24.26"
      integrity="sha256-Zo2w7XQtgECsnom2xI33f5AFG1VFcuJm3gogYqlgeRA="
      crossorigin="anonymous"
    ></script>
  </body>
</html>`))

func (s *Server) Docs(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := scalarTemplate.Execute(&buf, s.config.BasePath+"/openapi.json"); err != nil {
		handleResponse[struct{}](s, w, r, "", nil, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

type ServerConfig struct {
	// BaseURL is similar to [ServerConfig.BasePath], however, only the path of the URL is used
	// to prefill BasePath. This is not required if BasePath is provided.
	BaseURL string

	// BasePath if provided, and the /openapi.json endpoint is enabled, will allow annotating
	// API responses with "Link" headers. See [ServerConfig.EnableLinks] for more information.
	BasePath string

	// DisableSpecHandler if set to true, will disable the /openapi.json endpoint. This will also
	// disable the embedded API reference documentation, see [ServerConfig.DisableDocs] for more
	// information.
	DisableSpecHandler bool

	// DisableDocsHandler if set to true, will disable the embedded API reference documentation
	// endpoint at /docs. Use this if you want to provide your own documentation functionality.
	// This is disabled by default if [ServerConfig.DisableSpecHandler] is true.
	DisableDocsHandler bool

	// EnableLinks if set to true, will enable the "Link" response header, which can be used to hint
	// to clients about the location of the OpenAPI spec, API documentation, how to auto-paginate
	// through results, and more.
	EnableLinks bool

	// MaskErrors if set to true, will mask the error message returned to the client,
	// returning a generic error message based on the HTTP status code.
	MaskErrors bool

	// ErrorHandler is invoked when an error occurs. If not provided, the default
	// error handling logic will be used. If you want to run logic on errors, but
	// not actually handle the error yourself, you can still call [Server.DefaultErrorHandler]
	// after your logic.
	ErrorHandler func(w http.ResponseWriter, r *http.Request, op Operation, err error)

	// GetReqID returns the request ID for the given request. If not provided, the
	// default implementation will use the X-Request-Id header, otherwise an empty
	// string will be returned. If using go-chi, middleware.GetReqID will be used.
	GetReqID func(r *http.Request) string
}

type Server struct {
	db     *ent.Client
	config *ServerConfig
}

// NewServer returns a new auto-generated server implementation for your ent schema.
// [Server.Handler] returns a ready-to-use http.Handler that mounts all of the
// necessary endpoints.
func NewServer(db *ent.Client, config *ServerConfig) (*Server, error) {
	s := &Server{
		db:     db,
		config: config,
	}
	if s.config == nil {
		s.config = &ServerConfig{}
	}
	if s.config.BaseURL != "" && s.config.BasePath == "" {
		uri, err := url.Parse(s.config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse BaseURL: %w", err)
		}
		s.config.BasePath = uri.Path
	}
	if s.config.BasePath != "" {
		if !strings.HasPrefix(s.config.BasePath, "/") {
			s.config.BasePath = "/" + s.config.BasePath
		}
		s.config.BasePath = strings.TrimRight(s.config.BasePath, "/")
	}
	return s, nil
}

// DefaultErrorHandler is the default error handler for the Server.
func (s *Server) DefaultErrorHandler(w http.ResponseWriter, r *http.Request, op Operation, err error) {
	ts := time.Now().UTC().Format(time.RFC3339)

	resp := ErrorResponse{
		Error:     err.Error(),
		Timestamp: ts,
	}

	var numErr *strconv.NumError

	switch {
	case IsEndpointNotFound(err):
		resp.Code = http.StatusNotFound
	case IsMethodNotAllowed(err):
		resp.Code = http.StatusMethodNotAllowed
	case IsBadRequest(err):
		resp.Code = http.StatusBadRequest
	case errors.Is(err, privacy.Deny):
		resp.Code = http.StatusForbidden
	case ent.IsNotFound(err):
		if op == OperationList {
			resp.Type = "No results found matching the given query"
		}
		resp.Code = http.StatusNotFound
	case ent.IsConstraintError(err), ent.IsNotSingular(err):
		resp.Code = http.StatusConflict
	case ent.IsValidationError(err):
		resp.Code = http.StatusBadRequest
	case errors.As(err, &numErr):
		resp.Code = http.StatusBadRequest
		resp.Error = fmt.Sprintf("invalid ID provided: %v", err)
	default:
		resp.Code = http.StatusInternalServerError
	}

	if resp.Type == "" {
		resp.Type = http.StatusText(resp.Code)
	}
	if s.config.MaskErrors {
		resp.Error = http.StatusText(resp.Code)
	}
	if s.config.GetReqID != nil {
		resp.RequestID = s.config.GetReqID(r)
	} else {
		resp.RequestID = r.Header.Get("X-Request-Id")
	}
	JSON(w, r, resp.Code, resp)
}

func handleResponse[Resp any](s *Server, w http.ResponseWriter, r *http.Request, op Operation, resp *Resp, err error) {
	if s.config.EnableLinks {
		links := Links{}
		if !s.config.DisableSpecHandler {
			links["service-desc"] = s.config.BasePath + "/openapi.json"
			links["describedby"] = s.config.BasePath + "/openapi.json"
		}

		if err == nil && resp != nil && op == OperationList {
			if lr, ok := any(resp).(linkablePagedResource); ok {
				query := r.URL.Query()
				if page := lr.GetPage(); page > 1 {
					query.Set("page", strconv.Itoa(page-1))
					r.URL.RawQuery = query.Encode()
					links["prev"] = r.URL.String()
					if !strings.HasPrefix(links["prev"], s.config.BasePath) {
						links["prev"] = s.config.BasePath + links["prev"]
					}
				}
				if !lr.GetIsLastPage() {
					query.Set("page", strconv.Itoa(lr.GetPage()+1))
					r.URL.RawQuery = query.Encode()
					links["next"] = r.URL.String()
					if !strings.HasPrefix(links["next"], s.config.BasePath) {
						links["next"] = s.config.BasePath + links["next"]
					}
				}
			}
		}

		if v := links.String(); v != "" {
			w.Header().Set("Link", v)
		}
	}
	if err != nil {
		if s.config.ErrorHandler != nil {
			s.config.ErrorHandler(w, r, op, err)
			return
		}
		s.DefaultErrorHandler(w, r, op, err)
		return
	}
	if resp != nil {
		if r.Method == http.MethodPost {
			JSON(w, r, http.StatusCreated, resp)
			return
		}
		JSON(w, r, http.StatusOK, resp)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handler returns a ready-to-use http.Handler that mounts all of the necessary endpoints.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /categories", ReqParam(s, OperationList, s.ListCategory))
	mux.HandleFunc("GET /categories/{id}", ReqID(s, OperationRead, s.ReadCategoryByID))
	mux.HandleFunc("GET /categories/{id}/pets", ReqIDParam(s, OperationList, s.EdgeListCategoryPetsByID))
	mux.HandleFunc("POST /categories", ReqParam(s, OperationCreate, s.CreateCategory))
	mux.HandleFunc("PATCH /categories/{id}", ReqIDParam(s, OperationUpdate, s.UpdateCategoryByID))
	mux.HandleFunc("DELETE /categories/{id}", ReqID(s, OperationDelete, s.DeleteCategoryByID))
	mux.HandleFunc("GET /follows", ReqParam(s, OperationList, s.ListFollow))
	mux.HandleFunc("POST /follows", ReqParam(s, OperationCreate, s.CreateFollow))
	mux.HandleFunc("GET /friendships", ReqParam(s, OperationList, s.ListFriendship))
	mux.HandleFunc("GET /friendships/{id}", ReqID(s, OperationRead, s.ReadFriendshipByID))
	mux.HandleFunc("GET /friendships/{id}/user", ReqID(s, OperationRead, s.EdgeReadFriendshipUserByID))
	mux.HandleFunc("GET /friendships/{id}/friend", ReqID(s, OperationRead, s.EdgeReadFriendshipFriendByID))
	mux.HandleFunc("POST /friendships", ReqParam(s, OperationCreate, s.CreateFriendship))
	mux.HandleFunc("PATCH /friendships/{id}", ReqIDParam(s, OperationUpdate, s.UpdateFriendshipByID))
	mux.HandleFunc("DELETE /friendships/{id}", ReqID(s, OperationDelete, s.DeleteFriendshipByID))
	mux.HandleFunc("GET /pets", ReqParam(s, OperationList, s.ListPet))
	mux.HandleFunc("GET /pets/{id}", ReqID(s, OperationRead, s.ReadPetByID))
	mux.HandleFunc("GET /pets/{id}/categories", ReqIDParam(s, OperationList, s.EdgeListPetCategoriesByID))
	mux.HandleFunc("GET /pets/{id}/owner", ReqID(s, OperationRead, s.EdgeReadPetOwnerByID))
	mux.HandleFunc("GET /pets/{id}/friends", ReqIDParam(s, OperationList, s.EdgeListPetFriendsByID))
	mux.HandleFunc("GET /pets/{id}/followed-by", ReqIDParam(s, OperationList, s.EdgeListPetFollowedByByID))
	mux.HandleFunc("POST /pets", ReqParam(s, OperationCreate, s.CreatePet))
	mux.HandleFunc("PATCH /pets/{id}", ReqIDParam(s, OperationUpdate, s.UpdatePetByID))
	mux.HandleFunc("DELETE /pets/{id}", ReqID(s, OperationDelete, s.DeletePetByID))
	mux.HandleFunc("GET /settings", ReqParam(s, OperationList, s.ListSetting))
	mux.HandleFunc("GET /settings/{id}", ReqID(s, OperationRead, s.ReadSettingByID))
	mux.HandleFunc("GET /settings/{id}/admins", ReqIDParam(s, OperationList, s.EdgeListSettingAdminsByID))
	mux.HandleFunc("PATCH /settings/{id}", ReqIDParam(s, OperationUpdate, s.UpdateSettingByID))
	mux.HandleFunc("GET /users", ReqParam(s, OperationList, s.ListUser))
	mux.HandleFunc("GET /users/{id}", ReqID(s, OperationRead, s.ReadUserByID))
	mux.HandleFunc("GET /users/{id}/pets", ReqIDParam(s, OperationList, s.EdgeListUserPetsByID))
	mux.HandleFunc("GET /users/{id}/followed-pets", ReqIDParam(s, OperationList, s.EdgeListUserFollowedPetsByID))
	mux.HandleFunc("GET /users/{id}/friends", ReqIDParam(s, OperationList, s.EdgeListUserFriendsByID))
	mux.HandleFunc("GET /users/{id}/friendships", ReqIDParam(s, OperationList, s.EdgeListUserFriendshipsByID))
	mux.HandleFunc("POST /users", ReqParam(s, OperationCreate, s.CreateUser))
	mux.HandleFunc("PATCH /users/{id}", ReqIDParam(s, OperationUpdate, s.UpdateUserByID))
	mux.HandleFunc("DELETE /users/{id}", ReqID(s, OperationDelete, s.DeleteUserByID))

	if !s.config.DisableSpecHandler {
		mux.HandleFunc("GET /openapi.json", s.Spec)
	}

	if !s.config.DisableSpecHandler && !s.config.DisableDocsHandler {
		// If specs are enabled, it's safe to provide documentation, and if they don't override the
		// root endpoint, we can redirect to the docs.
		mux.HandleFunc("GET /", http.RedirectHandler("/docs", http.StatusTemporaryRedirect).ServeHTTP)
		mux.HandleFunc("GET /docs", s.Docs)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleResponse[struct{}](s, w, r, "", nil, ErrEndpointNotFound)
	})
	return mux
}

// ListCategory maps to "GET /categories".
func (s *Server) ListCategory(r *http.Request, p *ListCategoryParams) (*PagedResponse[ent.Category], error) {
	return p.Exec(r.Context(), s.db.Category.Query())
}

// ReadCategoryByID maps to "GET /categories/{id}".
func (s *Server) ReadCategoryByID(r *http.Request, id int) (*ent.Category, error) {
	return EagerLoadCategory(s.db.Category.Query().Where(category.ID(id))).Only(r.Context())
}

// EdgeListCategoryPetsByID maps to "GET /categories/{id}/pets".
func (s *Server) EdgeListCategoryPetsByID(r *http.Request, id int, p *ListPetParams) (*PagedResponse[ent.Pet], error) {
	return p.Exec(r.Context(), s.db.Category.Query().Where(category.ID(id)).QueryPets())
}

// CreateCategory maps to "POST /categories".
func (s *Server) CreateCategory(r *http.Request, p *CreateCategoryParams) (*ent.Category, error) {
	return p.Exec(r.Context(), s.db.Category.Create(), s.db.Category.Query())
}

// UpdateCategoryByID maps to "PATCH /categories/{id}".
func (s *Server) UpdateCategoryByID(r *http.Request, id int, p *UpdateCategoryParams) (*ent.Category, error) {
	return p.Exec(r.Context(), s.db.Category.UpdateOneID(id), s.db.Category.Query())
}

// DeleteCategoryByID maps to "DELETE /categories/{id}".
func (s *Server) DeleteCategoryByID(r *http.Request, id int) (*struct{}, error) {
	return nil, s.db.Category.DeleteOneID(id).Exec(r.Context())
}

// ListFollow maps to "GET /follows".
func (s *Server) ListFollow(r *http.Request, p *ListFollowParams) (*PagedResponse[ent.Follows], error) {
	return p.Exec(r.Context(), s.db.Follows.Query())
}

// CreateFollow maps to "POST /follows".
func (s *Server) CreateFollow(r *http.Request, p *CreateFollowParams) (*ent.Follows, error) {
	return p.Exec(r.Context(), s.db.Follows.Create(), s.db.Follows.Query())
}

// ListFriendship maps to "GET /friendships".
func (s *Server) ListFriendship(r *http.Request, p *ListFriendshipParams) (*PagedResponse[ent.Friendship], error) {
	return p.Exec(r.Context(), s.db.Friendship.Query())
}

// ReadFriendshipByID maps to "GET /friendships/{id}".
func (s *Server) ReadFriendshipByID(r *http.Request, id int) (*ent.Friendship, error) {
	return EagerLoadFriendship(s.db.Friendship.Query().Where(friendship.ID(id))).Only(r.Context())
}

// EdgeReadFriendshipUserByID maps to "GET /friendships/{id}/user".
func (s *Server) EdgeReadFriendshipUserByID(r *http.Request, id int) (*ent.User, error) {
	return EagerLoadUser(s.db.Friendship.Query().Where(friendship.ID(id)).QueryUser()).Only(r.Context())
}

// EdgeReadFriendshipFriendByID maps to "GET /friendships/{id}/friend".
func (s *Server) EdgeReadFriendshipFriendByID(r *http.Request, id int) (*ent.User, error) {
	return EagerLoadUser(s.db.Friendship.Query().Where(friendship.ID(id)).QueryFriend()).Only(r.Context())
}

// CreateFriendship maps to "POST /friendships".
func (s *Server) CreateFriendship(r *http.Request, p *CreateFriendshipParams) (*ent.Friendship, error) {
	return p.Exec(r.Context(), s.db.Friendship.Create(), s.db.Friendship.Query())
}

// UpdateFriendshipByID maps to "PATCH /friendships/{id}".
func (s *Server) UpdateFriendshipByID(r *http.Request, id int, p *UpdateFriendshipParams) (*ent.Friendship, error) {
	return p.Exec(r.Context(), s.db.Friendship.UpdateOneID(id), s.db.Friendship.Query())
}

// DeleteFriendshipByID maps to "DELETE /friendships/{id}".
func (s *Server) DeleteFriendshipByID(r *http.Request, id int) (*struct{}, error) {
	return nil, s.db.Friendship.DeleteOneID(id).Exec(r.Context())
}

// ListPet maps to "GET /pets".
func (s *Server) ListPet(r *http.Request, p *ListPetParams) (*PagedResponse[ent.Pet], error) {
	return p.Exec(r.Context(), s.db.Pet.Query())
}

// ReadPetByID maps to "GET /pets/{id}".
func (s *Server) ReadPetByID(r *http.Request, id int) (*ent.Pet, error) {
	return EagerLoadPet(s.db.Pet.Query().Where(pet.ID(id))).Only(r.Context())
}

// EdgeListPetCategoriesByID maps to "GET /pets/{id}/categories".
func (s *Server) EdgeListPetCategoriesByID(r *http.Request, id int, p *ListCategoryParams) (*PagedResponse[ent.Category], error) {
	return p.Exec(r.Context(), s.db.Pet.Query().Where(pet.ID(id)).QueryCategories())
}

// EdgeReadPetOwnerByID maps to "GET /pets/{id}/owner".
func (s *Server) EdgeReadPetOwnerByID(r *http.Request, id int) (*ent.User, error) {
	return EagerLoadUser(s.db.Pet.Query().Where(pet.ID(id)).QueryOwner()).Only(r.Context())
}

// EdgeListPetFriendsByID maps to "GET /pets/{id}/friends".
func (s *Server) EdgeListPetFriendsByID(r *http.Request, id int, p *ListPetParams) (*PagedResponse[ent.Pet], error) {
	return p.Exec(r.Context(), s.db.Pet.Query().Where(pet.ID(id)).QueryFriends())
}

// EdgeListPetFollowedByByID maps to "GET /pets/{id}/followed-by".
func (s *Server) EdgeListPetFollowedByByID(r *http.Request, id int, p *ListUserParams) (*PagedResponse[ent.User], error) {
	return p.Exec(r.Context(), s.db.Pet.Query().Where(pet.ID(id)).QueryFollowedBy())
}

// CreatePet maps to "POST /pets".
func (s *Server) CreatePet(r *http.Request, p *CreatePetParams) (*ent.Pet, error) {
	return p.Exec(r.Context(), s.db.Pet.Create(), s.db.Pet.Query())
}

// UpdatePetByID maps to "PATCH /pets/{id}".
func (s *Server) UpdatePetByID(r *http.Request, id int, p *UpdatePetParams) (*ent.Pet, error) {
	return p.Exec(r.Context(), s.db.Pet.UpdateOneID(id), s.db.Pet.Query())
}

// DeletePetByID maps to "DELETE /pets/{id}".
func (s *Server) DeletePetByID(r *http.Request, id int) (*struct{}, error) {
	return nil, s.db.Pet.DeleteOneID(id).Exec(r.Context())
}

// ListSetting maps to "GET /settings".
func (s *Server) ListSetting(r *http.Request, p *ListSettingParams) (*PagedResponse[ent.Settings], error) {
	return p.Exec(r.Context(), s.db.Settings.Query())
}

// ReadSettingByID maps to "GET /settings/{id}".
func (s *Server) ReadSettingByID(r *http.Request, id int) (*ent.Settings, error) {
	return EagerLoadSetting(s.db.Settings.Query().Where(settings.ID(id))).Only(r.Context())
}

// EdgeListSettingAdminsByID maps to "GET /settings/{id}/admins".
func (s *Server) EdgeListSettingAdminsByID(r *http.Request, id int, p *ListUserParams) (*PagedResponse[ent.User], error) {
	return p.Exec(r.Context(), s.db.Settings.Query().Where(settings.ID(id)).QueryAdmins())
}

// UpdateSettingByID maps to "PATCH /settings/{id}".
func (s *Server) UpdateSettingByID(r *http.Request, id int, p *UpdateSettingParams) (*ent.Settings, error) {
	return p.Exec(r.Context(), s.db.Settings.UpdateOneID(id), s.db.Settings.Query())
}

// ListUser maps to "GET /users".
func (s *Server) ListUser(r *http.Request, p *ListUserParams) (*PagedResponse[ent.User], error) {
	return p.Exec(r.Context(), s.db.User.Query())
}

// ReadUserByID maps to "GET /users/{id}".
func (s *Server) ReadUserByID(r *http.Request, id int) (*ent.User, error) {
	return EagerLoadUser(s.db.User.Query().Where(user.ID(id))).Only(r.Context())
}

// EdgeListUserPetsByID maps to "GET /users/{id}/pets".
func (s *Server) EdgeListUserPetsByID(r *http.Request, id int, p *ListPetParams) (*PagedResponse[ent.Pet], error) {
	return p.Exec(r.Context(), s.db.User.Query().Where(user.ID(id)).QueryPets())
}

// EdgeListUserFollowedPetsByID maps to "GET /users/{id}/followed-pets".
func (s *Server) EdgeListUserFollowedPetsByID(r *http.Request, id int, p *ListPetParams) (*PagedResponse[ent.Pet], error) {
	return p.Exec(r.Context(), s.db.User.Query().Where(user.ID(id)).QueryFollowedPets())
}

// EdgeListUserFriendsByID maps to "GET /users/{id}/friends".
func (s *Server) EdgeListUserFriendsByID(r *http.Request, id int, p *ListUserParams) (*PagedResponse[ent.User], error) {
	return p.Exec(r.Context(), s.db.User.Query().Where(user.ID(id)).QueryFriends())
}

// EdgeListUserFriendshipsByID maps to "GET /users/{id}/friendships".
func (s *Server) EdgeListUserFriendshipsByID(r *http.Request, id int, p *ListFriendshipParams) (*PagedResponse[ent.Friendship], error) {
	return p.Exec(r.Context(), s.db.User.Query().Where(user.ID(id)).QueryFriendships())
}

// CreateUser maps to "POST /users".
func (s *Server) CreateUser(r *http.Request, p *CreateUserParams) (*ent.User, error) {
	return p.Exec(r.Context(), s.db.User.Create(), s.db.User.Query())
}

// UpdateUserByID maps to "PATCH /users/{id}".
func (s *Server) UpdateUserByID(r *http.Request, id int, p *UpdateUserParams) (*ent.User, error) {
	return p.Exec(r.Context(), s.db.User.UpdateOneID(id), s.db.User.Query())
}

// DeleteUserByID maps to "DELETE /users/{id}".
func (s *Server) DeleteUserByID(r *http.Request, id int) (*struct{}, error) {
	return nil, s.db.User.DeleteOneID(id).Exec(r.Context())
}

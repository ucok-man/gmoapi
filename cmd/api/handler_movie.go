package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ucok-man/gmoapi/internal/data"
	"github.com/ucok-man/gmoapi/internal/validator"
)

// @Summary      List All Movies
// @Description  Retrieve a paginated list of movies with optional filtering by title (full-text search) and genres. Supports sorting by multiple fields.
// @Description
// @Description  **Permissions Required:** `movies:read`
// @Description
// @Description  **Filtering:**
// @Description  - Title: Partial match using PostgreSQL full-text search
// @Description  - Genres: Multiple genres can be specified (comma-separated)
// @Description
// @Description  **Sorting:**
// @Description  - Prefix with `-` for descending order (e.g., `-year`)
// @Description  - Available fields: id, title, year, runtime
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        title      query     string  false  "Filter by movie title (partial match, case-insensitive)"  example(Godfather)
// @Param        genres     query     string  false  "Filter by genres (comma-separated)"  example(drama,crime)
// @Param        page       query     int     false  "Page number (minimum: 1, maximum: 10,000,000)"  default(1)  minimum(1)  maximum(10000000)
// @Param        page_size  query     int     false  "Items per page (minimum: 1, maximum: 100)"  default(20)  minimum(1)  maximum(100)
// @Param        sort       query     string  false  "Sort field"  default(id)  Enums(id, title, year, runtime, -id, -title, -year, -runtime)
// @Security     BearerAuth
// @Success      200  {object}  object{movies=[]data.Movie, metadata=data.Metadata}  "List of movies with pagination metadata"
// @Failure      400  {object}  object{error=string}  "Bad request - invalid query parameters"
// @Failure      401  {object}  object{error=string}  "Unauthorized - missing or invalid authentication token"
// @Failure      403  {object}  object{error=string}  "Forbidden - user account not activated or insufficient permissions"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - validation errors"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /movies [get]
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		Filter data.Filter
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readQueryString(qs, "title", "")
	input.Genres = app.readQueryStrings(qs, "genres", []string{})

	input.Filter.Page = app.readQueryInt(qs, "page", 1, v)
	input.Filter.PageSize = app.readQueryInt(qs, "page_size", 20, v)
	input.Filter.Sort = app.readQueryString(qs, "sort", "id")
	input.Filter.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// @Summary      Get Movie by ID
// @Description  Retrieve detailed information about a specific movie by its unique ID.
// @Description
// @Description  **Permissions Required:** `movies:read`
// @Tags         Movies
// @Produce      json
// @Param        id   path      int  true  "Movie ID"  minimum(1)  example(1)
// @Security     BearerAuth
// @Success      200  {object}  object{movie=data.Movie}  "Movie details"
// @Failure      404  {object}  object{error=string}  "Movie not found"
// @Failure      401  {object}  object{error=string}  "Unauthorized - missing or invalid authentication token"
// @Failure      403  {object}  object{error=string}  "Forbidden - user account not activated or insufficient permissions"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /movies/{id} [get]
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Create New Movie (require movies:write permission)
// @Description  Create a new movie entry in the catalog. All fields are required.
// @Description
// @Description  **Permissions Required:** `movies:write`
// @Description
// @Description  **Validation Rules:**
// @Description  - Title: Required, max 500 characters
// @Description  - Year: Required, between 1888 and current year
// @Description  - Runtime: Required, positive integer, format "123 mins"
// @Description  - Genres: Required, 1-5 unique genres
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie  body      object{title=string, year=int32, runtime=string, genres=[]string}  true  "Movie creation data"
// @Security     BearerAuth
// @Success      201  {object}  object{movie=data.Movie}  "Movie created successfully"
// @Header       201  {string}  Location  "URL of the created movie"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON or invalid data types"
// @Failure      401  {object}  object{error=string}  "Unauthorized - missing or invalid authentication token"
// @Failure      403  {object}  object{error=string}  "Forbidden - user account not activated or insufficient permissions"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - validation errors"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /movies [post]
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Update Movie (require movies:write permission)
// @Description  Update an existing movie using partial update (PATCH). Only provided fields will be updated. Uses optimistic locking to prevent concurrent modification conflicts.
// @Description
// @Description  **Permissions Required:** `movies:write`
// @Description
// @Description  **Validation Rules:** Same as create operation
// @Description
// @Description  **Concurrency Control:** Uses version field for optimistic locking. If the movie has been modified by another request, a 409 Conflict will be returned.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        id     path      int     true  "Movie ID"  minimum(1)  example(1)
// @Param        movie  body      object{title=string, year=int32, runtime=string, genres=[]string}  true  "Movie update data (all fields optional)"
// @Security     BearerAuth
// @Success      200  {object}  object{movie=data.Movie}  "Movie updated successfully"
// @Failure      400  {object}  object{error=string}  "Bad request - malformed JSON or invalid data types"
// @Failure      401  {object}  object{error=string}  "Unauthorized - missing or invalid authentication token"
// @Failure      403  {object}  object{error=string}  "Forbidden - user account not activated or insufficient permissions"
// @Failure      404  {object}  object{error=string}  "Movie not found"
// @Failure      409  {object}  object{error=string}  "Edit conflict - movie has been modified by another request"
// @Failure      422  {object}  object{error=map[string]string}  "Unprocessable entity - validation errors"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /movies/{id} [patch]
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update hanya field yang tidak nil
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// @Summary      Delete Movie (require movies:write permission)
// @Description  Permanently delete a movie from the catalog by its ID. This action cannot be undone.
// @Description
// @Description  **Permissions Required:** `movies:write`
// @Tags         Movies
// @Produce      json
// @Param        id   path      int  true  "Movie ID"  minimum(1)  example(1)
// @Security     BearerAuth
// @Success      200  {object}  object{message=string}  "Movie deleted successfully"
// @Failure      401  {object}  object{error=string}  "Unauthorized - missing or invalid authentication token"
// @Failure      403  {object}  object{error=string}  "Forbidden - user account not activated or insufficient permissions"
// @Failure      404  {object}  object{error=string}  "Movie not found"
// @Failure      429  {object}  object{error=string}  "Too many requests - rate limit exceeded"
// @Failure      500  {object}  object{error=string}  "Internal server error"
// @Router       /movies/{id} [delete]
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

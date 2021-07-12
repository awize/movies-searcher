package controller

import (
	"errors"
	"net/http"
	"testing"

	mocks "github.com/awize/movies-searcher/controller/mock"
	"github.com/awize/movies-searcher/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=movie.go -destination=mock/movie.go -package=mocks

func TestGetMovie(t *testing.T) {
	type result struct {
		statusCode int
	}

	tests := []struct {
		name              string
		id                string
		expectedResult    result
		expectedInt       int
		shouldCallUseCase bool
		expectedError     error
		expectedMovie     *model.Movie
		wantError         bool
	}{
		{
			name:              "NotFound",
			id:                "notInt",
			shouldCallUseCase: true,
			expectedResult: result{
				statusCode: http.StatusNotFound,
			},
			expectedError: errors.New("404 not found"),
			expectedMovie: nil,
			wantError:     true,
		},
		{
			name:              "Successfully",
			id:                "12",
			expectedInt:       12,
			shouldCallUseCase: true,
			expectedResult: result{
				statusCode: http.StatusOK,
			},
			expectedError: nil,
			expectedMovie: &model.Movie{
				ID: 12,
			},
			wantError: false,
		},
		{
			name:              "Successfully",
			id:                "1111",
			expectedInt:       1111,
			shouldCallUseCase: true,
			expectedResult: result{
				statusCode: http.StatusOK,
			},
			expectedError: model.ErrorUnexpected.Err,
			expectedMovie: &model.Movie{
				ID: 12,
			},
			wantError: true,
		},
		// {
		// 	name:              "Malformed",
		// 	id:                "-1",
		// 	shouldCallUseCase: true,
		// 	expectedResult: result{
		// 		statusCode: http.StatusBadRequest,
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			movieUseCase := mocks.NewMockMovieUseCase(ctrl)
			// defer ctrl.Finish()
			// Controller which is tested
			if tt.shouldCallUseCase {
				movieUseCase.EXPECT().GetMovie(tt.expectedInt).Return(tt.expectedMovie, tt.expectedError)
			}

			c := NewMovieController(movieUseCase)

			response, err := c.mu.GetMovie(tt.expectedInt)

			assert.Equal(t, response, tt.expectedMovie)
			// cf := c.GetMovie()
			// response := httptest.NewRecorder()

			// context, _ := gin.CreateTestContext(response)
			// context.Params = []gin.Param{gin.Param{Key: "id", Value: tt.id}}
			// cf(context)

			if tt.wantError {
				assert.NotNil(t, err)
				assert.Equal(t, err, tt.expectedError)
			} else {
				assert.Nil(t, err)
			}

		})
	}

}

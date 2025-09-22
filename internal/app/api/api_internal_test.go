package api

// func TestAPIServer_HandleHello(t *testing.T) {
// 	s := newServer(
// 		NewConfig(),
// 		storage.New(),
// 	)

// 	rec := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/hello", nil)

// 	s.ServeHTTP(rec, req)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// 	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
// 	assert.NotEmpty(t, rec.Body.String())
// 	assert.JSONEq(t, `{"message":"Hello, World!"}`, rec.Body.String())
// }
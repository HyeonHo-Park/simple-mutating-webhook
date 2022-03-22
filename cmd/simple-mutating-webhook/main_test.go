package main

var router = setupRouter()

//func TestMutate(t *testing.T) {
//	tests := []struct {
//		body           gin.H
//		want 		   string
//		message        string
//		wantStatusCode int
//	}{
//		{
//			body: gin.H{
//				"": "",
//			},
//			want: "",
//			message: "OK",
//			wantStatusCode: http.StatusOK,
//		},
//	}
//
//	for _, tt := range tests {
//		test := NewHttpTest(router, t)
//		test.Post("/api/v1/deployment/mutate")
//		test.Send(tt.body)
//		test.Set("Content-Type", "application/json")
//		// test.Set("Authorization", "Bearer "+accessToken)
//		test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {
//			response := parse(rr.Body.Bytes())
//			t.Log(response)
//
//			var mutatedObj interface{}
//			encoded := stringify(response.Data)
//			json.Unmarshal(encoded, &mutatedObj)
//			t.Log(mutatedObj)
//
//			assert.Equal(t, tt.wantStatusCode, rr.Code)
//			assert.Equal(t, tt.want, response.Message)
//		})
//	}
//}
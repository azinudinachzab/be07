package handler

import (
	"a21hc3NpZ25tZW50/client"
	"a21hc3NpZ25tZW50/model"
	//"bufio"
	"context"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"path/filepath"
	"strings"
)

var UserLogin = make(map[string]model.User)

// DESC: func Auth is a middleware to check user login id, only user that already login can pass this middleware
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("user_login_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
			return
		}

		if _, ok := UserLogin[c.Value]; !ok || c.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", c.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DESC: func AuthAdmin is a middleware to check user login role, only admin can pass this middleware
func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("user_login_role")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login role not Admin"})
			return
		}

		if c.Value != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login role not Admin"})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userRole", c.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
		return
	}
	defer r.Body.Close()

	var user model.UserLogin
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
		return
	}

	if user.ID == "" || user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "ID or name is empty"})
		return
	}

	if d, exist := UserLogin[user.ID]; exist {
		http.SetCookie(w, &http.Cookie{Name: "user_login_id", Value: d.ID})
		http.SetCookie(w, &http.Cookie{Name: "user_login_role", Value: d.Role})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.SuccessResponse{
			Username: user.ID,
			Message:  "login success",
		})
	}

	contentUser, err := os.ReadFile("data/users.txt")
	if err != nil {
		panic(err)
		return
	}
	txtUser := string(contentUser)
	txtUserPerBaris := strings.Split(txtUser, "\n")
	userAda := false

	var recordedUser model.User
	for _, val := range txtUserPerBaris {
		usr := strings.Split(val, "_")
		if len(usr) != 4 {
			continue
		}
		if usr[0] == user.ID && usr[1] == user.Name {
			recordedUser = model.User{
				ID:        usr[0],
				Name:      usr[1],
				Role:      usr[3],
				StudyCode: usr[2],
			}
			userAda = true
			break
		}
	}

	if !userAda {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user not found"})
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "user_login_id", Value: recordedUser.ID})
	http.SetCookie(w, &http.Cookie{Name: "user_login_role", Value: recordedUser.Role})
	UserLogin[user.ID] = recordedUser

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: user.ID,
		Message:  "login success",
	})
	return
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
		return
	}
	defer r.Body.Close()

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
		return
	}

	if user.ID == "" || user.Name == "" || user.StudyCode == "" || user.Role == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "ID, name, study code or role is empty"})
		return
	}

	if !validateRole(user.Role) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "role must be admin or user"})
		return
	}

	contentJurusan, err := os.ReadFile("data/list-study.txt")
	if err != nil {
		panic(err)
		return
	}
	txtJurusan := string(contentJurusan)
	txtJurusanPerBaris := strings.Split(txtJurusan, "\n")

	jurusanAda := false
	for _, val := range txtJurusanPerBaris {
		jrs := strings.Split(val, "_")
		if user.StudyCode == jrs[0] {
			jurusanAda = true
			break
		}
	}

	if !jurusanAda {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "study code not found"})
		return
	}

	contentUser, err := os.ReadFile("data/users.txt")
	if err != nil {
		panic(err)
		return
	}
	txtUser := string(contentUser)
	txtUserPerBaris := strings.Split(txtUser, "\n")

	for _, val := range txtUserPerBaris {
		usr := strings.Split(val, "_")
		if len(usr) != 4 {
			continue
		}
		if usr[0] == user.ID {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user id already exist"})
			return
		}
	}

	if txtUser != "" {
		txtUser = txtUser + "\n"
	}
	txtUser = txtUser + user.ID + "_" + user.Name + "_" + user.StudyCode + "_" + user.Role
	err = os.WriteFile("data/users.txt", []byte(txtUser), 0644)
	if err != nil {
		panic(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: user.ID,
		Message:  "register success",
	})
	return
}

func validateRole(role string) bool {
	if role == "admin" || role == "user" {
		return true
	}
	return false
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	if _, exist := UserLogin[userID]; !exist {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
		return
	}

	delete(UserLogin, userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: userID,
		Message:  "logout success",
	})
	return
}

func GetStudyProgram(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	if _, exist := UserLogin[userID]; !exist {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
		return
	}

	content, err := os.ReadFile("data/list-study.txt")
	if err != nil {
		panic(err)
		return
	}

	txt := string(content)
	txtPerBaris := strings.Split(txt, "\n")
	hasilJurusan := make([]model.StudyData, 0)

	for _, val := range txtPerBaris {
		jurusan := strings.Split(val, "_")
		hasilJurusan = append(hasilJurusan, model.StudyData{
			Code: jurusan[0],
			Name: jurusan[1],
		})
	}

	w.WriteHeader(200)
	resp, err := json.Marshal(hasilJurusan)
	if err != nil {
		panic(err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		panic(err)
		return
	}
	return
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	if _, exist := UserLogin[userID]; !exist {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
		return
	}

	err = r.Body.Close()
	if err != nil {
		panic(err)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
		return
	}

	if user.ID == "" || user.Name == "" || user.StudyCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "ID, name, or study code is empty"})
		return
	}

	contentUser, err := os.ReadFile("data/users.txt")
	if err != nil {
		panic(err)
		return
	}
	txtUser := string(contentUser)
	txtUserPerBaris := strings.Split(txtUser, "\n")

	contentJurusan, err := os.ReadFile("data/list-study.txt")
	if err != nil {
		panic(err)
		return
	}
	txtJurusan := string(contentJurusan)
	txtJurusanPerBaris := strings.Split(txtJurusan, "\n")

	for _, val := range txtUserPerBaris {
		usr := strings.Split(val, "_")
		if len(usr) != 3 {
			continue
		}
		if usr[0] == user.ID {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user id already exist"})
			return
		}
	}

	jurusanAda := false
	for _, val := range txtJurusanPerBaris {
		jrs := strings.Split(val, "_")
		if user.StudyCode == jrs[0] {
			jurusanAda = true
			break
		}
	}

	if jurusanAda == false {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "study code not found"})
		return
	}

	if txtUser != "" {
		txtUser = txtUser + "\n"
	}
	txtUser = txtUser + user.ID + "_" + user.Name + "_" + user.StudyCode + "_" + user.Role
	err = os.WriteFile("data/users.txt", []byte(txtUser), 0644)
	if err != nil {
		panic(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: userID,
		Message:  "add user success",
	})
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Method is not allowed!"})
		return
	}

	if _, exist := UserLogin[userID]; !exist {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user login id not found"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user id is empty"})
		return
	}

	contentUser, err := os.ReadFile("data/users.txt")
	if err != nil {
		panic(err)
	}
	txtUser := string(contentUser)
	txtUserPerBaris := strings.Split(txtUser, "\n")

	var idYangDihapus string
	for _, val := range txtUserPerBaris {
		usr := strings.Split(val, "_")
		if len(usr) != 4 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user id not found"})
			return
		}

		if usr[0] == id {
			idYangDihapus = usr[0]
			break
		}
	}

	if idYangDihapus == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "user id not found"})
		return
	}

	var txtBaru string
	for _, val := range txtUserPerBaris {
		usr := strings.Split(val, "_")
		if len(usr) != 3 {
			continue
		}

		if usr[0] == idYangDihapus {
			continue
		}

		txtBaru = val + "\n"
	}
	if txtBaru != "" {
		txtBaru = txtBaru[0 : len(txtBaru)-1]
		err = os.WriteFile("data/users.txt", []byte(txtBaru), 0644)
		if err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Username: userID,
		Message:  "delete success",
	})
	return
}

// DESC: Gunakan variable ini sebagai goroutine di handler GetWeather
var GetWetherByRegionAPI = client.GetWeatherByRegion

func GetWeather(w http.ResponseWriter, r *http.Request) {
	var listRegion = []string{"jakarta", "bandung", "surabaya", "yogyakarta", "medan", "makassar", "manado", "palembang", "semarang", "bali"}
	chErr := make(chan error, 0)
	chMw := make(chan model.MainWeather, len(listRegion))
	for _, val := range listRegion {
		go wrapClientWeather(val, chMw, chErr)
	}

	respBody := make([]model.MainWeather, 0)
	for i:=0;i<len(listRegion);i++{
		err := <- chErr
		weather := <- chMw
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
			return
		}

		respBody = append(respBody, weather)
	}

	w.WriteHeader(200)
	resp, err := json.Marshal(respBody)
	if err != nil {
		panic(err)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		panic(err)
		return
	}
	return
}

func wrapClientWeather(region string, chMw chan model.MainWeather, chErr chan error) {
	mw, err := client.GetWeatherByRegion(region)
	if err != nil {
		chErr <- err
		return
	}

	chMw <- mw
	chErr <- nil
}




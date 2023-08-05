package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/peksinsara/e-voting-RDBMS/database"
	"github.com/peksinsara/e-voting-RDBMS/function"
)

// Candidate represents the candidate data
type Candidate struct {
	CandidateID int    `json:"candidate_id"`
	FullName    string `json:"full_name"`
	District    string `json:"district"`
	ShortBio    string `json:"short_bio"`
}

// AddCandidate handles the admin endpoint to add a new candidate
func AddCandidate(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body into a Candidate struct
	var candidate Candidate
	err := json.NewDecoder(r.Body).Decode(&candidate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	stmt, err := db.Prepare("INSERT INTO Candidate (full_name, district, short_bio) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(candidate.FullName, candidate.District, candidate.ShortBio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	candidateID, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		CandidateID int64 `json:"candidate_id"`
	}{
		CandidateID: candidateID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllCandidates handles the endpoint to get a list of all candidates
func GetAllCandidates(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()

	rows, err := db.Query("SELECT candidate_id, full_name, district, short_bio FROM Candidate")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	candidates := []Candidate{}
	for rows.Next() {
		var candidate Candidate
		err := rows.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.District, &candidate.ShortBio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		candidates = append(candidates, candidate)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}


// DeleteCandidate handles the admin endpoint to delete a candidate by ID
func DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	// Get the candidate ID from the path parameters
	vars := mux.Vars(r)
	candidateID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the candidate from the database
	db := database.GetDB()

	stmt, err := db.Prepare("DELETE FROM Candidate WHERE candidate_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(candidateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Candidate deleted successfully"))
}

// GetAllData handles the admin endpoint to retrieve all data from the database
func GetAllData(w http.ResponseWriter, r *http.Request) {
	// Retrieve all data from the User, Candidate, and Vote tables
	db := database.GetDB()

	userQuery := "SELECT user_id, full_name, mothers_name, email, phone_number, jmbg FROM User"
	rows, err := db.Query(userQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.FullName, &user.MothersName, &user.Email, &user.PhoneNumber, &user.JMBG)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	candidateQuery := "SELECT candidate_id, full_name, district, short_bio FROM Candidate"
	rows, err = db.Query(candidateQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var candidates []Candidate

	for rows.Next() {
		var candidate Candidate
		err = rows.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.District, &candidate.ShortBio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		candidates = append(candidates, candidate)
	}

	voteQuery := "SELECT vote_id, user_id, candidate_id, num_votes FROM Vote"
	rows, err = db.Query(voteQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var votes []function.Vote

	for rows.Next() {
		var vote function.Vote
		err = rows.Scan(&vote.VoteID, &vote.UserID, &vote.CandidateID, &vote.NumVotes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		votes = append(votes, vote)
	}

	// Create a combined data structure
	type AllData struct {
		Users      []User          `json:"users"`
		Candidates []Candidate     `json:"candidates"`
		Votes      []function.Vote `json:"votes"`
	}

	allData := AllData{
		Users:      users,
		Candidates: candidates,
		Votes:      votes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allData)
}

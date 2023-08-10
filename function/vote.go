package function

import (
	"encoding/json"
	"net/http"

	"github.com/peksinsara/e-voting-RDBMS/database"
)

// Vote represents the vote data
type Vote struct {
	VoteID      int `json:"vote_id"`
	UserID      int `json:"user_id"`
	CandidateID int `json:"candidate_id"`
	NumVotes    int `json:"num_votes"`
}

// CastVote handles the user endpoint to cast a vote for a candidate
func CastVote(w http.ResponseWriter, r *http.Request) {
	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM Vote WHERE user_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(vote.UserID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "User has already voted", http.StatusBadRequest)
		return
	}

	stmt, err = db.Prepare("INSERT INTO Vote (user_id, candidate_id, num_votes) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(vote.UserID, vote.CandidateID, vote.NumVotes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the total_votes for the voted candidate
	_, err = db.Exec("UPDATE Candidate SET total_votes = total_votes + 1 WHERE candidate_id = ?", vote.CandidateID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Vote casted successfully"))
}

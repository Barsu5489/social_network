package models

type User struct {
	ID           string `json:"id"`            
	Email        string `json:"email"`          
	PasswordHash string `json:"password_hash"`  
	FirstName    string `json:"first_name"`    
	LastName     string `json:"last_name"`     
	Nickname     string `json:"nickname"`       
	DateOfBirth  string `json:"date_of_birth"` 
	AboutMe      string `json:"about_me"`     
	AvatarURL    string `json:"avatar_url"`    
	IsPrivate    bool   `json:"is_private"`     
	CreatedAt    int64  `json:"created_at"`    
	UpdatedAt    int64  `json:"updated_at"`     
	DeletedAt    *int64 `json:"deleted_at"`     
}
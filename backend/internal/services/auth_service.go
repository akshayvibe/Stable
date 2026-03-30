package services

import (
	"fmt"
	"impact5-backend/internal/database"
	"impact5-backend/internal/models"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

func RegisterUser(client *supabase.Client, email string, password string, name string) (*models.User, error) {
	// 1. Ask Supabase Admin endpoint to create user
	signUpDetails := types.SignupRequest{
		Email:    email,
		Password: password,
		Data: map[string]interface{}{
			"name": name,
		},
	}
	
	// Create user directly (Signup only takes 1 argument in the Supabase Go True client)
	user, err := client.Auth.Signup(signUpDetails)
	if err != nil {
		return nil, fmt.Errorf("supabase registration failed: %v", err)
	}

	// 2. Mirror into local DB. The ID mapped is Supabase's UUID.
	localUser := models.User{
		ID:    user.User.ID.String(),
		Email: email,
		Name:  name,
		Role:  "player",
	}

	if err := database.DB.Create(&localUser).Error; err != nil {
		// Rollback on external provider is complex. Ideally handled with webhook,
		// but since we proxy it, we handle it sequentially here.
		return nil, fmt.Errorf("local db sync failed: %v", err)
	}

	return &localUser, nil
}

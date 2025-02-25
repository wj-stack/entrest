// Code generated by ent, DO NOT EDIT.

package rest

import (
	"context"
	"time"

	github "github.com/google/go-github/v66/github"
	uuid "github.com/google/uuid"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/category"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/follows"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/friendship"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/pet"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/settings"
	"github.com/lrstanley/entrest/_examples/kitchensink/internal/database/ent/user"
	schema "github.com/lrstanley/entrest/_examples/kitchensink/internal/database/schema"
)

// CreateCategoryParams defines parameters for creating a Category via a POST request.
type CreateCategoryParams struct {
	Name     string   `json:"name"`
	Nillable *string  `json:"nillable"`
	Strings  []string `json:"strings,omitempty"`
	Ints     []int    `json:"ints,omitempty"`
	Pets     []int    `json:"pets,omitempty"`
}

func (c *CreateCategoryParams) ApplyInputs(builder *ent.CategoryCreate) *ent.CategoryCreate {
	builder.SetName(c.Name)
	if c.Nillable != nil {
		builder.SetNillable(*c.Nillable)
	}
	if c.Strings != nil {
		builder.SetStrings(c.Strings)
	}
	if c.Ints != nil {
		builder.SetInts(c.Ints)
	}
	builder.AddPetIDs(c.Pets...)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreateCategoryParams) Exec(ctx context.Context, builder *ent.CategoryCreate, query *ent.CategoryQuery) (*ent.Category, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	return EagerLoadCategory(query.Where(category.ID(result.ID))).Only(ctx)
}

// CreateFollowParams defines parameters for creating a Follow via a POST request.
type CreateFollowParams struct {
	UserID uuid.UUID `json:"user_id"`
	PetID  int       `json:"pet_id"`
}

func (c *CreateFollowParams) ApplyInputs(builder *ent.FollowsCreate) *ent.FollowsCreate {
	builder.SetUserID(c.UserID)
	builder.SetPetID(c.PetID)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreateFollowParams) Exec(ctx context.Context, builder *ent.FollowsCreate, query *ent.FollowsQuery) (*ent.Follows, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	// Since Follow entities have a composite ID, we have to query by all known FK fields.
	return EagerLoadFollow(query.Where(
		follows.UserIDEQ(result.UserID),
		follows.PetIDEQ(result.PetID),
	)).Only(ctx)
}

// CreateFriendshipParams defines parameters for creating a Friendship via a POST request.
type CreateFriendshipParams struct {
	CreatedAt *time.Time `json:"created_at"`
	UserID    uuid.UUID  `json:"user_id"`
	FriendID  uuid.UUID  `json:"friend_id"`
}

func (c *CreateFriendshipParams) ApplyInputs(builder *ent.FriendshipCreate) *ent.FriendshipCreate {
	if c.CreatedAt != nil {
		builder.SetCreatedAt(*c.CreatedAt)
	}
	builder.SetUserID(c.UserID)
	builder.SetFriendID(c.FriendID)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreateFriendshipParams) Exec(ctx context.Context, builder *ent.FriendshipCreate, query *ent.FriendshipQuery) (*ent.Friendship, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	return EagerLoadFriendship(query.Where(friendship.ID(result.ID))).Only(ctx)
}

// CreatePetParams defines parameters for creating a Pet via a POST request.
type CreatePetParams struct {
	Name      string   `json:"name"`
	Nicknames []string `json:"nicknames,omitempty"`
	Age       int      `json:"age"`
	Type      pet.Type `json:"type"`
	// Categories that the pet belongs to.
	Categories []int `json:"categories,omitempty"`
	// The user that owns the pet.
	Owner *uuid.UUID `json:"owner,omitempty"`
	// Pets that this pet is friends with.
	Friends []int `json:"friends,omitempty"`
	// Users that this pet is followed by.
	FollowedBy []uuid.UUID `json:"followed_by,omitempty"`
}

func (c *CreatePetParams) ApplyInputs(builder *ent.PetCreate) *ent.PetCreate {
	builder.SetName(c.Name)
	if c.Nicknames != nil {
		builder.SetNicknames(c.Nicknames)
	}
	builder.SetAge(c.Age)
	builder.SetType(c.Type)
	builder.AddCategoryIDs(c.Categories...)
	if c.Owner != nil {
		builder.SetOwnerID(*c.Owner)
	}
	builder.AddFriendIDs(c.Friends...)
	builder.AddFollowedByIDs(c.FollowedBy...)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreatePetParams) Exec(ctx context.Context, builder *ent.PetCreate, query *ent.PetQuery) (*ent.Pet, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	return EagerLoadPet(query.Where(pet.ID(result.ID))).Only(ctx)
}

// CreateSettingParams defines parameters for creating a Setting via a POST request.
type CreateSettingParams struct {
	// Global banner text to apply to the frontend.
	GlobalBanner *string `json:"global_banner,omitempty"`
	// Administrators for the platform.
	Admins []uuid.UUID `json:"admins,omitempty"`
}

func (c *CreateSettingParams) ApplyInputs(builder *ent.SettingsCreate) *ent.SettingsCreate {
	if c.GlobalBanner != nil {
		builder.SetGlobalBanner(*c.GlobalBanner)
	}
	builder.AddAdminIDs(c.Admins...)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreateSettingParams) Exec(ctx context.Context, builder *ent.SettingsCreate, query *ent.SettingsQuery) (*ent.Settings, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	return EagerLoadSetting(query.Where(settings.ID(result.ID))).Only(ctx)
}

// CreateUserParams defines parameters for creating a User via a POST request.
type CreateUserParams struct {
	ID *uuid.UUID `json:"id"`
	// Name of the user.
	Name string `json:"name"`
	// Type of object being defined (user or system which is for internal usecases).
	Type *user.Type `json:"type"`
	// Full name if USER, otherwise null.
	Description *string `json:"description,omitempty"`
	// If the user is still in the source system.
	Enabled *bool `json:"enabled"`
	// Email associated with the user. Note that not all users have an associated email address.
	Email *string `json:"email,omitempty"`
	// Avatar data for the user. This should generally only apply to the USER user type.
	Avatar []byte `json:"avatar,omitempty"`
	// Hashed password for the user, this shouldn't be readable in the spec anywhere.
	PasswordHashed string `json:"password_hashed"`
	// The github user raw JSON data.
	GithubData          *github.User          `json:"github_data,omitempty"`
	ProfileURL          *schema.ExampleValuer `json:"profile_url,omitempty"`
	LastAuthenticatedAt *time.Time            `json:"last_authenticated_at,omitempty"`
	// Pets owned by the user.
	Pets []int `json:"pets,omitempty"`
	// Pets that the user is following.
	FollowedPets []int `json:"followed_pets,omitempty"`
	// Friends of the user.
	Friends     []uuid.UUID `json:"friends,omitempty"`
	Friendships []int       `json:"friendships,omitempty"`
}

func (c *CreateUserParams) ApplyInputs(builder *ent.UserCreate) *ent.UserCreate {
	if c.ID != nil {
		builder.SetID(*c.ID)
	}
	builder.SetName(c.Name)
	if c.Type != nil {
		builder.SetType(*c.Type)
	}
	if c.Description != nil {
		builder.SetDescription(*c.Description)
	}
	if c.Enabled != nil {
		builder.SetEnabled(*c.Enabled)
	}
	if c.Email != nil {
		builder.SetEmail(*c.Email)
	}
	if c.Avatar != nil {
		builder.SetAvatar(c.Avatar)
	}
	builder.SetPasswordHashed(c.PasswordHashed)
	if c.GithubData != nil {
		builder.SetGithubData(c.GithubData)
	}
	if c.ProfileURL != nil {
		builder.SetProfileURL(c.ProfileURL)
	}
	if c.LastAuthenticatedAt != nil {
		builder.SetLastAuthenticatedAt(*c.LastAuthenticatedAt)
	}
	builder.AddPetIDs(c.Pets...)
	builder.AddFollowedPetIDs(c.FollowedPets...)
	builder.AddFriendIDs(c.Friends...)
	builder.AddFriendshipIDs(c.Friendships...)
	return builder
}

// Exec wraps all logic (mapping all provided values to the builder), creates the entity,
// and does another query (using provided query as base) to get the entity, with all eager
// loaded edges.
func (c *CreateUserParams) Exec(ctx context.Context, builder *ent.UserCreate, query *ent.UserQuery) (*ent.User, error) {
	result, err := c.ApplyInputs(builder).Save(ctx)
	if err != nil {
		return nil, err
	}
	return EagerLoadUser(query.Where(user.ID(result.ID))).Only(ctx)
}

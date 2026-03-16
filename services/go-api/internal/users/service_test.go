package users

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// --- Mock Repository ---

type mockRepo struct {
	listFn             func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error)
	getByIDFn          func(ctx context.Context, id uuid.UUID) (*UserRow, error)
	getByEmailFn       func(ctx context.Context, email string) (*UserRow, error)
	createFn           func(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error)
	updateFn           func(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error)
	deactivateFn       func(ctx context.Context, id uuid.UUID) error
	resetPasswordFn    func(ctx context.Context, id uuid.UUID, newPassword string) error
	updateSelfFn       func(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error)
	listDepartmentsFn  func(ctx context.Context) ([]DepartmentRow, error)
	createDepartmentFn func(ctx context.Context, name string) (*DepartmentRow, error)
}

func (m *mockRepo) List(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, limit, filter)
	}
	return nil, 0, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*UserRow, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*UserRow, error) {
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

func (m *mockRepo) Create(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req, tempPassword)
	}
	return nil, nil
}

func (m *mockRepo) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, req)
	}
	return nil, nil
}

func (m *mockRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	if m.deactivateFn != nil {
		return m.deactivateFn(ctx, id)
	}
	return nil
}

func (m *mockRepo) ResetPassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	if m.resetPasswordFn != nil {
		return m.resetPasswordFn(ctx, id, newPassword)
	}
	return nil
}

func (m *mockRepo) UpdateSelf(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error) {
	if m.updateSelfFn != nil {
		return m.updateSelfFn(ctx, id, firstName, lastName)
	}
	return nil, nil
}

func (m *mockRepo) ListDepartments(ctx context.Context) ([]DepartmentRow, error) {
	if m.listDepartmentsFn != nil {
		return m.listDepartmentsFn(ctx)
	}
	return nil, nil
}

func (m *mockRepo) CreateDepartment(ctx context.Context, name string) (*DepartmentRow, error) {
	if m.createDepartmentFn != nil {
		return m.createDepartmentFn(ctx, name)
	}
	return nil, nil
}

// --- Test Helpers ---

func testUserRow(id uuid.UUID, email, role string) UserRow {
	return UserRow{
		ID:                 id,
		Email:              email,
		FirstName:          "Test",
		LastName:           "User",
		Role:               role,
		IsActive:           true,
		MustChangePassword: false,
	}
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

// ==================== List Tests ====================

func TestList_Success(t *testing.T) {
	uid1 := uuid.New()
	uid2 := uuid.New()
	rows := []UserRow{
		testUserRow(uid1, "user1@test.com", "hr"),
		testUserRow(uid2, "user2@test.com", "supervisor"),
	}

	repo := &mockRepo{
		listFn: func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
			if page != 1 || limit != 20 {
				t.Errorf("List() page=%d limit=%d, want 1/20", page, limit)
			}
			return rows, 2, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	result, err := svc.List(context.Background(), 1, 20, ListFilter{})

	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(result.Data) != 2 {
		t.Errorf("List() returned %d users, want 2", len(result.Data))
	}
	if result.Total != 2 {
		t.Errorf("List() total = %d, want 2", result.Total)
	}
	if result.TotalPages != 1 {
		t.Errorf("List() totalPages = %d, want 1", result.TotalPages)
	}
}

func TestList_PaginationDefaults(t *testing.T) {
	repo := &mockRepo{
		listFn: func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
			if page != 1 {
				t.Errorf("List() should default page to 1, got %d", page)
			}
			if limit != 20 {
				t.Errorf("List() should default limit to 20, got %d", limit)
			}
			return nil, 0, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.List(context.Background(), 0, -5, ListFilter{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
}

func TestList_LimitClamped(t *testing.T) {
	repo := &mockRepo{
		listFn: func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
			if limit != 20 {
				t.Errorf("List() should clamp limit > 100 to 20, got %d", limit)
			}
			return nil, 0, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	_, _ = svc.List(context.Background(), 1, 200, ListFilter{})
}

func TestList_TotalPagesCalculation(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		limit     int
		wantPages int
	}{
		{"exact division", 40, 20, 2},
		{"remainder", 41, 20, 3},
		{"single page", 5, 20, 1},
		{"empty", 0, 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				listFn: func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
					return make([]UserRow, 0), tt.total, nil
				},
			}
			svc := NewServiceWithDeps(repo)
			result, _ := svc.List(context.Background(), 1, tt.limit, ListFilter{})
			if result.TotalPages != tt.wantPages {
				t.Errorf("TotalPages = %d, want %d", result.TotalPages, tt.wantPages)
			}
		})
	}
}

func TestList_RepoError(t *testing.T) {
	repo := &mockRepo{
		listFn: func(ctx context.Context, page, limit int, filter ListFilter) ([]UserRow, int, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.List(context.Background(), 1, 20, ListFilter{})

	if err == nil {
		t.Fatal("List() should return error on repo failure")
	}
}

// ==================== GetByID Tests ====================

func TestGetByID_Success(t *testing.T) {
	uid := uuid.New()
	row := testUserRow(uid, "hr@test.com", "hr")

	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			if id == uid {
				return &row, nil
			}
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	user, err := svc.GetByID(context.Background(), uid)

	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if user.ID != uid {
		t.Errorf("GetByID() ID = %v, want %v", user.ID, uid)
	}
	if user.Email != "hr@test.com" {
		t.Errorf("GetByID() Email = %q, want %q", user.Email, "hr@test.com")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.GetByID(context.Background(), uuid.New())

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("GetByID() error = %v, want ErrUserNotFound", err)
	}
}

// ==================== Create Tests ====================

func TestCreate_Success(t *testing.T) {
	createdRow := testUserRow(uuid.New(), "new@test.com", "hr")
	createdRow.MustChangePassword = true

	var capturedTempPwd string
	repo := &mockRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return nil, nil // no existing user
		},
		createFn: func(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error) {
			capturedTempPwd = tempPassword
			return &createdRow, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	user, tempPwd, err := svc.Create(context.Background(), CreateUserRequest{
		Email:     "new@test.com",
		FirstName: "New",
		LastName:  "User",
		Role:      "hr",
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.Email != "new@test.com" {
		t.Errorf("Create() Email = %q, want %q", user.Email, "new@test.com")
	}
	if tempPwd == "" {
		t.Error("Create() returned empty temporary password")
	}
	if len(tempPwd) != 12 {
		t.Errorf("Create() temp password length = %d, want 12", len(tempPwd))
	}
	if tempPwd != capturedTempPwd {
		t.Errorf("Create() temp password mismatch: returned %q, passed to repo %q", tempPwd, capturedTempPwd)
	}
}

func TestCreate_EmailExists(t *testing.T) {
	existingRow := testUserRow(uuid.New(), "existing@test.com", "hr")

	repo := &mockRepo{
		getByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
			return &existingRow, nil // email already taken
		},
	}

	svc := NewServiceWithDeps(repo)
	_, _, err := svc.Create(context.Background(), CreateUserRequest{
		Email:     "existing@test.com",
		FirstName: "New",
		LastName:  "User",
		Role:      "hr",
	})

	if !errors.Is(err, ErrEmailExists) {
		t.Errorf("Create() error = %v, want ErrEmailExists", err)
	}
}

func TestCreate_InvalidRole(t *testing.T) {
	repo := &mockRepo{}

	svc := NewServiceWithDeps(repo)

	invalidRoles := []string{"viewer", "manager", "applicant", ""}
	for _, role := range invalidRoles {
		_, _, err := svc.Create(context.Background(), CreateUserRequest{
			Email:     "new@test.com",
			FirstName: "New",
			LastName:  "User",
			Role:      role,
		})
		if !errors.Is(err, ErrInvalidRole) {
			t.Errorf("Create(role=%q) error = %v, want ErrInvalidRole", role, err)
		}
	}
}

func TestCreate_ValidRoles(t *testing.T) {
	validRoles := []string{"admin", "hr", "supervisor"}
	for _, role := range validRoles {
		createdRow := testUserRow(uuid.New(), "new@test.com", role)

		repo := &mockRepo{
			getByEmailFn: func(ctx context.Context, email string) (*UserRow, error) {
				return nil, nil
			},
			createFn: func(ctx context.Context, req CreateUserRequest, tempPassword string) (*UserRow, error) {
				return &createdRow, nil
			},
		}

		svc := NewServiceWithDeps(repo)
		_, _, err := svc.Create(context.Background(), CreateUserRequest{
			Email:     "new@test.com",
			FirstName: "New",
			LastName:  "User",
			Role:      role,
		})
		if err != nil {
			t.Errorf("Create(role=%q) should succeed, got error = %v", role, err)
		}
	}
}

// ==================== Update Tests ====================

func TestUpdate_Success(t *testing.T) {
	uid := uuid.New()
	updatedRow := testUserRow(uid, "hr@test.com", "supervisor")
	updatedRow.FirstName = "Updated"

	repo := &mockRepo{
		updateFn: func(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error) {
			return &updatedRow, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	user, err := svc.Update(context.Background(), uid, UpdateUserRequest{
		FirstName: strPtr("Updated"),
		Role:      strPtr("supervisor"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if user.FirstName != "Updated" {
		t.Errorf("Update() FirstName = %q, want %q", user.FirstName, "Updated")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error) {
			return nil, pgx.ErrNoRows
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.Update(context.Background(), uuid.New(), UpdateUserRequest{
		FirstName: strPtr("Updated"),
	})

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Update() error = %v, want ErrUserNotFound", err)
	}
}

func TestUpdate_InvalidRole(t *testing.T) {
	repo := &mockRepo{}

	svc := NewServiceWithDeps(repo)
	badRole := "invalid_role"
	_, err := svc.Update(context.Background(), uuid.New(), UpdateUserRequest{
		Role: &badRole,
	})

	if !errors.Is(err, ErrInvalidRole) {
		t.Errorf("Update() error = %v, want ErrInvalidRole", err)
	}
}

func TestUpdate_NilRow(t *testing.T) {
	repo := &mockRepo{
		updateFn: func(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*UserRow, error) {
			return nil, nil // returns nil without error
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.Update(context.Background(), uuid.New(), UpdateUserRequest{
		FirstName: strPtr("Updated"),
	})

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Update() with nil row error = %v, want ErrUserNotFound", err)
	}
}

// ==================== Delete Tests ====================

func TestDelete_Success(t *testing.T) {
	targetID := uuid.New()
	callerID := uuid.New()

	repo := &mockRepo{
		deactivateFn: func(ctx context.Context, id uuid.UUID) error {
			if id != targetID {
				t.Errorf("Deactivate() called with %v, want %v", id, targetID)
			}
			return nil
		},
	}

	svc := NewServiceWithDeps(repo)
	err := svc.Delete(context.Background(), targetID, callerID.String())

	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestDelete_SelfDeactivation(t *testing.T) {
	userID := uuid.New()

	repo := &mockRepo{}

	svc := NewServiceWithDeps(repo)
	err := svc.Delete(context.Background(), userID, userID.String())

	if !errors.Is(err, ErrSelfDeactivate) {
		t.Errorf("Delete() self-deactivation error = %v, want ErrSelfDeactivate", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepo{
		deactivateFn: func(ctx context.Context, id uuid.UUID) error {
			return pgx.ErrNoRows
		},
	}

	svc := NewServiceWithDeps(repo)
	err := svc.Delete(context.Background(), uuid.New(), uuid.New().String())

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("Delete() error = %v, want ErrUserNotFound", err)
	}
}

// ==================== AdminResetPassword Tests ====================

func TestAdminResetPassword_Success(t *testing.T) {
	targetID := uuid.New()

	var capturedPwd string
	repo := &mockRepo{
		resetPasswordFn: func(ctx context.Context, id uuid.UUID, newPassword string) error {
			capturedPwd = newPassword
			return nil
		},
	}

	svc := NewServiceWithDeps(repo)
	err := svc.AdminResetPassword(context.Background(), targetID, "newpassword1")

	if err != nil {
		t.Fatalf("AdminResetPassword() error = %v", err)
	}
	if capturedPwd != "newpassword1" {
		t.Errorf("AdminResetPassword() password = %q, want %q", capturedPwd, "newpassword1")
	}
}

func TestAdminResetPassword_NotFound(t *testing.T) {
	repo := &mockRepo{
		resetPasswordFn: func(ctx context.Context, id uuid.UUID, newPassword string) error {
			return pgx.ErrNoRows
		},
	}

	svc := NewServiceWithDeps(repo)
	err := svc.AdminResetPassword(context.Background(), uuid.New(), "newpassword1")

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("AdminResetPassword() error = %v, want ErrUserNotFound", err)
	}
}

// ==================== GetMe Tests ====================

func TestGetMe_Success(t *testing.T) {
	uid := uuid.New()
	row := testUserRow(uid, "me@test.com", "hr")

	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*UserRow, error) {
			if id == uid {
				return &row, nil
			}
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	user, err := svc.GetMe(context.Background(), uid.String())

	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}
	if user.Email != "me@test.com" {
		t.Errorf("GetMe() Email = %q, want %q", user.Email, "me@test.com")
	}
}

func TestGetMe_InvalidID(t *testing.T) {
	svc := NewServiceWithDeps(&mockRepo{})
	_, err := svc.GetMe(context.Background(), "not-a-uuid")

	if err == nil {
		t.Fatal("GetMe() should return error for invalid UUID")
	}
}

// ==================== UpdateMe Tests ====================

func TestUpdateMe_Success(t *testing.T) {
	uid := uuid.New()
	updatedRow := testUserRow(uid, "me@test.com", "hr")
	updatedRow.FirstName = "NewFirst"
	updatedRow.LastName = "NewLast"

	repo := &mockRepo{
		updateSelfFn: func(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error) {
			if firstName != "NewFirst" || lastName != "NewLast" {
				t.Errorf("UpdateSelf() names = %q %q, want NewFirst NewLast", firstName, lastName)
			}
			return &updatedRow, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	user, err := svc.UpdateMe(context.Background(), uid.String(), UpdateMeRequest{
		FirstName: "NewFirst",
		LastName:  "NewLast",
	})

	if err != nil {
		t.Fatalf("UpdateMe() error = %v", err)
	}
	if user.FirstName != "NewFirst" {
		t.Errorf("UpdateMe() FirstName = %q, want %q", user.FirstName, "NewFirst")
	}
}

func TestUpdateMe_NotFound(t *testing.T) {
	repo := &mockRepo{
		updateSelfFn: func(ctx context.Context, id uuid.UUID, firstName, lastName string) (*UserRow, error) {
			return nil, nil // user not found
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.UpdateMe(context.Background(), uuid.New().String(), UpdateMeRequest{
		FirstName: "Test",
		LastName:  "User",
	})

	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("UpdateMe() error = %v, want ErrUserNotFound", err)
	}
}

// ==================== ListDepartments Tests ====================

func TestListDepartments_Success(t *testing.T) {
	rows := []DepartmentRow{
		{ID: uuid.New(), Name: "Engineering"},
		{ID: uuid.New(), Name: "Marketing"},
	}

	repo := &mockRepo{
		listDepartmentsFn: func(ctx context.Context) ([]DepartmentRow, error) {
			return rows, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	depts, err := svc.ListDepartments(context.Background())

	if err != nil {
		t.Fatalf("ListDepartments() error = %v", err)
	}
	if len(depts) != 2 {
		t.Errorf("ListDepartments() returned %d, want 2", len(depts))
	}
	if depts[0].Name != "Engineering" {
		t.Errorf("ListDepartments()[0].Name = %q, want %q", depts[0].Name, "Engineering")
	}
}

func TestListDepartments_Empty(t *testing.T) {
	repo := &mockRepo{
		listDepartmentsFn: func(ctx context.Context) ([]DepartmentRow, error) {
			return nil, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	depts, err := svc.ListDepartments(context.Background())

	if err != nil {
		t.Fatalf("ListDepartments() error = %v", err)
	}
	if len(depts) != 0 {
		t.Errorf("ListDepartments() returned %d, want 0", len(depts))
	}
}

// ==================== CreateDepartment Tests ====================

func TestCreateDepartment_Success(t *testing.T) {
	deptID := uuid.New()
	repo := &mockRepo{
		createDepartmentFn: func(ctx context.Context, name string) (*DepartmentRow, error) {
			return &DepartmentRow{ID: deptID, Name: name}, nil
		},
	}

	svc := NewServiceWithDeps(repo)
	dept, err := svc.CreateDepartment(context.Background(), CreateDepartmentRequest{Name: "Engineering"})

	if err != nil {
		t.Fatalf("CreateDepartment() error = %v", err)
	}
	if dept.Name != "Engineering" {
		t.Errorf("CreateDepartment() Name = %q, want %q", dept.Name, "Engineering")
	}
	if dept.ID != deptID {
		t.Errorf("CreateDepartment() ID = %v, want %v", dept.ID, deptID)
	}
}

func TestCreateDepartment_DuplicateName(t *testing.T) {
	repo := &mockRepo{
		createDepartmentFn: func(ctx context.Context, name string) (*DepartmentRow, error) {
			return nil, errors.New("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)")
		},
	}

	svc := NewServiceWithDeps(repo)
	_, err := svc.CreateDepartment(context.Background(), CreateDepartmentRequest{Name: "Engineering"})

	if !errors.Is(err, ErrDeptExists) {
		t.Errorf("CreateDepartment() error = %v, want ErrDeptExists", err)
	}
}

// ==================== rowToUser Tests ====================

func TestRowToUser_MapsAllFields(t *testing.T) {
	deptID := uuid.New()
	row := UserRow{
		ID:                 uuid.New(),
		Email:              "test@test.com",
		PasswordHash:       "should_not_appear",
		FirstName:          "First",
		LastName:           "Last",
		Role:               "hr",
		DepartmentID:       &deptID,
		IsActive:           true,
		MustChangePassword: false,
	}

	user := rowToUser(row)

	if user.ID != row.ID {
		t.Errorf("rowToUser() ID mismatch")
	}
	if user.Email != row.Email {
		t.Errorf("rowToUser() Email mismatch")
	}
	if user.FirstName != row.FirstName {
		t.Errorf("rowToUser() FirstName mismatch")
	}
	if user.DepartmentID == nil || *user.DepartmentID != deptID {
		t.Errorf("rowToUser() DepartmentID mismatch")
	}
}

// ==================== generateTempPassword Tests ====================

func TestGenerateTempPassword_Length(t *testing.T) {
	pwd, err := generateTempPassword(12)
	if err != nil {
		t.Fatalf("generateTempPassword() error = %v", err)
	}
	if len(pwd) != 12 {
		t.Errorf("generateTempPassword() length = %d, want 12", len(pwd))
	}
}

func TestGenerateTempPassword_Alphanumeric(t *testing.T) {
	pwd, err := generateTempPassword(100)
	if err != nil {
		t.Fatalf("generateTempPassword() error = %v", err)
	}
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, c := range pwd {
		found := false
		for _, valid := range chars {
			if c == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("generateTempPassword() contains invalid char: %c", c)
		}
	}
}

func TestGenerateTempPassword_Unique(t *testing.T) {
	pwd1, _ := generateTempPassword(12)
	pwd2, _ := generateTempPassword(12)
	if pwd1 == pwd2 {
		t.Error("generateTempPassword() should produce unique values")
	}
}

// ==================== isDuplicateKeyError Tests ====================

func TestIsDuplicateKeyError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"duplicate key", errors.New("ERROR: duplicate key value violates unique constraint"), true},
		{"sqlstate 23505", errors.New("SQLSTATE 23505"), true},
		{"unrelated error", errors.New("connection refused"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicateKeyError(tt.err)
			if got != tt.want {
				t.Errorf("isDuplicateKeyError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
